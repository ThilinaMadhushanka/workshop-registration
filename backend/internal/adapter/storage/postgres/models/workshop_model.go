package models

import "time"

type WorkshopModel struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Code       string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Title      string    `gorm:"type:varchar(255);not null"`
	Instructor string    `gorm:"type:varchar(150);not null"`
	StartsAt   time.Time `gorm:"not null;index"`
	Capacity   int       `gorm:"not null;check:capacity_positive,capacity > 0"`
	Status     string    `gorm:"type:varchar(50);not null;default:'scheduled';index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (WorkshopModel) TableName() string {
	return "workshops"
}
