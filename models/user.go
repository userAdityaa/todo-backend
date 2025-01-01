package models

type User struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"username"`
	Email   string `json:"email" bson:"email"`
	Picture string `json:"picture" bson:"picture"`
}
