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
	"github.com/userAdityaa/todo-backend/utils"
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

func CreateTodo(collection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request payload", http.StatusNotAcceptable)
			return
		}

		todo.ID = primitive.NewObjectID()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Invalid token claims: email missing", http.StatusUnauthorized)
			return
		}

		var user models.User
		err = userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
		if err != nil {
			log.Println("User not found:", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if todo.Name == "" {
			http.Error(w, "Name and DueDate are required field", http.StatusBadRequest)
			return
		}

		_, err = collection.InsertOne(context.TODO(), todo)
		if err != nil {
			log.Fatal(err)
			return
		}

		user.Todo = append(user.Todo, todo)

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"todos": user.Todo}},
		)

		if err != nil {
			log.Println("Error updating user:", err)
			http.Error(w, "Failed to update user with new Todo", http.StatusInternalServerError)
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
		id := chi.URLParam(r, "id")
		filterID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			log.Println("Invalid ID format:", err)
			return
		}

		filter := bson.M{"_id": filterID}

		result, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			log.Println("Error deleting task:", err)
			return
		}

		if result.DeletedCount == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"message": "Task Deleted Successfully",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error encoding response:", err)
		}
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

// func GetATodo(collection *mongo.Collection) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
