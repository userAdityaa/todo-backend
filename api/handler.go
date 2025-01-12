package handler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/userAdityaa/todo-backend/config"
	"github.com/userAdityaa/todo-backend/pkg/auth"
	"github.com/userAdityaa/todo-backend/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	router     *chi.Mux
	database   *mongo.Database
	setupOnce  sync.Once
	setupError error
	client     *mongo.Client
)

// Initialize connection pool settings
const (
	maxPoolSize     uint64 = 5
	minPoolSize     uint64 = 1
	maxConnIdleTime        = 30 * time.Second
)

func init() {
	setupOnce.Do(func() {
		if err := config.LoadConfig(); err != nil {
			setupError = err
			return
		}

		// Initialize MongoDB client with optimized connection pool
		clientOptions := options.Client().
			SetMaxPoolSize(maxPoolSize).
			SetMinPoolSize(minPoolSize).
			SetMaxConnIdleTime(maxConnIdleTime)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			setupError = err
			return
		}

		// Ping database to verify connection
		if err = client.Ping(ctx, nil); err != nil {
			setupError = err
			return
		}

		database = client.Database(config.GetConfig().DBName)
		auth.InitGoogleOAuth(config.GoogleClientID, config.GoogleClientSecret, config.GoogleRedirectURL)

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

		// Pre-initialize collections
		collections := initializeCollections(database)

		// Set up routes with dependency injection
		setupRoutes(router, collections, database)
	})
}

type Collections struct {
	Todo   *mongo.Collection
	User   *mongo.Collection
	Sticky *mongo.Collection
	List   *mongo.Collection
	Event  *mongo.Collection
}

func initializeCollections(db *mongo.Database) Collections {
	return Collections{
		Todo:   config.TodoCollection(db),
		User:   config.UserCollection(db),
		Sticky: config.StickyCollection(db),
		List:   config.ListCollection(db),
		Event:  config.EventCollection(db),
	}
}

func setupRoutes(r *chi.Mux, collections Collections, db *mongo.Database) {
	// Auth routes
	r.HandleFunc("/auth/google/login", auth.GoogleLoginHandler)
	r.HandleFunc("/auth/google/callback", auth.GoogleCallBackHandler(db))
	r.Get("/auth/user", auth.GetUserDetailsHandler(db))

	// Application routes
	routes.SetUpTodoRoutes(r, collections.Todo, collections.User)
	routes.SetUpStickyRoutes(r, collections.Sticky, collections.User)
	routes.SetUpListRoutes(r, collections.List, collections.User)
	routes.SetUpEventRoutes(r, collections.Event, collections.User)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if setupError != nil {
		log.Printf("Setup error: %v", setupError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	router.ServeHTTP(w, r)
}

// Cleanup function to be called when the application shuts down
func Cleanup() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}
