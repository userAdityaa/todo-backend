package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/userAdityaa/todo-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTodo(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		log.Println(json.NewDecoder(r.Body))
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request payload", http.StatusNotAcceptable)
			return
		}

		if todo.Name == "" {
			http.Error(w, "Name and DueDate are required field", http.StatusBadRequest)
			return
		}

		result, err := collection.InsertOne(context.TODO(), todo)
		if err != nil {
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message":     "Todo Created Successfully",
			"inserted_id": result.InsertedID,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func DeleteTodo(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			log.Fatal(err)
			return
		}

		_, err := collection.DeleteOne(context.TODO(), todo)
		if err != nil {
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message": "Task Deleted Successfully",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateTodo(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "Invalid Todo ID", http.StatusBadRequest)
			return
		}

		var updatedTodo models.Todo

		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		todoId, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"name":        updatedTodo.Name,
				"description": updatedTodo.Description,
				"list":        updatedTodo.List,
				"due_date":    updatedTodo.DueDate,
				"sub_task":    updatedTodo.Subtask,
			},
		}

		filter := bson.M{"_id": todoId}

		result, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"message": "Todo updated successfully",
			"matched": result.MatchedCount,
			"updated": result.ModifiedCount,
		}
		json.NewEncoder(w).Encode(response)
	}
}
