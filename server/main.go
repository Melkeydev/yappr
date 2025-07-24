package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/melkeydev/chat-go/db"
	"github.com/melkeydev/chat-go/db/migrations"
	coreHandler "github.com/melkeydev/chat-go/internal/api/handler/core"
	statsHandler "github.com/melkeydev/chat-go/internal/api/handler/stats"
	userHandler "github.com/melkeydev/chat-go/internal/api/handler/user"
	roomRepo "github.com/melkeydev/chat-go/internal/repo/room"
	statsRepo "github.com/melkeydev/chat-go/internal/repo/stats"
	repository "github.com/melkeydev/chat-go/internal/repo/user"
	"github.com/melkeydev/chat-go/internal/service/pinnedrooms"
	statsService "github.com/melkeydev/chat-go/internal/service/stats"
	service "github.com/melkeydev/chat-go/internal/service/user"
	"github.com/melkeydev/chat-go/internal/ws"
	"github.com/melkeydev/chat-go/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize DB connection: %s", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Run migrations
	if err := migration.RunMigrations(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set up Repositories
	userRepo := repository.NewUserRepository(dbConn)
	statsRepository := statsRepo.NewStatsRepository(dbConn)

	// Set up Services
	userService := service.NewUserService(userRepo)
	statsServ := statsService.NewStatsService(statsRepository)
	wsService := ws.NewCore(dbConn)

	// Set up Handlers
	userHandler := userHandler.NewUserHandler(userService)
	coreHandler := coreHandler.NewCoreHandler(wsService)
	statsHand := statsHandler.NewStatsHandler(statsServ)

	// run it in a separate go routine
	go wsService.Run()
	
	// Initialize pinned rooms on startup
	pinnedRoomsService := pinnedrooms.NewPinnedRoomsService(dbConn, wsService)
	if err := pinnedRoomsService.CheckAndRefreshPinnedRooms(context.Background()); err != nil {
		log.Printf("Failed to initialize pinned rooms: %v", err)
	}
	
	// Start background job to clean up expired rooms
	go startRoomCleanupJob(dbConn, wsService)

	router := router.SetupRouter(userHandler, coreHandler, statsHand)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// startRoomCleanupJob runs a background job that deletes expired rooms every 5 minutes
func startRoomCleanupJob(db *sql.DB, wsCore *ws.Core) {
	roomRepository := roomRepo.NewRoomRepository(db)
	pinnedRoomsService := pinnedrooms.NewPinnedRoomsService(db, wsCore)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	// Run cleanup immediately on startup
	cleanupRooms(roomRepository, pinnedRoomsService)
	
	for range ticker.C {
		cleanupRooms(roomRepository, pinnedRoomsService)
	}
}

func cleanupRooms(roomRepository *roomRepo.RoomRepository, pinnedRoomsService *pinnedrooms.PinnedRoomsService) {
	ctx := context.Background()
	deletedCount, err := roomRepository.DeleteExpiredRooms(ctx)
	if err != nil {
		log.Printf("Error deleting expired rooms: %v", err)
		return
	}
	
	if deletedCount > 0 {
		log.Printf("Deleted %d expired rooms", deletedCount)
	}
	
	// Check and refresh pinned rooms if needed
	if err := pinnedRoomsService.CheckAndRefreshPinnedRooms(ctx); err != nil {
		log.Printf("Error refreshing pinned rooms: %v", err)
	}
}
