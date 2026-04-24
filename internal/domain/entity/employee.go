package entity

import "time"

type Employee struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"-"`
	Name         string       `json:"name"`
	Role         EmployeeRole `json:"role"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type EmployeeRole string

const (
	RoleAdmin    EmployeeRole = "admin"
	RoleEmployee EmployeeRole = "employee"
)

func (r EmployeeRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleEmployee:
		return true
	default:
		return false
	}
}
