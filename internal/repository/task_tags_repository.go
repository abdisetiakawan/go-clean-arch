package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/sirupsen/logrus"
)

type taskTagRepository struct {
	Repository[entity.TaskTag]
	Log *logrus.Logger
}

func NewtaskTagRepository(log *logrus.Logger) *taskTagRepository {
	return &taskTagRepository{
		Log: log,
	}
}
