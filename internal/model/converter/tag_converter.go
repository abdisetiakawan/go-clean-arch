package converter

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
)

func TagToResponse(tag *entity.Tag) *model.TagResponse {
	return &model.TagResponse{
		ID: tag.ID,
		Email: tag.Email,
		Name: tag.Name,
	}
}