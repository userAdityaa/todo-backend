package routes

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/userAdityaa/todo-backend/internal/event"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpEventRoutes(router *chi.Mux, eventCollection *mongo.Collection, userCollection *mongo.Collection) {
	router.Get("/all-event", handlers.GetAllEvent(eventCollection, userCollection))
	router.Post("/create-event", handlers.CreateEvent(eventCollection, userCollection))
}
