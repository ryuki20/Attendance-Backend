package repository

import (
	"context"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.Attendance) error
	FindByID(ctx context.Context, id int) (*entity.Attendance, error)
	FindByUserAndDate(ctx context.Context, userID int, date time.Time) (*entity.Attendance, error)
	Update(ctx context.Context, attendance *entity.Attendance) error
	Delete(ctx context.Context, id int) error
	ListByUser(ctx context.Context, userID int, startDate, endDate time.Time) ([]*entity.Attendance, error)
	ListByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.Attendance, error)
}
