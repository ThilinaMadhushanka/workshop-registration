package models

import "time"

type UserModel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"type:varchar(150);not null"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	Role         string    `gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (UserModel) TableName() string {
	return "users"
}
