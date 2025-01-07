package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/internal/auth"
	"github.com/userAdityaa/todo-backend/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler is the entrypoint for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize configurations only once
	initializeOnce.Do(func() {
		config.LoadConfig()
		auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)
		database = config.SetUpDataBase() // Store in package-level variable
	})

	// Create a new router instance for each request
	router := chi.NewRouter()

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://minimal-planner.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(corsHandler.Handler)

	// Set up routes using the shared database connection
	routes.SetUpTodoRoutes(router, config.TodoCollection(database), config.UserCollection(database))
	routes.SetUpStickyRoutes(router, config.StickyCollection(database), config.UserCollection(database))
	routes.SetUpListRoutes(router, config.ListCollection(database), config.UserCollection(database))
	routes.SetUpEventRoutes(router, config.EventCollection(database), config.UserCollection(database))

	// Auth routes
	router.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	router.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler(database))
	router.Get("/auth/user", auth.GetUserDetailsHandler(database))

	// Serve the request
	router.ServeHTTP(w, r)
}

// Package-level variables for maintaining state between invocations
var (
	initializeOnce sync.Once
	database       *mongo.Database
)
