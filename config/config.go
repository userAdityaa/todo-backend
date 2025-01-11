package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GoogleClientID string
var GoogleClientSecret string
var GoogleRedirectURL string

func LoadPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return "0.0.0.0" + port
}

func LoadConfig() error {
	if os.Getenv("VERCEL") == "" {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("error loading .env file: %v", err)
		}
	}

	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")

	if GoogleClientID == "" || GoogleClientSecret == "" || GoogleRedirectURL == "" {
		return fmt.Errorf("missing required environment variables")
	}

	return nil
}

func SetUpDataBase() (*mongo.Database, error) {
	// Only load .env in development
	if os.Getenv("VERCEL") == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println(mongoURI)
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable not set")
	}

	opts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	database := client.Database("todoDB")
	return database, nil
}

func TodoCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("todo")
}

func UserCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("user")
}

func StickyCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("sticky")
}

func ListCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("list")
}

func EventCollection(database *mongo.Database) *mongo.Collection {
	return database.Collection("event")
}
