package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/userAdityaa/todo-backend/models"
	"github.com/userAdityaa/todo-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateEvent(eventCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event models.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			log.Println("Error decoding request body:", err)
			http.Error(w, "Invalid request payload", http.StatusNotAcceptable)
			return
		}

		event.ID = primitive.NewObjectID()

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

		if event.Title == "" || event.Start.IsZero() || event.End.IsZero() {
			http.Error(w, "Title, Start, and End are required fields", http.StatusBadRequest)
			return
		}

		_, err = eventCollection.InsertOne(context.TODO(), event)
		if err != nil {
			log.Println("Error inserting event:", err)
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}

		user.Event = append(user.Event, event)

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"event": user.Event}},
		)

		if err != nil {
			log.Println("Error updating user:", err)
			http.Error(w, "Failed to update user with new event", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message": "Event Created Successfully",
			"id":      event.ID.Hex(),
		}

		json.NewEncoder(w).Encode(response)
	}
}

func GetAllEvent(eventCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if len(user.Event) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "No events found for this user"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user.Event)
		if err != nil {
			log.Println("Error encoding events:", err)
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}
	}
}
