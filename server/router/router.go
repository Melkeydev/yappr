package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	corehandler "github.com/melkeydev/chat-go/internal/api/handler/core"
	statshandler "github.com/melkeydev/chat-go/internal/api/handler/stats"
	userhandler "github.com/melkeydev/chat-go/internal/api/handler/user"
	authmiddleware "github.com/melkeydev/chat-go/middleware"
)

func SetupRouter(userH *userhandler.UserHandler, coreH *corehandler.CoreHandler, statsH *statshandler.StatsHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/users", func(u chi.Router) {
		u.Post("/signup", userH.CreateUser)
		u.Post("/login", userH.Login)
		u.Get("/logout", userH.Logout)
		
		// Protected routes
		u.Group(func(r chi.Router) {
			r.Use(authmiddleware.JWTAuth)
			r.Put("/username", userH.UpdateUsername)
		})
	})

	r.Route("/api/stats", func(s chi.Router) {
		// Protected routes requiring authentication
		s.Group(func(r chi.Router) {
			r.Use(authmiddleware.JWTAuth)
			r.Post("/checkin", statsH.CheckIn)
			r.Post("/upvote", statsH.GiveUpvote)
		})
		
		// Public routes (with optional auth for viewing permissions)
		s.Group(func(r chi.Router) {
			r.Use(authmiddleware.OptionalJWTAuth)
			r.Get("/profile/{userId}", statsH.GetUserProfile)
		})
	})

	r.Route("/ws", func(u chi.Router) {
		// Protected route for creating rooms
		u.Group(func(r chi.Router) {
			r.Use(authmiddleware.OptionalJWTAuth)
			r.Post("/createRoom", coreH.CreateRoom)
		})
		
		u.Get("/joinRoom/{roomId}", coreH.JoinRoom)
		u.Get("/getRooms", coreH.GetRooms)
		u.Get("/getClients/{roomId}", coreH.GetClients)
	})

	// simple health
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
