package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
	"github.com/myuto/attendance-backend/pkg/utils"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password, name string, role entity.EmployeeRole) (*entity.Employee, error)
	Login(ctx context.Context, email, password string) (string, *entity.Employee, error)
	ValidateToken(ctx context.Context, token string) (*entity.Employee, error)
}

type authUseCase struct {
	employeeRepo repository.EmployeeRepository
	jwtSecret    string
	jwtExp       time.Duration
}

func NewAuthUseCase(employeeRepo repository.EmployeeRepository, jwtSecret string, jwtExp time.Duration) AuthUseCase {
	return &authUseCase{
		employeeRepo: employeeRepo,
		jwtSecret:    jwtSecret,
		jwtExp:       jwtExp,
	}
}

func (uc *authUseCase) Register(ctx context.Context, email, password, name string, role entity.EmployeeRole) (*entity.Employee, error) {
	existingEmployee, _ := uc.employeeRepo.FindByEmail(ctx, email)
	if existingEmployee != nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if !role.IsValid() {
		role = entity.RoleEmployee
	}

	employeeID := uuid.New().String()

	employee := &entity.Employee{
		ID:           employeeID,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.employeeRepo.Create(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return employee, nil
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (string, *entity.Employee, error) {
	employee, err := uc.employeeRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(employee.ID, string(employee.Role), uc.jwtSecret, uc.jwtExp)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, employee, nil
}

func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*entity.Employee, error) {
	claims, err := utils.ValidateJWT(token, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	employee, err := uc.employeeRepo.FindByID(ctx, claims.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return employee, nil
}
