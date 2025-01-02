package routes

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/userAdityaa/todo-backend/internal/todo"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpTodoRoutes(router *chi.Mux, collection *mongo.Collection, userCollection *mongo.Collection) {
	router.Post("/create-todo", handlers.CreateTodo(collection, userCollection))
	router.Delete("/delete-todo/{id}", handlers.DeleteTodo(collection))
	router.Put("/update-todo/{id}", handlers.UpdateTodo(collection))
	router.Get("/all-todo", handlers.GetAllTodo(collection))
}
