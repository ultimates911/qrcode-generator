package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"qrcodegen/internal/dto"
	"qrcodegen/internal/repository/postgres"
	sqldb "qrcodegen/sqlc/generated"
)

var ErrUserAlreadyExists = errors.New("user with this email already exists")

// UserUseCase инкапсулирует бизнес-логику для сущности User.
type UserUseCase struct {
	repo postgres.Repository
}

func NewUserUseCase(repo postgres.Repository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	_, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	params := sqldb.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	}
	_, err = uc.repo.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.RegisterResponse{Message: "User registered successfully"}, nil
}