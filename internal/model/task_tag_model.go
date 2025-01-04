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