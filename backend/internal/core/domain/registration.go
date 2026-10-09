package domain

import "time"

type RegistrationStatus string

const (
	RegistrationStatusPending   RegistrationStatus = "active"
	RegistrationStatusCancelled RegistrationStatus = "cancelled"
)

type Registration struct {
	ID            uint
	WorkshopID    uint
	AttendeeName  string
	AttendeeEmail string
	Status        RegistrationStatus
	RegisteredBy  string
	RegisteredAt  time.Time
	CancelledBy   string
	CancelledAt   time.Time
}
