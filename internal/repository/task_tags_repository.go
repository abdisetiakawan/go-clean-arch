package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
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

func (r *TaskTagRepository) SearchTaskTag(db *gorm.DB, request *model.SearchTaskTagRequest) ([]model.TaskTagResult, int64, error) {
    var count int64
    var taskTags []model.TaskTagResult
    query := db.Table("tasks").
        Select("tasks.id, tasks.title, tasks.description, tasks.status, tasks.due_date, task_tags.tag_id").
        Joins("INNER JOIN task_tags ON tasks.id = task_tags.task_id").
        Where("tasks.email = ?", request.Email)
    if err := query.Count(&count).Error; err != nil {
        r.Log.WithError(err).Error("failed to count tasks")
        return nil, 0, err
    }
    if count == 0 {
        return nil, 0, gorm.ErrRecordNotFound
    }
    if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Scan(&taskTags).Error; err != nil {
        r.Log.WithError(err).Error("failed to fetch task tags")
        return nil, 0, err
    }

    return taskTags, count, nil
}

func (r *TaskTagRepository) SearchTaskTagRequestWithTagId(db *gorm.DB, request *model.SearchTaskTagRequestWithTagId) ([]model.TaskTagResult, int64, error) {
    var count int64
    var taskTags []model.TaskTagResult
    query := db.Table("tasks").
    Select("tasks.id, tasks.title, tasks.description, tasks.status, tasks.due_date, task_tags.tag_id").
    Joins("INNER JOIN task_tags ON tasks.id = task_tags.task_id").
    Where("tasks.email = ?", request.Email).
    Where("task_tags.tag_id = ?", request.TagId)
    if err := query.Count(&count).Error; err != nil {
        r.Log.WithError(err).Error("failed to count tasks")
        return nil, 0, err
    }
    if count == 0 {
        return nil, 0, gorm.ErrRecordNotFound
    }
    if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Scan(&taskTags).Error; err != nil {
        r.Log.WithError(err).Error("failed to fetch task tags")
        return nil, 0, err
    }

    return taskTags, count, nil
}