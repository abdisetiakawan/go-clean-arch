package entity

import "time"

type User struct {
	Email       string    `gorm:"column:email;PrimaryKey;uniqueIndex"`
	Name        string    `gorm:"column:name;"`
	Password    string    `gorm:"column:password;"`
	Token       string    `gorm:"column:token;"`
	AccessToken string    `gorm:"-"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime;"`
	Task        []Task    `gorm:"foreignKey:email;references:email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "users"
}