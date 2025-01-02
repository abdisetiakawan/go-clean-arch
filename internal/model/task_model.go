package model

import (
	"time"
)

type CreateTaskRequest struct {
	Email       string `json:"email" validate:"required,max=100"`
	Title       string `json:"title" validate:"required,max=150"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status"`
	DueDate     time.Time	`json:"due_date" validate:"required"`
}

type TaskResponse struct {
	Email 		string `json:"email,omitempty"`
	Title 		string `json:"name"`
	Description string `json:"description"`
	Status		string `json:"status"`
	DueDate		time.Time `json:"due_date"`
}

type SearchTaskRequest struct {
	Email string `json:"-"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status		string `json:"status"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}