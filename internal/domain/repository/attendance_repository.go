package repository

import (
	"context"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.Attendance) error
	FindByID(ctx context.Context, id string) (*entity.Attendance, error)
	FindByUserAndDate(ctx context.Context, userID string, date time.Time) (*entity.Attendance, error)
	Update(ctx context.Context, attendance *entity.Attendance) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string, startDate, endDate time.Time) ([]*entity.Attendance, error)
	ListByUserPaged(ctx context.Context, userID string, startDate, endDate time.Time, limit, offset int) ([]*entity.Attendance, error)
	CountByUser(ctx context.Context, userID string, startDate, endDate time.Time) (int, error)
}
