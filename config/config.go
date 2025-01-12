package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration constants
const (
	defaultPort     = ":8080"
	defaultDBName   = "todoDB"
	connectTimeout  = 10 * time.Second
	maxPoolSize     = 5
	minPoolSize     = 1
	maxConnIdleTime = 30 * time.Second
)

// Config holds all configuration variables
type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	MongoURI           string
	Port               string
	DBName             string
}

var (
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	appConfig          *Config
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
		port = defaultPort
	}
	return "0.0.0.0" + port
}

func LoadConfig() error {
	if err := loadEnv(); err != nil {
		return err
	}

	appConfig = &Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		MongoURI:           os.Getenv("MONGO_URI"),
		Port:               LoadPort(),
		DBName:             getDBName(),
	}

	// Set global variables for backward compatibility
	GoogleClientID = appConfig.GoogleClientID
	GoogleClientSecret = appConfig.GoogleClientSecret
	GoogleRedirectURL = appConfig.GoogleRedirectURL

	if err := validateConfig(appConfig); err != nil {
		return err
	}

	return nil
}

func getDBName() string {
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		return defaultDBName
	}
	return dbName
}

func validateConfig(cfg *Config) error {
	if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" || cfg.GoogleRedirectURL == "" {
		return fmt.Errorf("missing required Google OAuth environment variables")
	}
	if cfg.MongoURI == "" {
		return fmt.Errorf("MONGO_URI environment variable not set")
	}
	return nil
}

func SetUpDataBase() (*mongo.Database, error) {
	if appConfig == nil {
		return nil, fmt.Errorf("config not initialized, call LoadConfig first")
	}

	// Configure MongoDB client options
	opts := options.Client().
		ApplyURI(appConfig.MongoURI).
		SetMaxPoolSize(maxPoolSize).
		SetMinPoolSize(minPoolSize).
		SetMaxConnIdleTime(maxConnIdleTime)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client.Database(appConfig.DBName), nil
}

// Collection getters with proper error handling
func getCollection(database *mongo.Database, name string) *mongo.Collection {
	if database == nil {
		return nil
	}
	return database.Collection(name)
}

func TodoCollection(database *mongo.Database) *mongo.Collection {
	return getCollection(database, "todo")
}

func UserCollection(database *mongo.Database) *mongo.Collection {
	return getCollection(database, "user")
}

func StickyCollection(database *mongo.Database) *mongo.Collection {
	return getCollection(database, "sticky")
}

func ListCollection(database *mongo.Database) *mongo.Collection {
	return getCollection(database, "list")
}

func EventCollection(database *mongo.Database) *mongo.Collection {
	return getCollection(database, "event")
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return appConfig
}
