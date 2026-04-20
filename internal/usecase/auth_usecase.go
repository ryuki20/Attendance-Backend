package usecase

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/internal/domain/repository"
	"github.com/myuto/attendance-backend/pkg/utils"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password, name string, role entity.UserRole) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, *entity.User, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}

type authUseCase struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExp    time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExp time.Duration) AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExp:    jwtExp,
	}
}

func (uc *authUseCase) Register(ctx context.Context, email, password, name string, role entity.UserRole) (*entity.User, error) {
	// メールアドレスの重複チェック
	existingUser, _ := uc.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ロールのバリデーション
	if !role.IsValid() {
		role = entity.RoleEmployee
	}

	// ユーザーの作成
	user := &entity.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	// ユーザーの取得
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// JWTトークンの生成
	token, err := utils.GenerateJWT(user.ID, string(user.Role), uc.jwtSecret, uc.jwtExp)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	// トークンの検証
	claims, err := utils.ValidateJWT(token, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// ユーザーの取得
	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
