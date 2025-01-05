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

func TaskWithTagsToResponse(taskWithTags *model.TaskTagResult) *model.TaskTagResult {
	return &model.TaskTagResult{
		ID:     taskWithTags.ID,
		Title:  taskWithTags.Title,
		Description: taskWithTags.Description,
		TagID:  taskWithTags.TagID,
		Status: taskWithTags.Status,
		DueDate: taskWithTags.DueDate,
	}
}