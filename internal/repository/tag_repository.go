package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/abdisetiakawan/go-clean-arch/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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


func (r *TagRepository) Search(db *gorm.DB, request *model.SearchTagRequest) ([]entity.Tag, int64, error) {
	var tags []entity.Tag
	if err := db.Scopes(r.FilterTag(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Tag{}).Scopes(r.FilterTag(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

func (r *TagRepository) FilterTag(request *model.SearchTagRequest) func(tx *gorm.DB) *gorm.DB {
    return func(tx *gorm.DB) *gorm.DB {
        tx = tx.Where("email = ?", request.Email)

        if name := request.Name; name != "" {
            name = "%" + name + "%"
            tx = tx.Where("name LIKE ?", name )
        }
        return tx
    }
}

func (r *TagRepository) FindByEmailAndId(db *gorm.DB, tag *entity.Tag, id string, email string) error {
	return db.Where("id = ? AND email = ?", id, email).Take(tag).Error
}