package usecase

import (
	"context"
	"fmt"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
)

type AdminUseCase interface {
	ListEmployees(ctx context.Context, page, perPage int, role *entity.UserRole) ([]*entity.User, int, error)
}

type adminUseCase struct {
	userRepo repository.UserRepository
}

func NewAdminUseCase(userRepo repository.UserRepository) AdminUseCase {
	return &adminUseCase{userRepo: userRepo}
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
