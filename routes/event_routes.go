package routes

import (
	"github.com/go-chi/chi/v5"
	Event "github.com/userAdityaa/todo-backend/pkg/event"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpEventRoutes(router *chi.Mux, eventCollection *mongo.Collection, userCollection *mongo.Collection) {
	router.Get("/all-event", Event.GetAllEvent(eventCollection, userCollection))
	router.Post("/create-event", Event.CreateEvent(eventCollection, userCollection))
}
