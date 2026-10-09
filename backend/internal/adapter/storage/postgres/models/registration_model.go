package models

import "time"

type RegistrationModel struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	WorkshopID    uint      `gorm:"not null;index"`
	AttendeeName  string    `gorm:"type:varchar(150);not null"`
	AttendeeEmail string    `gorm:"type:varchar(255);not null"`
	Status        string    `gorm:"type:varchar(20);not null;default:active;index"`
	RegisteredBy  uint      `gorm:"not null;index"`
	RegisteredAt  time.Time `gorm:"autoCreateTime"`
	CancelledBy   *uint     `gorm:"index"`
	CancelledAt   *time.Time

	Workshop  WorkshopModel `gorm:"foreignKey:WorkshopID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Registrar UserModel     `gorm:"foreignKey:RegisteredBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Canceller *UserModel    `gorm:"foreignKey:CancelledBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (RegistrationModel) TableName() string {
	return "registrations"
}
