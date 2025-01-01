package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/userAdityaa/todo-backend/models"
)

var jwtSecret = []byte("your-secret-key")

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
