package entity

import "time"

type User struct {
    Email       string    `gorm:"column:email;primaryKey;type:varchar(150);uniqueIndex"`
    Name        string    `gorm:"column:name;type:varchar(100);not null"`
    Password    string    `gorm:"column:password;type:varchar(255);not null"`
    Token       string    `gorm:"column:token;type:varchar(255)"`
    AccessToken string    `gorm:"-"`
    CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
    Tasks       []Task    `gorm:"foreignKey:email;references:email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Tags        []Tag     `gorm:"foreignKey:email;references:email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}