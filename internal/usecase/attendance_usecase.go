package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
)

type AttendanceUseCase interface {
	ClockIn(ctx context.Context, employeeID string) (*entity.Attendance, error)
	ClockOut(ctx context.Context, employeeID string) (*entity.Attendance, error)
	GetAttendancesByMonth(ctx context.Context, employeeID string, yearMonth string) ([]*entity.Attendance, error)
}

type attendanceUseCase struct {
	attendanceRepo repository.AttendanceRepository
}

func NewAttendanceUseCase(attendanceRepo repository.AttendanceRepository) AttendanceUseCase {
	return &attendanceUseCase{
		attendanceRepo: attendanceRepo,
	}
}

func (uc *attendanceUseCase) GetAttendancesByMonth(ctx context.Context, employeeID string, yearMonth string) ([]*entity.Attendance, error) {
	t, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		return nil, fmt.Errorf("invalid year_month format")
	}
	startDate := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	attendances, err := uc.attendanceRepo.ListByEmployee(ctx, employeeID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	return attendances, nil
}

func (uc *attendanceUseCase) ClockIn(ctx context.Context, employeeID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	attendance, err := uc.attendanceRepo.FindByEmployeeAndDate(ctx, employeeID, today)
	if err == nil && attendance.ClockIn != nil {
		return nil, fmt.Errorf("already clocked in today")
	}

	if attendance == nil {
		attendanceID := uuid.New().String()
		attendance = &entity.Attendance{
			ID:         attendanceID,
			EmployeeID: employeeID,
			Date:       today,
			ClockIn:    &now,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := uc.attendanceRepo.Create(ctx, attendance); err != nil {
			return nil, fmt.Errorf("failed to clock in: %w", err)
		}
	} else {
		attendance.ClockIn = &now
		attendance.UpdatedAt = now
		if err := uc.attendanceRepo.Update(ctx, attendance); err != nil {
			return nil, fmt.Errorf("failed to clock in: %w", err)
		}
	}

	return attendance, nil
}

func (uc *attendanceUseCase) ClockOut(ctx context.Context, employeeID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	attendance, err := uc.attendanceRepo.FindByEmployeeAndDate(ctx, employeeID, today)
	if err != nil {
		return nil, fmt.Errorf("no clock-in record found for today")
	}

	if attendance.ClockIn == nil {
		return nil, fmt.Errorf("must clock in before clocking out")
	}

	if attendance.ClockOut != nil {
		return nil, fmt.Errorf("already clocked out today")
	}

	attendance.ClockOut = &now
	attendance.UpdatedAt = now

	if err := uc.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to clock out: %w", err)
	}

	return attendance, nil
}
