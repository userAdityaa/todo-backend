package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/userAdityaa/todo-backend/models"
	"github.com/userAdityaa/todo-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateList(listCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
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

		var newList models.List
		if err := json.NewDecoder(r.Body).Decode(&newList); err != nil {
			http.Error(w, "Error creating new List", http.StatusBadRequest)
			return
		}

		newList.ID = primitive.NewObjectID()

		if newList.Color == "" || newList.Name == "" {
			http.Error(w, "Color or Name for the list is missing", http.StatusBadRequest)
			return
		}

		_, err = listCollection.InsertOne(context.TODO(), newList)
		if err != nil {
			http.Error(w, "Failed to create list", http.StatusInternalServerError)
			return
		}

		insertedID := newList.ID
		newList.ID = insertedID

		user.List = append(user.List, newList)

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"list": user.List}},
		)

		if err != nil {
			http.Error(w, "Failed to update user list", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message": "List Created Successfully",
			"id":      newList.ID.Hex(),
		}

		json.NewEncoder(w).Encode(response)
	}
}

func DeleteList(listCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
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

		var deleteRequest struct {
			ID primitive.ObjectID `json:"id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if deleteRequest.ID.IsZero() {
			http.Error(w, "Invalide list ID", http.StatusBadRequest)
			return
		}

		result, err := listCollection.DeleteOne(
			context.TODO(),
			bson.M{"_id": deleteRequest.ID},
		)
		if err != nil {
			http.Error(w, "Error Deleting List", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			http.Error(w, "Sticky not found", http.StatusNotFound)
			return
		}

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"email": email},
			bson.M{
				"$pull": bson.M{"list": bson.M{"_id": deleteRequest.ID}},
			},
		)

		if err != nil {
			http.Error(w, "Error removing list", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "List Deleted Successfully",
		})
	}
}

func GetAllList(listCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
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

		if len(user.List) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "No list found for this user"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user.List)
		if err != nil {
			http.Error(w, "Failed to fetch List", http.StatusInternalServerError)
			return
		}
	}
}

func FindAList(listCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		listID := chi.URLParam(r, "id")
		if listID == "" {
			http.Error(w, "List ID is required", http.StatusBadRequest)
			return
		}

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

		var foundList models.List
		found := false
		for _, list := range user.List {
			if list.ID.Hex() == listID {
				foundList = list
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "List not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(foundList)
		if err != nil {
			http.Error(w, "Failed to encode list", http.StatusInternalServerError)
			return
		}
	}
}
