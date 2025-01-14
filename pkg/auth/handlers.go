package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/userAdityaa/todo-backend/models"
	"github.com/userAdityaa/todo-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := GetGoogleAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallBackHandler(database *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		userInfo, err := HandleGoogleCallBack(code)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}

		var googleResponse map[string]interface{}
		json.Unmarshal([]byte(userInfo), &googleResponse)

		user := models.User{
			ID:      googleResponse["id"].(string),
			Name:    googleResponse["name"].(string),
			Email:   googleResponse["email"].(string),
			Picture: googleResponse["picture"].(string),
		}

		if err := storeInDatabase(database, user); err != nil {
			log.Println("Database error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		token, err := utils.GenerateJWT(user)
		if err != nil {
			http.Error(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		redirectURL := "https://minimal-planner.vercel.app/home?token=" + token
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

func storeInDatabase(database *mongo.Database, user models.User) error {
	collection := database.Collection("user")

	filter := bson.M{"email": user.Email}
	var existingUser models.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		_, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func GetUserDetailsHandler(database *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Invalid token claims: email missing", http.StatusUnauthorized)
			return
		}

		collection := database.Collection("user")
		var user models.User
		filter := bson.M{"email": email}

		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
