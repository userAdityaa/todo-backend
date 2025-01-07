package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/internal/auth"
	"github.com/userAdityaa/todo-backend/routes"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	config.LoadConfig()
	auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)
	database := config.SetUpDataBase()

	router := chi.NewMux()

	// CORS Configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://minimal-planner.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsHandler.Handler)

	// Routes
	routes.SetUpTodoRoutes(router, config.TodoCollection(database), config.UserCollection(database))
	routes.SetUpStickyRoutes(router, config.StickyCollection(database), config.UserCollection(database))
	routes.SetUpListRoutes(router, config.ListCollection(database), config.UserCollection(database))
	routes.SetUpEventRoutes(router, config.EventCollection(database), config.UserCollection(database))
	router.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	router.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler(database))
	router.Get("/auth/user", auth.GetUserDetailsHandler(database))

	// Serve the request
	router.ServeHTTP(w, r)
}
