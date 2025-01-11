package routes

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/userAdityaa/todo-backend/pkg/container"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpListRoutes(router *chi.Mux, listCollection *mongo.Collection, userCollection *mongo.Collection) {
	router.Post("/create-list", handlers.CreateList(listCollection, userCollection))
	router.Delete("/delete-list", handlers.DeleteList(listCollection, userCollection))
	router.Get("/all-list", handlers.GetAllList(listCollection, userCollection))
	router.Get("/lists/{id}", handlers.FindAList(listCollection, userCollection))
}
