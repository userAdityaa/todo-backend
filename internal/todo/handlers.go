package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/userAdityaa/todo-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllTodo(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todos []models.Todo
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var todo models.Todo
			if err := cursor.Decode(&todo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		if err := cursor.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}

func CreateTodo(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request payload", http.StatusNotAcceptable)
			return
		}

		todo.ID = primitive.NewObjectID()

		if todo.Name == "" {
			http.Error(w, "Name and DueDate are required field", http.StatusBadRequest)
			return
		}

		_, err := collection.InsertOne(context.TODO(), todo)
		if err != nil {
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message": "Todo Created Successfully",
			"id":      todo.ID.Hex(),
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		var updatedTodo models.Todo

		if err := json.Unmarshal(body, &updatedTodo); err != nil {
			http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
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

		filterID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": filterID}

		result, err := collection.UpdateOne(context.TODO(), filter, update)

		fmt.Println(err)
		if err != nil {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
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
