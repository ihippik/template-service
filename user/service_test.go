package user

import (
	"bytes"
	"context"
	"errors"
	"github.com/ihippik/template-service/config"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestService_GetUser(t *testing.T) {
	type args struct {
		id uuid.UUID
	}

	repo := new(MockRepo)

	setGet := func(id uuid.UUID, user *User, err error) {
		repo.On("Get", mock.Anything, id).Return(user, err).Once()
	}

	tests := []struct {
		name    string
		setup   func()
		args    args
		want    *User
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)
			},
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
			},
			want: &User{
				ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  "1971-06-28",
				CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: nil,
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					nil,
					errors.New("some error"),
				)
			},
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
			},
			want:    nil,
			wantErr: errors.New("some error"),
		},
		{
			name: "not found",
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					nil,
					nil,
				)
			},
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
			},
			want:    nil,
			wantErr: newNotFoundErr(NotFound, "user not found"),
		},
	}

	svc := &Service{logger: zap.NewNop(), repo: repo}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.GetUser(context.Background(), tt.args.id)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	type args struct {
		id uuid.UUID
	}

	repo := new(MockRepo)

	setDelete := func(id uuid.UUID, err error) {
		repo.On("Delete", mock.Anything, id).Return(err).Once()
	}

	tests := []struct {
		name    string
		setup   func()
		args    args
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setDelete(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					nil,
				)
			},
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setDelete(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					errors.New("some error"),
				)
			},
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
			},
			wantErr: errors.New("delete user: some error"),
		},
	}

	svc := &Service{logger: zap.NewNop(), repo: repo}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			err := svc.DeleteUser(context.Background(), tt.args.id)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	type args struct {
		dto DTO
	}

	wayback := time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	repo := new(MockRepo)

	setCreate := func(user *User, err error) {
		repo.On("Create", mock.Anything, user).Return(err).Once()
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *User
		wantErr error
	}{
		{
			name: "success",
			args: args{dto: DTO{
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  "1971-06-28",
			}},
			setup: func() {
				reader := bytes.NewReader([]byte("1111111111111111"))
				uuid.SetRand(reader)

				setCreate(
					&User{
						ID:        uuid.MustParse("31313131-3131-4131-b131-313131313131"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)
			},
			want: &User{
				ID:        uuid.MustParse("31313131-3131-4131-b131-313131313131"),
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  "1971-06-28",
				CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: nil,
			},
			wantErr: nil,
		},
		{
			name: "validation error",
			args: args{dto: DTO{
				FirstName: "Elon",
				LastName:  "Musk",
			}},
			setup:   func() {},
			want:    nil,
			wantErr: errors.New("Key: 'DTO.Birthday' Error:Field validation for 'Birthday' failed on the 'required' tag"),
		},
		{
			name: "some error",
			args: args{dto: DTO{
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  "1971-06-28",
			}},
			setup: func() {
				reader := bytes.NewReader([]byte("1111111111111111"))
				uuid.SetRand(reader)

				setCreate(
					&User{
						ID:        uuid.MustParse("31313131-3131-4131-b131-313131313131"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					errors.New("some error"),
				)
			},
			want:    nil,
			wantErr: errors.New("could not create user: some error"),
		},
	}

	svc := &Service{logger: zap.NewNop(), repo: repo}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.CreateUser(context.Background(), tt.args.dto)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	type args struct {
		id  uuid.UUID
		dto DTO
	}

	wayback := time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	repo := new(MockRepo)

	setGet := func(id uuid.UUID, user *User, err error) {
		repo.On("Get", mock.Anything, id).Return(user, err).Once()
	}

	setUpdate := func(user *User, err error) {
		repo.On("Update", mock.Anything, user).Return(err).Once()
	}

	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *User
		wantErr error
	}{
		{
			name: "success",
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				dto: DTO{
					FirstName: "Elon",
					LastName:  "Rogozin",
					Birthday:  "1971-06-28",
				},
			},
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)

				setUpdate(
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
					},
					nil,
				)
			},
			want: &User{
				ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				FirstName: "Elon",
				LastName:  "Rogozin",
				Birthday:  "1971-06-28",
				CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				dto: DTO{
					FirstName: "Elon",
					LastName:  "Rogozin",
					Birthday:  "1971-06-28",
				},
			},
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					nil,
					nil,
				)
			},
			want:    nil,
			wantErr: newNotFoundErr(NotFound, "user not found"),
		},
		{
			name: "validation error",
			args: args{dto: DTO{
				FirstName: "Elon",
				LastName:  "Musk",
			}},
			setup:   func() {},
			want:    nil,
			wantErr: errors.New("Key: 'DTO.Birthday' Error:Field validation for 'Birthday' failed on the 'required' tag"),
		},
		{
			name: "get: some error",
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				dto: DTO{
					FirstName: "Elon",
					LastName:  "Musk",
					Birthday:  "1971-06-28",
				},
			},
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					errors.New("some error"),
				)
			},
			want:    nil,
			wantErr: errors.New("could not get user: some error"),
		},
		{
			name: "update: some err",
			args: args{
				id: uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
				dto: DTO{
					FirstName: "Elon",
					LastName:  "Rogozin",
					Birthday:  "1971-06-28",
				},
			},
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)

				setUpdate(
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: toPointer(time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC)),
					},
					errors.New("some err"),
				)
			},
			want:    nil,
			wantErr: errors.New("some err"),
		},
	}

	svc := &Service{logger: zap.NewNop(), repo: repo}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.UpdateUser(context.Background(), tt.args.id, tt.args.dto)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_ListUser(t *testing.T) {
	repo := new(MockRepo)

	setList := func(users []*User, err error) {
		repo.On("List", mock.Anything).Return(users, err).Once()
	}

	tests := []struct {
		name    string
		setup   func()
		want    []*User
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setList(
					[]*User{
						{
							ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  "1971-06-28",
							CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt: nil,
						},
					},
					nil,
				)
			},
			want: []*User{
				{
					ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					FirstName: "Elon",
					LastName:  "Musk",
					Birthday:  "1971-06-28",
					CreatedAt: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList(nil, errors.New("some error"))
			},
			want:    nil,
			wantErr: errors.New("list: some error"),
		},
		{
			name: "empty list",
			setup: func() {
				setList(nil, nil)
			},
			want:    nil,
			wantErr: nil,
		},
	}

	svc := &Service{logger: zap.NewNop(), repo: repo}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.ListUser(context.Background())
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func toPointer[T any](d T) *T {
	return &d
}

func TestNewService(t *testing.T) {
	type args struct {
		cfg    *config.Config
		logger *zap.Logger
		repo   repository
	}

	repo := new(MockRepo)

	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "success",
			args: args{
				cfg: &config.Config{
					ServerAddr: ":80",
					DB: config.DBCfg{
						Conn:         "user=pass",
						MaxOpenConns: 10,
						MaxIdleConns: 20,
					},
					Log: config.LogCfg{
						Level:      "info",
						Caller:     false,
						StackTrace: false,
					},
				},
				logger: &zap.Logger{},
				repo:   repo,
			},
			want: &Service{
				cfg: &config.Config{
					ServerAddr: ":80",
					DB: config.DBCfg{
						Conn:         "user=pass",
						MaxOpenConns: 10,
						MaxIdleConns: 20,
					},
					Log: config.LogCfg{
						Level:      "info",
						Caller:     false,
						StackTrace: false,
					},
				},
				logger: &zap.Logger{},
				repo:   repo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewService(tt.args.cfg, tt.args.logger, tt.args.repo), "NewService(%v, %v, %v)", tt.args.cfg, tt.args.logger, tt.args.repo)
		})
	}
}
