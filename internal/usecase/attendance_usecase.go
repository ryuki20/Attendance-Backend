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
	ClockIn(ctx context.Context, userID string) (*entity.Attendance, error)
	ClockOut(ctx context.Context, userID string) (*entity.Attendance, error)
	StartBreak(ctx context.Context, userID string) (*entity.Attendance, error)
	EndBreak(ctx context.Context, userID string) (*entity.Attendance, error)
	GetTodayAttendance(ctx context.Context, userID string) (*entity.Attendance, error)
	GetAttendanceByDate(ctx context.Context, userID string, date time.Time) (*entity.Attendance, error)
	GetAttendanceHistory(ctx context.Context, userID string, startDate, endDate time.Time) ([]*entity.Attendance, error)
	GetAllAttendances(ctx context.Context, startDate, endDate time.Time) ([]*entity.Attendance, error)
	GetAttendancesByMonth(ctx context.Context, userID string, yearMonth string) ([]*entity.Attendance, error)
}

type attendanceUseCase struct {
	attendanceRepo repository.AttendanceRepository
}

func NewAttendanceUseCase(attendanceRepo repository.AttendanceRepository) AttendanceUseCase {
	return &attendanceUseCase{
		attendanceRepo: attendanceRepo,
	}
}

func (uc *attendanceUseCase) ClockIn(ctx context.Context, userID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 既存の出勤記録をチェック
	attendance, err := uc.attendanceRepo.FindByUserAndDate(ctx, userID, today)
	if err == nil && attendance.ClockIn != nil {
		return nil, fmt.Errorf("already clocked in today")
	}

	if attendance == nil {
		// 新規作成
		attendanceID := uuid.New().String()
		attendance = &entity.Attendance{
			ID:        attendanceID,
			UserID:    userID,
			Date:      today,
			ClockIn:   &now,
			Status:    entity.StatusPresent,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := uc.attendanceRepo.Create(ctx, attendance); err != nil {
			return nil, fmt.Errorf("failed to clock in: %w", err)
		}
	} else {
		// 更新
		attendance.ClockIn = &now
		attendance.Status = entity.StatusPresent
		attendance.UpdatedAt = now
		if err := uc.attendanceRepo.Update(ctx, attendance); err != nil {
			return nil, fmt.Errorf("failed to clock in: %w", err)
		}
	}

	return attendance, nil
}

func (uc *attendanceUseCase) ClockOut(ctx context.Context, userID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	attendance, err := uc.attendanceRepo.FindByUserAndDate(ctx, userID, today)
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

func (uc *attendanceUseCase) StartBreak(ctx context.Context, userID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	attendance, err := uc.attendanceRepo.FindByUserAndDate(ctx, userID, today)
	if err != nil {
		return nil, fmt.Errorf("no clock-in record found for today")
	}

	if attendance.ClockIn == nil {
		return nil, fmt.Errorf("must clock in before starting break")
	}

	if attendance.BreakStart != nil {
		return nil, fmt.Errorf("break already started")
	}

	attendance.BreakStart = &now
	attendance.UpdatedAt = now

	if err := uc.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to start break: %w", err)
	}

	return attendance, nil
}

func (uc *attendanceUseCase) EndBreak(ctx context.Context, userID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	attendance, err := uc.attendanceRepo.FindByUserAndDate(ctx, userID, today)
	if err != nil {
		return nil, fmt.Errorf("no clock-in record found for today")
	}

	if attendance.BreakStart == nil {
		return nil, fmt.Errorf("must start break before ending it")
	}

	if attendance.BreakEnd != nil {
		return nil, fmt.Errorf("break already ended")
	}

	attendance.BreakEnd = &now
	attendance.UpdatedAt = now

	if err := uc.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to end break: %w", err)
	}

	return attendance, nil
}

func (uc *attendanceUseCase) GetTodayAttendance(ctx context.Context, userID string) (*entity.Attendance, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	attendance, err := uc.attendanceRepo.FindByUserAndDate(ctx, userID, today)
	if err != nil {
		return nil, fmt.Errorf("no attendance record found for today")
	}

	return attendance, nil
}

func (uc *attendanceUseCase) GetAttendanceByDate(ctx context.Context, userID string, date time.Time) (*entity.Attendance, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	attendance, err := uc.attendanceRepo.FindByUserAndDate(ctx, userID, dateOnly)
	if err != nil {
		return nil, fmt.Errorf("no attendance record found for the date")
	}

	return attendance, nil
}

func (uc *attendanceUseCase) GetAttendanceHistory(ctx context.Context, userID string, startDate, endDate time.Time) ([]*entity.Attendance, error) {
	attendances, err := uc.attendanceRepo.ListByUser(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance history: %w", err)
	}

	return attendances, nil
}

func (uc *attendanceUseCase) GetAttendancesByMonth(ctx context.Context, userID string, yearMonth string) ([]*entity.Attendance, error) {
	t, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		return nil, fmt.Errorf("invalid year_month format")
	}
	startDate := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	attendances, err := uc.attendanceRepo.ListByUser(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	return attendances, nil
}

func (uc *attendanceUseCase) GetAllAttendances(ctx context.Context, startDate, endDate time.Time) ([]*entity.Attendance, error) {
	attendances, err := uc.attendanceRepo.ListByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get all attendances: %w", err)
	}

	return attendances, nil
}
