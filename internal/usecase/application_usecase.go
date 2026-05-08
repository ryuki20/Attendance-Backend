package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
)

var (
	ErrApplicationNotFound      = errors.New("application not found")
	ErrApplicationNotPending    = errors.New("application is not in pending status")
	ErrApplicationNotOwned      = errors.New("application does not belong to the employee")
	ErrCorrectionFieldsRequired = errors.New("requested_clock_in or requested_clock_out is required")
)

type ApplicationUseCase interface {
	CreateAttendanceCorrectionApplication(ctx context.Context, employeeID string, date time.Time, requestedClockIn, requestedClockOut *time.Time, reason string) (*entity.Application, error)
	GetApplications(ctx context.Context, employeeID string, status *entity.ApplicationStatus, page, perPage int) ([]*entity.Application, int, error)
	CancelApplication(ctx context.Context, id, employeeID string) (*entity.Application, error)
	AdminListApplications(ctx context.Context, status *entity.ApplicationStatus, page, perPage int) ([]*entity.Application, int, error)
	AdminGetApplication(ctx context.Context, id string) (*entity.Application, error)
	ApproveApplication(ctx context.Context, id, approvedBy string, adminComment *string) (*entity.Application, error)
	RejectApplication(ctx context.Context, id, approvedBy string, adminComment *string) (*entity.Application, error)
}

type applicationUseCase struct {
	applicationRepo repository.ApplicationRepository
}

func NewApplicationUseCase(applicationRepo repository.ApplicationRepository) ApplicationUseCase {
	return &applicationUseCase{applicationRepo: applicationRepo}
}

func (u *applicationUseCase) CreateAttendanceCorrectionApplication(ctx context.Context, employeeID string, date time.Time, requestedClockIn, requestedClockOut *time.Time, reason string) (*entity.Application, error) {
	if requestedClockIn == nil && requestedClockOut == nil {
		return nil, ErrCorrectionFieldsRequired
	}

	now := time.Now()
	app := &entity.Application{
		ID:         uuid.New().String(),
		EmployeeID: employeeID,
		Type:       entity.ApplicationTypeAttendanceCorrection,
		Status:     entity.ApplicationStatusPending,
		Reason:     reason,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	correction := &entity.AttendanceCorrectionRequest{
		ID:                uuid.New().String(),
		ApplicationID:     app.ID,
		Date:              date,
		RequestedClockIn:  requestedClockIn,
		RequestedClockOut: requestedClockOut,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := u.applicationRepo.Create(ctx, app, correction); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return u.applicationRepo.FindByID(ctx, app.ID)
}

func (u *applicationUseCase) GetApplications(ctx context.Context, employeeID string, status *entity.ApplicationStatus, page, perPage int) ([]*entity.Application, int, error) {
	offset := (page - 1) * perPage

	apps, err := u.applicationRepo.ListByEmployeeID(ctx, employeeID, status, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list applications: %w", err)
	}

	total, err := u.applicationRepo.CountByEmployeeID(ctx, employeeID, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}

	return apps, total, nil
}

func (u *applicationUseCase) CancelApplication(ctx context.Context, id, employeeID string) (*entity.Application, error) {
	app, err := u.applicationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	if app.EmployeeID != employeeID {
		return nil, ErrApplicationNotOwned
	}
	if app.Status != entity.ApplicationStatusPending {
		return nil, ErrApplicationNotPending
	}

	return u.applicationRepo.Delete(ctx, id)
}

func (u *applicationUseCase) AdminListApplications(ctx context.Context, status *entity.ApplicationStatus, page, perPage int) ([]*entity.Application, int, error) {
	offset := (page - 1) * perPage

	apps, err := u.applicationRepo.List(ctx, status, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list applications: %w", err)
	}

	total, err := u.applicationRepo.Count(ctx, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}

	return apps, total, nil
}

func (u *applicationUseCase) AdminGetApplication(ctx context.Context, id string) (*entity.Application, error) {
	app, err := u.applicationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	return app, nil
}

func (u *applicationUseCase) ApproveApplication(ctx context.Context, id, approvedBy string, adminComment *string) (*entity.Application, error) {
	app, err := u.applicationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	if app.Status != entity.ApplicationStatusPending {
		return nil, ErrApplicationNotPending
	}

	now := time.Now()
	if err := u.applicationRepo.UpdateStatus(ctx, id, entity.ApplicationStatusApproved, approvedBy, now, adminComment); err != nil {
		return nil, fmt.Errorf("failed to approve application: %w", err)
	}

	return u.applicationRepo.FindByID(ctx, id)
}

func (u *applicationUseCase) RejectApplication(ctx context.Context, id, approvedBy string, adminComment *string) (*entity.Application, error) {
	app, err := u.applicationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrApplicationNotFound
	}
	if app.Status != entity.ApplicationStatusPending {
		return nil, ErrApplicationNotPending
	}

	now := time.Now()
	if err := u.applicationRepo.UpdateStatus(ctx, id, entity.ApplicationStatusRejected, approvedBy, now, adminComment); err != nil {
		return nil, fmt.Errorf("failed to reject application: %w", err)
	}

	return u.applicationRepo.FindByID(ctx, id)
}
