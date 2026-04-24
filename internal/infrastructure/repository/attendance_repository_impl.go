package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
	"github.com/myuto/attendance-backend/internal/infrastructure/database"
)

type attendanceRepository struct {
	db *database.DB
}

func NewAttendanceRepository(db *database.DB) repository.AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *entity.Attendance) error {
	query := `
		INSERT INTO attendances (id, user_id, date, clock_in, clock_out, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		attendance.ID, attendance.UserID, attendance.Date, attendance.ClockIn, attendance.ClockOut,
		attendance.CreatedAt, attendance.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create attendance: %w", err)
	}
	return nil
}

func (r *attendanceRepository) FindByID(ctx context.Context, id string) (*entity.Attendance, error) {
	query := `
		SELECT id, user_id, date, clock_in, clock_out, created_at, updated_at
		FROM attendances
		WHERE id = $1
	`
	attendance := &entity.Attendance{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attendance.ID, &attendance.UserID, &attendance.Date,
		&attendance.ClockIn, &attendance.ClockOut,
		&attendance.CreatedAt, &attendance.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("attendance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find attendance: %w", err)
	}
	return attendance, nil
}

func (r *attendanceRepository) FindByUserAndDate(ctx context.Context, userID string, date time.Time) (*entity.Attendance, error) {
	query := `
		SELECT id, user_id, date, clock_in, clock_out, created_at, updated_at
		FROM attendances
		WHERE user_id = $1 AND date = $2
	`
	attendance := &entity.Attendance{}
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(
		&attendance.ID, &attendance.UserID, &attendance.Date,
		&attendance.ClockIn, &attendance.ClockOut,
		&attendance.CreatedAt, &attendance.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("attendance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find attendance: %w", err)
	}
	return attendance, nil
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *entity.Attendance) error {
	query := `
		UPDATE attendances
		SET clock_in = $1, clock_out = $2, updated_at = $3
		WHERE id = $4
	`
	result, err := r.db.ExecContext(
		ctx, query,
		attendance.ClockIn, attendance.ClockOut, attendance.UpdatedAt, attendance.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("attendance not found")
	}
	return nil
}

func (r *attendanceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM attendances WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete attendance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("attendance not found")
	}
	return nil
}

func (r *attendanceRepository) ListByUser(ctx context.Context, userID string, startDate, endDate time.Time) ([]*entity.Attendance, error) {
	query := `
		SELECT id, user_id, date, clock_in, clock_out, created_at, updated_at
		FROM attendances
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
		ORDER BY date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to list attendances: %w", err)
	}
	defer rows.Close()

	return r.scanAttendances(rows)
}

func (r *attendanceRepository) ListByUserPaged(ctx context.Context, userID string, startDate, endDate time.Time, limit, offset int) ([]*entity.Attendance, error) {
	query := `
		SELECT id, user_id, date, clock_in, clock_out, created_at, updated_at
		FROM attendances
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
		ORDER BY date DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list attendances: %w", err)
	}
	defer rows.Close()

	return r.scanAttendances(rows)
}

func (r *attendanceRepository) CountByUser(ctx context.Context, userID string, startDate, endDate time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
	`
	var count int
	if err := r.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count attendances: %w", err)
	}
	return count, nil
}

func (r *attendanceRepository) scanAttendances(rows *sql.Rows) ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance
	for rows.Next() {
		attendance := &entity.Attendance{}
		err := rows.Scan(
			&attendance.ID, &attendance.UserID, &attendance.Date,
			&attendance.ClockIn, &attendance.ClockOut,
			&attendance.CreatedAt, &attendance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return attendances, nil
}
