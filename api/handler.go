package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/pkg/auth"
	"github.com/userAdityaa/todo-backend/routes"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := config.LoadConfig(); err != nil {
		log.Printf("Config error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)

	router := chi.NewRouter()
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://minimal-planner.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsHandler.Handler)
	database, err := config.SetUpDataBase()
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	router.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	router.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler(database))
	router.Get("/auth/user", auth.GetUserDetailsHandler(database))

	todoCollection := config.TodoCollection(database)
	userCollection := config.UserCollection(database)
	stickyCollection := config.StickyCollection(database)
	listCollection := config.ListCollection(database)
	eventCollection := config.EventCollection(database)

	routes.SetUpTodoRoutes(router, todoCollection, userCollection)
	routes.SetUpStickyRoutes(router, stickyCollection, userCollection)
	routes.SetUpListRoutes(router, listCollection, userCollection)
	routes.SetUpEventRoutes(router, eventCollection, userCollection)

	router.ServeHTTP(w, r)
}
