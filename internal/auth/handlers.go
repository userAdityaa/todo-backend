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

var user struct {
	Profile string `json:"profile"`
	Name    string `json:"name"`
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := GetGoogleAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallBackHandler(database *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code not found", http.StatusBadRequest)
			return
		}

		userInfo, err := HandleGoogleCallBack(code)
		if err != nil {
			log.Println("Failed to authenticate user:", err)
			http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
			return
		}

		var googleResponse map[string]interface{}
		if err := json.Unmarshal([]byte(userInfo), &googleResponse); err != nil {
			log.Println("Failed to parse user info:", err)
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}

		var user models.User
		user.ID = googleResponse["id"].(string)
		user.Name = googleResponse["name"].(string)
		user.Email = googleResponse["email"].(string)
		user.Picture = googleResponse["picture"].(string)

		token, err := utils.GenerateJWT(user)
		if err != nil {
			log.Println("Failed to generate Jwt: ", err)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		err = storeInDatabase(database, user)
		if err != nil {
			http.Error(w, "Error inserting user", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "google_user_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "http://localhost:3000/home", http.StatusSeeOther)
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

func GetUserDetailsHandler(database *mongo.Database, tokenString string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		filter := bson.M{"email": email}

		collection := database.Collection("user")
		var user models.User

		err = collection.FindOne(context.Background(), filter).Decode(&user)
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
