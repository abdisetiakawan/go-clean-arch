package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id int) error {
	return db.Where("id = ?", id).First(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var count int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count, err
}

func (r *Repository[T]) CountByEmail(db *gorm.DB, email string) (int64, error) {
	var count int64
	err := db.Model(new(T)).Where("email = ?", email).Count(&count).Error
	return count, err
}

func (r *Repository[T]) FindByEmail(db *gorm.DB, entity *T, email string) error {
	return db.Where("email = ?", email).First(entity).Error
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}