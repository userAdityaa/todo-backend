package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Color string             `json:"color" bson:"color"`
}

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
	List    []List   `json:"list" bson:"list"`
	Event   []Event  `json:"event" bson:"event"`
}

type Sticky struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Topic   string             `json:"topic" bson:"topic"`
	Content string             `json:"content" bson:"content"`
	Color   string             `json:"color" bson:"color"`
}

type Event struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Title string             `json:"title" bson:"title"`
	Date  time.Time          `json:"date" bson:"date"`
	Color string             `json:"color" bson:"color"`
	Start time.Time          `json:"start" bson:"start"`
	End   time.Time          `json:"end" bson:"end"`
}
