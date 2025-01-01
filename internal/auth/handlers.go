package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/userAdityaa/todo-backend/models"
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

		err = storeInDatabase(database, user)
		if err != nil {
			http.Error(w, "Error inserting user", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "http://localhost:3000", http.StatusSeeOther)
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
