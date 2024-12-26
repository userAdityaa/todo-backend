package routes

import (
	"github.com/go-chi/chi/v5"
	handlers "github.com/userAdityaa/todo-backend/internal/todo"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpTodoRoutes(router *chi.Mux, collection *mongo.Collection) {
	router.Get("/create-todo", handlers.CreateTodo(collection))
	router.Delete("/delete-todo", handlers.DeleteTodo(collection))
	router.Post("/update-todo", handlers.UpdateTodo(collection))
}
