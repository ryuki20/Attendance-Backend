package repository

import (
	"context"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *entity.Application, correction *entity.AttendanceCorrectionRequest) error
	FindByID(ctx context.Context, id string) (*entity.Application, error)
	ListByEmployeeID(ctx context.Context, employeeID string, status *entity.ApplicationStatus, limit, offset int) ([]*entity.Application, error)
	CountByEmployeeID(ctx context.Context, employeeID string, status *entity.ApplicationStatus) (int, error)
	List(ctx context.Context, status *entity.ApplicationStatus, limit, offset int) ([]*entity.Application, error)
	Count(ctx context.Context, status *entity.ApplicationStatus) (int, error)
	UpdateStatus(ctx context.Context, id string, status entity.ApplicationStatus, approvedBy string, approvedAt time.Time, adminComment *string) error
	Delete(ctx context.Context, id string) (*entity.Application, error)
}
