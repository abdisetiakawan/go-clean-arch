package entity

import "time"

type Tag struct {
    ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
    Email     string    `gorm:"column:email;type:varchar(150);index"`
    Name      string    `gorm:"column:name;type:varchar(50);not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
    Tasks     []Task    `gorm:"many2many:task_tags"`
    User      User      `gorm:"foreignKey:email;references:email"`
}

func (Tag) TableName() string {
	return "tags"
}