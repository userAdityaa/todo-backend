package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/userAdityaa/todo-backend/pkg/todo"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpTodoRoutes(router *chi.Mux, collection *mongo.Collection, userCollection *mongo.Collection) {
	router.Post("/create-todo", todo.CreateTodo(collection, userCollection))
	router.Delete("/delete-todo/{id}", todo.DeleteTodo(collection, userCollection))
	router.Put("/update-todo/{id}", todo.UpdateTodo(collection, userCollection))
	router.Get("/all-todo", todo.GetAllTodo(collection, userCollection))
}
