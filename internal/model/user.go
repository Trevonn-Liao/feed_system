package model

import "time"

type User struct {
	ID           int64     `json:"id" gorm:"primaryKey;column:id"`
	Username     string    `json:"username" gorm:"column:username;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
