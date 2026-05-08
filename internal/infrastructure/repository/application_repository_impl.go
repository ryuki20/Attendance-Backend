package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	domainrepo "github.com/myuto/attendance-backend/internal/domain/repository"
	"github.com/myuto/attendance-backend/internal/infrastructure/database"
)

type applicationRepository struct {
	db *database.DB
}

func NewApplicationRepository(db *database.DB) domainrepo.ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *entity.Application, correction *entity.AttendanceCorrectionRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO applications (id, employee_id, type, status, reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, app.ID, app.EmployeeID, app.Type, app.Status, app.Reason, app.CreatedAt, app.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO attendance_correction_requests (id, application_id, date, requested_clock_in, requested_clock_out, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, correction.ID, correction.ApplicationID, correction.Date, correction.RequestedClockIn, correction.RequestedClockOut, correction.CreatedAt, correction.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create attendance correction request: %w", err)
	}

	return tx.Commit()
}

const applicationSelectQuery = `
	SELECT
		a.id, a.employee_id, a.type, a.status, a.reason,
		a.approved_by, a.approved_at, a.admin_comment,
		a.created_at, a.updated_at,
		e.name AS employee_name,
		approver.id AS approver_id, approver.name AS approver_name,
		acr.id AS correction_id, acr.date, acr.requested_clock_in, acr.requested_clock_out
	FROM applications a
	JOIN employees e ON a.employee_id = e.id
	LEFT JOIN employees approver ON a.approved_by = approver.id
	LEFT JOIN attendance_correction_requests acr ON acr.application_id = a.id
`

func (r *applicationRepository) FindByID(ctx context.Context, id string) (*entity.Application, error) {
	query := applicationSelectQuery + `WHERE a.id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	app, err := scanApplication(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find application: %w", err)
	}
	return app, nil
}

func (r *applicationRepository) ListByEmployeeID(ctx context.Context, employeeID string, status *entity.ApplicationStatus, limit, offset int) ([]*entity.Application, error) {
	where, args := buildWhereClause([]string{"a.employee_id = $1"}, []any{employeeID}, status)
	query := applicationSelectQuery + where + fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	defer rows.Close()
	return scanApplications(rows)
}

func (r *applicationRepository) CountByEmployeeID(ctx context.Context, employeeID string, status *entity.ApplicationStatus) (int, error) {
	where, args := buildWhereClause([]string{"employee_id = $1"}, []any{employeeID}, status)
	query := "SELECT COUNT(*) FROM applications" + where
	var count int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count applications: %w", err)
	}
	return count, nil
}

func (r *applicationRepository) List(ctx context.Context, status *entity.ApplicationStatus, limit, offset int) ([]*entity.Application, error) {
	where, args := buildWhereClause(nil, nil, status)
	query := applicationSelectQuery + where + fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	defer rows.Close()
	return scanApplications(rows)
}

func (r *applicationRepository) Count(ctx context.Context, status *entity.ApplicationStatus) (int, error) {
	where, args := buildWhereClause(nil, nil, status)
	query := "SELECT COUNT(*) FROM applications" + where
	var count int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count applications: %w", err)
	}
	return count, nil
}

func (r *applicationRepository) UpdateStatus(ctx context.Context, id string, status entity.ApplicationStatus, approvedBy string, approvedAt time.Time, adminComment *string) error {
	query := `
		UPDATE applications
		SET status = $1, approved_by = $2, approved_at = $3, admin_comment = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.ExecContext(ctx, query, status, approvedBy, approvedAt, adminComment, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("application not found")
	}
	return nil
}

func (r *applicationRepository) Delete(ctx context.Context, id string) (*entity.Application, error) {
	app, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM applications WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete application: %w", err)
	}
	return app, nil
}

// buildWhereClause はWHERE句とバインド変数を構築するヘルパー
func buildWhereClause(conditions []string, args []any, status *entity.ApplicationStatus) (string, []any) {
	if args == nil {
		args = []any{}
	}
	if status != nil {
		args = append(args, *status)
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanApplication(row rowScanner) (*entity.Application, error) {
	app := &entity.Application{
		Employee:   &entity.EmployeeSummary{},
		Correction: &entity.AttendanceCorrectionRequest{},
	}
	var (
		approvedBy   sql.NullString
		approvedAt   sql.NullTime
		adminComment sql.NullString
		approverID   sql.NullString
		approverName sql.NullString
		correctionID sql.NullString
		corrDate     sql.NullTime
		corrClockIn  sql.NullTime
		corrClockOut sql.NullTime
	)

	err := row.Scan(
		&app.ID, &app.EmployeeID, &app.Type, &app.Status, &app.Reason,
		&approvedBy, &approvedAt, &adminComment,
		&app.CreatedAt, &app.UpdatedAt,
		&app.Employee.Name,
		&approverID, &approverName,
		&correctionID, &corrDate, &corrClockIn, &corrClockOut,
	)
	if err != nil {
		return nil, err
	}

	app.Employee.ID = app.EmployeeID

	if approvedBy.Valid {
		app.ApprovedBy = &approvedBy.String
	}
	if approvedAt.Valid {
		app.ApprovedAt = &approvedAt.Time
	}
	if adminComment.Valid {
		app.AdminComment = &adminComment.String
	}

	if approverID.Valid && approverName.Valid {
		app.Approver = &entity.EmployeeSummary{
			ID:   approverID.String,
			Name: approverName.String,
		}
	}

	if correctionID.Valid {
		app.Correction.ID = correctionID.String
		app.Correction.ApplicationID = app.ID
		app.Correction.Date = corrDate.Time
		if corrClockIn.Valid {
			app.Correction.RequestedClockIn = &corrClockIn.Time
		}
		if corrClockOut.Valid {
			app.Correction.RequestedClockOut = &corrClockOut.Time
		}
	}

	return app, nil
}

func scanApplications(rows *sql.Rows) ([]*entity.Application, error) {
	var apps []*entity.Application
	for rows.Next() {
		app, err := scanApplication(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		apps = append(apps, app)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return apps, nil
}
