package repository

import (
	"github.com/abdisetiakawan/go-clean-arch/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByToken (db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}