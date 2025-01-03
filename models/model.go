package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	List        string             `json:"list" bson:"list"`
	DueDate     string             `json:"due_date" bson:"due_date"`
	Subtask     []string           `json:"sub_task" bson:"sub_task"`
}

type User struct {
	ID      string   `json:"id" bson:"_id"`
	Name    string   `json:"name" bson:"username"`
	Email   string   `json:"email" bson:"email"`
	Picture string   `json:"picture" bson:"picture"`
	Todo    []Todo   `json:"todos" bson:"todos"`
	Stick   []Sticky `json:"sticky" bson:"sticky"`
}

type Sticky struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Topic   string             `json:"topic" bson:"topic"`
	Content string             `json:"content" bson:"content"`
	Color   string             `json:"color" bson:"color"`
}
