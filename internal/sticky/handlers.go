package sticky

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

func CreateSticky(stickyCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sticky models.Sticky
		if err := json.NewDecoder(r.Body).Decode(&sticky); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request payload", http.StatusNotAcceptable)
			return
		}

		sticky.ID = primitive.NewObjectID()

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

		if sticky.Topic == "" || sticky.Content == "" || sticky.Color == "" {
			http.Error(w, "Topic, Content, and Color are required fields", http.StatusBadRequest)
			return
		}

		_, err = stickyCollection.InsertOne(context.TODO(), sticky)
		if err != nil {
			log.Println("Error inserting sticky:", err)
			http.Error(w, "Failed to create Sticky", http.StatusInternalServerError)
			return
		}

		insertedID := sticky.ID
		sticky.ID = insertedID
		user.Stick = append(user.Stick, sticky)

		_, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"sticky": user.Stick}},
		)

		if err != nil {
			log.Println("Error updating user:", err)
			http.Error(w, "Failed to update user with new Sticky", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message": "Sticky Created Successfully",
			"id":      sticky.ID.Hex(),
		}

		json.NewEncoder(w).Encode(response)
	}
}

func GetAllSticky(stickyCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
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

		if len(user.Stick) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "No Sticky found for this user"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user.Stick)
		if err != nil {
			log.Println("Error encoding sticky: ", err)
			http.Error(w, "Failed to fetch sticky", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateSticky(stickyCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
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

		var partialUpdate struct {
			ID      primitive.ObjectID `json:"id"`
			Topic   *string            `json:"topic,omitempty"`
			Content *string            `json:"content,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&partialUpdate); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if partialUpdate.ID.IsZero() {
			http.Error(w, "Invalid sticky ID", http.StatusBadRequest)
			return
		}

		updateFields := bson.M{}
		if partialUpdate.Topic != nil {
			updateFields["sticky.$.topic"] = *partialUpdate.Topic
		}
		if partialUpdate.Content != nil {
			updateFields["sticky.$.content"] = *partialUpdate.Content
		}

		if len(updateFields) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		result, err := stickyCollection.UpdateOne(
			context.TODO(),
			bson.M{
				"_id": partialUpdate.ID,
			},
			bson.M{
				"$set": updateFields,
			},
		)

		result, err = userCollection.UpdateOne(
			context.TODO(),
			bson.M{
				"email":      email,
				"sticky._id": partialUpdate.ID,
			},
			bson.M{
				"$set": updateFields,
			},
		)

		if err != nil {
			log.Printf("Error updating sticky: %v", err)
			http.Error(w, "Error updating sticky", http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "Sticky not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sticky updated successfully",
			"sticky":  partialUpdate,
		})
	}
}

func DeleteSticky(stickyCollection *mongo.Collection, userCollection *mongo.Collection) http.HandlerFunc {
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
			http.Error(w, "Invalid sticky ID", http.StatusBadRequest)
			return
		}

		result, err := stickyCollection.DeleteOne(
			context.TODO(),
			bson.M{"_id": deleteRequest.ID},
		)
		if err != nil {
			log.Printf("Error deleting sticky: %v", err)
			http.Error(w, "Error deleting sticky", http.StatusInternalServerError)
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
				"$pull": bson.M{"sticky": bson.M{"_id": deleteRequest.ID}},
			},
		)
		if err != nil {
			log.Printf("Error updating user collection: %v", err)
			http.Error(w, "Error removing sticky from user data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sticky deleted successfully",
		})
	}
}
