package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GoogleClientID string
var GoogleClientSecret string
var GoogleRedirectURL string

func LoadPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	return port
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
}

func SetUpDataBase() *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongoDB.")

	database := client.Database("todoDB")
	collection := database.Collection("todo")

	return collection
}
