package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskRepository struct {
	Repository[entity.Task]
	Log *logrus.Logger
}

func NewTaskRepository(log *logrus.Logger) *TaskRepository {
	return &TaskRepository{
		Log: log,
	}
}

func (r *TaskRepository) Search(db *gorm.DB, request *model.SearchTaskRequest) ([]entity.Task, int64, error) {
	var tasks []entity.Task
	if err := db.Scopes(r.FilterTask(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Task{}).Scopes(r.FilterTask(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) FilterTask(request *model.SearchTaskRequest) func(tx *gorm.DB) *gorm.DB {
    return func(tx *gorm.DB) *gorm.DB {
        tx = tx.Where("email = ?", request.Email)

        if title := request.Title; title != "" {
            title = "%" + title + "%"
            tx = tx.Where("title LIKE ?", title )
        }
        if description := request.Description; description != "" {
            description = "%" + description + "%"
            tx = tx.Where("description LIKE ?", description )
        }
        if status := request.Status; status != "" {
            status = "%" + status + "%"
            tx = tx.Where("status LIKE ?", status )
        }

        return tx
    }
}

func (r *TaskRepository) FindByEmailAndId(db *gorm.DB, task *entity.Task, id string, email string) error {
	return db.Where("id = ? AND email = ?", id, email).Take(task).Error
}