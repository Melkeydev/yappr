package main

import (
	"log"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/melkeydev/chat-go/db"
	"github.com/melkeydev/chat-go/db/migrations"
	coreHandler "github.com/melkeydev/chat-go/internal/api/handler/core"
	userHandler "github.com/melkeydev/chat-go/internal/api/handler/user"
	repository "github.com/melkeydev/chat-go/internal/repo/user"
	service "github.com/melkeydev/chat-go/internal/service/user"
	"github.com/melkeydev/chat-go/internal/ws"
	"github.com/melkeydev/chat-go/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not initialize .env filed: %s", err)
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

	// Set up Services
	userService := service.NewUserService(userRepo)
	wsService := ws.NewCore()

	// Set up Handlers
	userHandler := userHandler.NewUserHandler(userService)
	coreHandler := coreHandler.NewCoreHandler(wsService)

	// run it in a separate go routine
	go wsService.Run()

	router := router.SetupRouter(userHandler, coreHandler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
