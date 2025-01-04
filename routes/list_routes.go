package routes

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpListRoutes(router *chi.Mux, listCollection *mongo.Collection, userCollection *mongo.Collection) {
	// router.Post("/create-list", handlers.)
}
