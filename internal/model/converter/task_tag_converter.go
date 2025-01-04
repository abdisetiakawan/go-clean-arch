package converter

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
)

func TaskTagToResponse(taskTag *entity.TaskTag) *model.TaskTagResponse {
	return &model.TaskTagResponse{
		ID:    taskTag.ID,
		TaskId: taskTag.TaskId,
		TagId:  taskTag.TagId,
	}
}