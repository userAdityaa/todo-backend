package routes

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/userAdityaa/todo-backend/internal/sticky"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpStickyRoutes(router *chi.Mux, stickCollection *mongo.Collection, userCollection *mongo.Collection) {
	router.Post("/create-sticky", handlers.CreateSticky(stickCollection, userCollection))
	router.Get("/all-sticky", handlers.GetAllSticky(stickCollection, userCollection))
	router.Put("/update-sticky", handlers.UpdateSticky(stickCollection, userCollection))
}
