package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/pkg/auth"
	"github.com/userAdityaa/todo-backend/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	database *mongo.Database
)

func InitializeRouter() *chi.Mux {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)

	var err error
	database, err = config.SetUpDataBase()
	if err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

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

	return router
}

func Handler(router *chi.Mux, w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
