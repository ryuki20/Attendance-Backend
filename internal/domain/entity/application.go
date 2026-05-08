package entity

import "time"

type ApplicationType string
type ApplicationStatus string

const (
	ApplicationTypeAttendanceCorrection ApplicationType = "ATTENDANCE_CORRECTION"
)

const (
	ApplicationStatusPending  ApplicationStatus = "PENDING"
	ApplicationStatusApproved ApplicationStatus = "APPROVED"
	ApplicationStatusRejected ApplicationStatus = "REJECTED"
)

type EmployeeSummary struct {
	ID   string
	Name string
}

type AttendanceCorrectionRequest struct {
	ID                string
	ApplicationID     string
	Date              time.Time
	RequestedClockIn  *time.Time
	RequestedClockOut *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Application struct {
	ID           string
	EmployeeID   string
	Type         ApplicationType
	Status       ApplicationStatus
	Reason       string
	ApprovedBy   *string
	ApprovedAt   *time.Time
	AdminComment *string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// JOINで取得する関連データ
	Employee   *EmployeeSummary
	Approver   *EmployeeSummary
	Correction *AttendanceCorrectionRequest
}
