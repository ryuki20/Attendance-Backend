package repository

import (
	"context"

	"github.com/myuto/attendance-backend/internal/domain/entity"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *entity.Employee) error
	FindByID(ctx context.Context, id string) (*entity.Employee, error)
	FindByEmail(ctx context.Context, email string) (*entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	Delete(ctx context.Context, id string) (*entity.Employee, error)
	List(ctx context.Context, limit, offset int, role *entity.EmployeeRole) ([]*entity.Employee, error)
	Count(ctx context.Context, role *entity.EmployeeRole) (int, error)
}
