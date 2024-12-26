package models

import (
	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	List        []string  `json:"list" bson:"list"`
	DueDate     string    `json:"due_date" bson:"due_date"`
	Subtask     []string  `json:"sub_task" bson:"sub_task"`
}
