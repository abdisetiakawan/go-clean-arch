package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskTagRepository struct {
	Repository[entity.TaskTag]
	Log *logrus.Logger
}

func NewtaskTagRepository(log *logrus.Logger) *TaskTagRepository {
	return &TaskTagRepository{
		Log: log,
	}
}

func (r *TaskTagRepository) CreateTaskTag(db *gorm.DB, taskTag *entity.TaskTag, email string) error {
    var count int64
    err := db.Table("tasks").Where("id = ? AND email = ?", taskTag.TaskId, email).Count(&count).Error
    if err != nil {
        r.Log.WithError(err).Error("failed to validate task id and email")
        return err
    }
    if count == 0 {
        return gorm.ErrRecordNotFound
    }

    err = db.Table("tags").Where("id = ? AND email = ?", taskTag.TagId, email).Count(&count).Error
    if err != nil {
        r.Log.WithError(err).Error("failed to validate tag id and email")
        return err
    }
    if count == 0 {
        return gorm.ErrRecordNotFound
    }
    if err := db.Create(taskTag).Error; err != nil {
        r.Log.WithError(err).Error("failed to create task tag")
        return err
    }
    
    return nil
}