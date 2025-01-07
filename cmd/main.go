package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/internal/auth"
	"github.com/userAdityaa/todo-backend/routes"
)

func main() {
	config.LoadConfig()
	auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)
	database := config.SetUpDataBase()
	todoCollection := config.TodoCollection(database)
	userCollection := config.UserCollection(database)
	stickyCollection := config.StickyCollection(database)
	listCollection := config.ListCollection(database)
	eventCollection := config.EventCollection(database)
	router := chi.NewMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://minimal-planner.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(corsHandler.Handler)
	routes.SetUpTodoRoutes(router, todoCollection, userCollection)
	routes.SetUpStickyRoutes(router, stickyCollection, userCollection)
	routes.SetUpListRoutes(router, listCollection, userCollection)
	routes.SetUpEventRoutes(router, eventCollection, userCollection)
	router.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	router.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler(database))
	router.Get("/auth/user", auth.GetUserDetailsHandler(database))
	port := config.LoadPort()

	fmt.Printf("Connected Locally to port number: %s\n", port)
	if err := http.ListenAndServe(config.LoadPort(), router); err != nil {
		log.Fatal(err)
		return
	}
}
