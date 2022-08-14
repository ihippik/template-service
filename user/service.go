package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ihippik/template-service/config"
)

type repository interface {
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Service represent the main application structure.
type Service struct {
	cfg    *config.Config
	logger *zap.Logger
	repo   repository
}

// NewService creates new Service entity.
func NewService(cfg *config.Config, logger *zap.Logger, repo repository) *Service {
	return &Service{cfg: cfg, logger: logger, repo: repo}
}

// GetUser get user entity by her identification.
func (svc *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	model, err := svc.repo.Get(ctx, id)
	if err != nil {
		svc.logger.Error("could not get user", zap.Error(err))
		return nil, err
	}

	if model == nil {
		svc.logger.Warn("user not found", zap.String("id", id.String()))
		return nil, newNotFoundErr(NotFound, "user not found")
	}

	return model, nil
}

// ListUser fetch all users.
func (svc *Service) ListUser(ctx context.Context) ([]*User, error) {
	models, err := svc.repo.List(ctx)
	if err != nil {
		svc.logger.Error("could not fetch users", zap.Error(err))
		return nil, fmt.Errorf("list: %w", err)
	}

	return models, nil
}

// UpdateUser update user entity by her identification.
func (svc *Service) UpdateUser(ctx context.Context, id uuid.UUID, dto DTO) (*User, error) {
	if err := dto.Validate(); err != nil {
		svc.logger.Warn("dto validation error", zap.Error(err))
		return nil, newValidationErr(ValidationError, err.Error())
	}

	model, err := svc.repo.Get(ctx, id)
	if err != nil {
		svc.logger.Error("could not get user", zap.Error(err))
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	if model == nil {
		svc.logger.Warn("user not found", zap.String("id", id.String()))
		return nil, newNotFoundErr(NotFound, "user not found")
	}

	now := time.Now().UTC()

	model.UpdatedAt = &now
	model.FirstName = dto.FirstName
	model.LastName = dto.LastName

	if err := svc.repo.Update(ctx, model); err != nil {
		svc.logger.Error("update user error", zap.Error(err))
		return nil, err
	}

	return model, nil
}

// CreateUser create new entity user.
func (svc *Service) CreateUser(ctx context.Context, dto DTO) (*User, error) {
	if err := dto.Validate(); err != nil {
		svc.logger.Warn("dto validation error", zap.Error(err))
		return nil, newValidationErr(ValidationError, err.Error())
	}

	model := User{
		ID:        uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Birthday:  dto.Birthday,
		CreatedAt: time.Now().UTC(),
	}

	if err := svc.repo.Create(ctx, &model); err != nil {
		svc.logger.Error("could not create user", zap.Error(err))
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	return &model, nil
}

// DeleteUser delete a user by her identification.
func (svc *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		svc.logger.Error("could not delete user", zap.Error(err))
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
