package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/sirupsen/logrus"
)

type TagRepository struct {
	Repository[entity.Tag]
	Log *logrus.Logger
}

func NewTagRepository(log *logrus.Logger) *TagRepository {
	return &TagRepository{
		Log: log,
	}
}
