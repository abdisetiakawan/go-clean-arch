package converter

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
)

func TaskToResponse(task *entity.Task) *model.TaskResponse {
	return &model.TaskResponse{
		Email: task.Email,
		Title: task.Title,
		Description: task.Description,
		Status: task.Status,
		DueDate: task.DueDate,
	}
}