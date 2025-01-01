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
	collection := config.SetUpDataBase()
	router := chi.NewMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(corsHandler.Handler)
	routes.SetUpTodoRoutes(router, collection)
	router.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	router.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler)
	port := config.LoadPort()

	fmt.Printf("Connected Locally to port number: %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
		return
	}
}
