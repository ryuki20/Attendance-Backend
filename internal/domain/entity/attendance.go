package entity

import (
	"time"
)

type Attendance struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Date      time.Time  `json:"date"`
	ClockIn   *time.Time `json:"clock_in"`
	ClockOut  *time.Time `json:"clock_out"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
