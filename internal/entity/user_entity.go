package entity

import "time"

type User struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;"`
	Email       string    `gorm:"column:email;"`
	Password    string    `gorm:"column:password;"`
	Token       string    `gorm:"column:token;"`
	AccessToken string    `gorm:"-"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime;"`
	Task        []Task    `gorm:"foreignKey:user_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "users"
}