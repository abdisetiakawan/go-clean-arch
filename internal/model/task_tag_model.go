package model

type CreateTaskTagRequest struct {
	TaskId uint `json:"-" validate:"required"`
	TagId  uint `json:"tag_id" validate:"required"`
}

type TaskTagResponse struct {
	ID     uint `json:"id"`
	TaskId uint `json:"taskId"`
	TagId  uint `json:"tag_id"`
}

type GetTaskTagRequest struct {
	TaskId uint `json:"-" validate:"required"`
	TagId  uint `json:"-" validate:"required"`
}

type TaskTagToResponse struct {
	ID     uint `json:"id"`
	TaskId uint `json:"task_id"`
	TagId  uint `json:"tag_id"`
}

type SearchTaskTagRequest struct {
	Email string `json:"-" validate:"required"`
	Page  int    `json:"page" validate:"min=1"`
	Size  int    `json:"size" validate:"min=1,max=100"`
}

type SearchTaskTagRequestWithTagId struct {
	Email string `json:"-" validate:"required"`
	TagId uint   `json:"-" validate:"required"`
	Page  int    `json:"page" validate:"min=1"`
	Size  int    `json:"size" validate:"min=1,max=100"`
}

type TaskTagResult struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
	TagID       uint   `json:"tag_id"`
}

type GetTaskTagForDelete struct {
	Email  string `json:"-" validate:"required"`
	TaskId uint   `json:"-" validate:"required"`
	TagId  uint   `json:"-" validate:"required"`
}