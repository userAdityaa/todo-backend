package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/pkg/auth"
	"github.com/userAdityaa/todo-backend/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	router     *chi.Mux
	database   *mongo.Database
	setupOnce  sync.Once
	setupError error
)

func init() {
	setupOnce.Do(func() {
		// Load configuration
		if err := config.LoadConfig(); err != nil {
			setupError = err
			return
		}

		// Initialize Google OAuth
		auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)

		// Set up database connection
		db, err := config.SetUpDataBase()
		if err != nil {
			setupError = err
			return
		}
		database = db

		// Initialize router with CORS
		router = chi.NewMux()
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   []string{"https://minimal-planner.vercel.app"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		router.Use(corsHandler.Handler)

		// Set up routes
		router.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
		router.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler(database))
		router.Get("/auth/user", auth.GetUserDetailsHandler(database))

		// Initialize collections
		todoCollection := config.TodoCollection(database)
		userCollection := config.UserCollection(database)
		stickyCollection := config.StickyCollection(database)
		listCollection := config.ListCollection(database)
		eventCollection := config.EventCollection(database)

		// Set up route handlers
		routes.SetUpTodoRoutes(router, todoCollection, userCollection)
		routes.SetUpStickyRoutes(router, stickyCollection, userCollection)
		routes.SetUpListRoutes(router, listCollection, userCollection)
		routes.SetUpEventRoutes(router, eventCollection, userCollection)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if setupError != nil {
		log.Printf("Setup error: %v", setupError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)
}
