package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockServer struct {
	mock.Mock
}

func (m *MockServer) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockServer) ListUser(ctx context.Context) ([]*User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*User), args.Error(1)
}

func (m *MockServer) UpdateUser(ctx context.Context, id uuid.UUID, dto DTO) (*User, error) {
	args := m.Called(ctx, id, dto)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockServer) CreateUser(ctx context.Context, dto DTO) (*User, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockServer) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
