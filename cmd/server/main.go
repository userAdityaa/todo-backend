package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/routes"
	// "github.com/userAdityaa/todo-backend/routes"
)

func main() {
	collection := config.SetUpDataBase()
	router := chi.NewMux()

	routes.SetUpTodoRoutes(router, collection)

	port := config.LoadPort()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All Set."))
	})

	fmt.Printf("Connected Locally to port number: %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
		return
	}

	// config.LoadConfig()

	// auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)

	// r := chi.NewRouter()
	// corsHandler := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	// 	AllowCredentials: true,
	// }).Handler

	// r.Use(corsHandler)

	// r.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	// r.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler)

	// log.Println("Server running at http://localhost:8000")
	// log.Fatal(http.ListenAndServe(":8000", r))
}
