package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
)

var ErrEmployeeNotFound = fmt.Errorf("employee not found")

type EmployeeDetail struct {
	User        *entity.User
	Attendances []*entity.Attendance
	Total       int
	Page        int
	PerPage     int
}

type AdminUseCase interface {
	ListEmployees(ctx context.Context, page, perPage int, role *entity.UserRole) ([]*entity.User, int, error)
	GetEmployee(ctx context.Context, id, yearMonth string, page, perPage int) (*EmployeeDetail, error)
}

type adminUseCase struct {
	userRepo       repository.UserRepository
	attendanceRepo repository.AttendanceRepository
}

func NewAdminUseCase(userRepo repository.UserRepository, attendanceRepo repository.AttendanceRepository) AdminUseCase {
	return &adminUseCase{userRepo: userRepo, attendanceRepo: attendanceRepo}
}

func (uc *adminUseCase) GetEmployee(ctx context.Context, id, yearMonth string, page, perPage int) (*EmployeeDetail, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	t, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		return nil, fmt.Errorf("invalid year_month format")
	}
	startDate := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	offset := (page - 1) * perPage

	attendances, err := uc.attendanceRepo.ListByUserPaged(ctx, id, startDate, endDate, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	total, err := uc.attendanceRepo.CountByUser(ctx, id, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to count attendances: %w", err)
	}

	return &EmployeeDetail{
		User:        user,
		Attendances: attendances,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
	}, nil
}

func (uc *adminUseCase) ListEmployees(ctx context.Context, page, perPage int, role *entity.UserRole) ([]*entity.User, int, error) {
	offset := (page - 1) * perPage

	users, err := uc.userRepo.List(ctx, perPage, offset, role)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}

	total, err := uc.userRepo.Count(ctx, role)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return users, total, nil
}
