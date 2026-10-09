package domain

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleStaff   Role = "staff"
)

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}
