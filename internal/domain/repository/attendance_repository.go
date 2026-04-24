package repository

import (
	"context"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.Attendance) error
	FindByID(ctx context.Context, id string) (*entity.Attendance, error)
	FindByEmployeeAndDate(ctx context.Context, employeeID string, date time.Time) (*entity.Attendance, error)
	Update(ctx context.Context, attendance *entity.Attendance) error
	Delete(ctx context.Context, id string) error
	ListByEmployee(ctx context.Context, employeeID string, startDate, endDate time.Time) ([]*entity.Attendance, error)
	ListByEmployeePaged(ctx context.Context, employeeID string, startDate, endDate time.Time, limit, offset int) ([]*entity.Attendance, error)
	CountByEmployee(ctx context.Context, employeeID string, startDate, endDate time.Time) (int, error)
}
