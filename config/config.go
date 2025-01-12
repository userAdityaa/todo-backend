package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
)

func loadEnv() error {
	if os.Getenv("VERCEL") != "1" {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("error loading .env file: %v", err)
		}
	}
	return nil
}

func LoadPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return "0.0.0.0" + port
}

func LoadConfig() error {
	if err := loadEnv(); err != nil {
		return err
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
	if err := loadEnv(); err != nil {
		return nil, err
	}

	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	mongoURI := os.Getenv("MONGO_URI")
	fmt.Printf("MONGO_URI value: '%s'\n", mongoURI)
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable not set")
	}

	opts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client.Database("todoDB"), nil
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
