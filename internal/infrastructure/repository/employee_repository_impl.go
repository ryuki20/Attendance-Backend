package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
	"github.com/myuto/attendance-backend/internal/infrastructure/database"
)

type employeeRepository struct {
	db *database.DB
}

func NewEmployeeRepository(db *database.DB) repository.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	query := `
		INSERT INTO employees (id, email, password_hash, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		employee.ID, employee.Email, employee.PasswordHash, employee.Name, employee.Role,
		employee.CreatedAt, employee.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	return nil
}

func (r *employeeRepository) FindByID(ctx context.Context, id string) (*entity.Employee, error) {
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM employees
		WHERE id = $1
	`
	employee := &entity.Employee{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&employee.ID, &employee.Email, &employee.PasswordHash, &employee.Name,
		&employee.Role, &employee.CreatedAt, &employee.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}
	return employee, nil
}

func (r *employeeRepository) FindByEmail(ctx context.Context, email string) (*entity.Employee, error) {
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM employees
		WHERE email = $1
	`
	employee := &entity.Employee{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&employee.ID, &employee.Email, &employee.PasswordHash, &employee.Name,
		&employee.Role, &employee.CreatedAt, &employee.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}
	return employee, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	query := `
		UPDATE employees
		SET email = $1, name = $2, role = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.db.ExecContext(
		ctx, query,
		employee.Email, employee.Name, employee.Role, employee.UpdatedAt, employee.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM employees WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *employeeRepository) List(ctx context.Context, limit, offset int, role *entity.EmployeeRole) ([]*entity.Employee, error) {
	var (
		query string
		args  []interface{}
	)

	if role != nil {
		query = `
			SELECT id, email, password_hash, name, role, created_at, updated_at
			FROM employees
			WHERE role = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{*role, limit, offset}
	} else {
		query = `
			SELECT id, email, password_hash, name, role, created_at, updated_at
			FROM employees
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}
	defer rows.Close()

	var employees []*entity.Employee
	for rows.Next() {
		employee := &entity.Employee{}
		err := rows.Scan(
			&employee.ID, &employee.Email, &employee.PasswordHash, &employee.Name,
			&employee.Role, &employee.CreatedAt, &employee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, employee)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return employees, nil
}

func (r *employeeRepository) Count(ctx context.Context, role *entity.EmployeeRole) (int, error) {
	var (
		query string
		args  []interface{}
	)

	if role != nil {
		query = `SELECT COUNT(*) FROM employees WHERE role = $1`
		args = []interface{}{*role}
	} else {
		query = `SELECT COUNT(*) FROM employees`
	}

	var count int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}
	return count, nil
}
