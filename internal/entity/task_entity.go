package entity

import "time"

type Task struct {
    ID          uint      `gorm:"column:id;primaryKey;autoIncrement"`
    Email       string    `gorm:"column:email;type:varchar(100);not null;index"`
    Title       string    `gorm:"column:title;type:varchar(150);not null"`
    Description string    `gorm:"column:description;type:text"`
    Status      string    `gorm:"column:status;type:enum('pending','in_progress','completed');default:pending"`
    DueDate     time.Time `gorm:"column:due_date;type:date"`
    CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
    Tags        []Tag     `gorm:"many2many:task_tags"`
    User        User      `gorm:"foreignKey:email;references:email"`
}

func (Task) TableName() string {
	return "tasks"
}