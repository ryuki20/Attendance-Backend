package entity

import (
	"time"
)

type Attendance struct {
	ID         string           `json:"id"`
	UserID     string           `json:"user_id"`
	Date       time.Time        `json:"date"`
	ClockIn    *time.Time       `json:"clock_in"`
	ClockOut   *time.Time       `json:"clock_out"`
	BreakStart *time.Time       `json:"break_start"`
	BreakEnd   *time.Time       `json:"break_end"`
	Status     AttendanceStatus `json:"status"`
	Notes      string           `json:"notes"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type AttendanceStatus string

const (
	StatusAbsent     AttendanceStatus = "absent"
	StatusPresent    AttendanceStatus = "present"
	StatusLate       AttendanceStatus = "late"
	StatusEarlyLeave AttendanceStatus = "early_leave"
	StatusHoliday    AttendanceStatus = "holiday"
)

func (s AttendanceStatus) IsValid() bool {
	switch s {
	case StatusAbsent, StatusPresent, StatusLate, StatusEarlyLeave, StatusHoliday:
		return true
	default:
		return false
	}
}
