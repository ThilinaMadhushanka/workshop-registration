package domain

import "time"

type WorkshopStatus string

const (
	WorkshopStatusScheduled WorkshopStatus = "scheduled"
	WorkshopStatusCompleted WorkshopStatus = "completed"
	WorkshopStatusCancelled WorkshopStatus = "cancelled"
)

type Workshop struct {
	ID         uint
	Code       string
	Title      string
	Instructor string
	StartsAt   time.Time
	Capacity   int
	Status     WorkshopStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
