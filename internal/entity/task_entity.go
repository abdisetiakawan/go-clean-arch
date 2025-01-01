package entity

import "time"

type Task struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      uint      `gorm:"column:user_id;not null"`
	Title       string    `gorm:"column:title;size:150;not null"`
	Description string    `gorm:"column:description;type:text"`
	Status      string    `gorm:"column:status;type:enum('pending','in_progress','completed');default:'pending'"`
	DueDate     *time.Time `gorm:"column:due_date;type:date"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
