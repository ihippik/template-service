package user

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestEndpoint_ListUsers(t *testing.T) {
	svc := new(MockServer)

	setList := func(users []*User, err error) {
		svc.On("ListUser", mock.Anything).Return(users, err).Once()
	}

	tests := []struct {
		name         string
		setup        func()
		wantHTTPCode int
		want         []byte
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
							CreatedAt: time.Date(2022, 11, 17, 20, 0, 0, 0, time.UTC),
							UpdatedAt: nil,
						},
					},
					nil,
				)
			},
			wantHTTPCode: http.StatusOK,
			want:         []byte(`{"data":[{"id":"ccae37ea-d41e-4371-a3a3-89203b9e2608","firstName":"Elon","lastName":"Musk","birthday":"1971-06-28","createdAt":"2022-11-17T20:00:00Z","updatedAt":null}]}`),
		},
		{
			name: "svc error",
			setup: func() {
				setList(
					[]*User{},
					newInternalServer(InternalServerError, "internal server error"),
				)
			},
			wantHTTPCode: http.StatusInternalServerError,
			want:         []byte(`{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`),
		},
	}

	e := &Endpoint{
		logger: zap.NewNop(),
		svc:    svc,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer svc.AssertExpectations(t)

			tt.setup()

			req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
			w := httptest.NewRecorder()

			e.ListUsers(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantHTTPCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, data)
		})
	}
}

func TestEndpoint_GetUser(t *testing.T) {
	type args struct {
		id string
	}

	svc := new(MockServer)

	setGet := func(id uuid.UUID, user *User, err error) {
		svc.On("GetUser", mock.Anything, id).Return(user, err).Once()
	}

	tests := []struct {
		name         string
		args         args
		setup        func()
		wantHTTPCode int
		want         []byte
	}{
		{
			name: "success",
			args: args{
				id: "ccae37ea-d41e-4371-a3a3-89203b9e2608",
			},
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Musk",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 11, 17, 20, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)
			},
			wantHTTPCode: http.StatusOK,
			want:         []byte(`{"data":[{"id":"ccae37ea-d41e-4371-a3a3-89203b9e2608","firstName":"Elon","lastName":"Musk","birthday":"1971-06-28","createdAt":"2022-11-17T20:00:00Z","updatedAt":null}]}`),
		},
		{
			name: "svc error",
			args: args{
				id: "ccae37ea-d41e-4371-a3a3-89203b9e2608",
			},
			setup: func() {
				setGet(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					nil,
					newInternalServer(InternalServerError, "internal server error"),
				)
			},
			wantHTTPCode: http.StatusInternalServerError,
			want:         []byte(`{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`),
		},
		{
			name: "invalid id",
			args: args{
				id: "invalid",
			},
			setup: func() {
			},
			wantHTTPCode: http.StatusBadRequest,
			want:         []byte(`{"code":"INVALID_USER_ID","message":"invalid UUID length: 7"}`),
		},
	}

	e := &Endpoint{
		logger: zap.NewNop(),
		svc:    svc,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer svc.AssertExpectations(t)

			tt.setup()

			req := httptest.NewRequest(http.MethodGet, "/v1/users/ccae37ea-d41e-4371-a3a3-89203b9e2608", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.args.id})
			w := httptest.NewRecorder()

			e.GetUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantHTTPCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, data)
		})
	}
}

func TestEndpoint_UpdateUser(t *testing.T) {
	type args struct {
		id  string
		dto []byte
	}

	svc := new(MockServer)

	setUpdate := func(id uuid.UUID, dto DTO, user *User, err error) {
		svc.On("UpdateUser", mock.Anything, id, dto).Return(user, err).Once()
	}

	tests := []struct {
		name         string
		args         args
		setup        func()
		wantHTTPCode int
		want         []byte
	}{
		{
			name: "success",
			args: args{
				id:  "ccae37ea-d41e-4371-a3a3-89203b9e2608",
				dto: []byte(`{"lastName":"Rogozin","firstName":"Elon","birthday":"1971-06-28"}`),
			},
			setup: func() {
				setUpdate(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					DTO{
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
					},
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 11, 17, 20, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)
			},
			wantHTTPCode: http.StatusOK,
			want:         []byte(`{"data":[{"id":"ccae37ea-d41e-4371-a3a3-89203b9e2608","firstName":"Elon","lastName":"Rogozin","birthday":"1971-06-28","createdAt":"2022-11-17T20:00:00Z","updatedAt":null}]}`),
		},
		{
			name: "svc error",
			args: args{
				id:  "ccae37ea-d41e-4371-a3a3-89203b9e2608",
				dto: []byte(`{"lastName":"Rogozin","firstName":"Elon","birthday":"1971-06-28"}`),
			},
			setup: func() {
				setUpdate(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					DTO{
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
					},
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 11, 17, 20, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					newInternalServer(InternalServerError, "internal server error"),
				)
			},
			wantHTTPCode: http.StatusInternalServerError,
			want:         []byte(`{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`),
		},
		{
			name: "invalid id",
			args: args{
				id:  "invalid",
				dto: []byte(`{"lastName":"Rogozin","firstName":"Elon","birthday":"1971-06-28"}`),
			},
			setup: func() {
			},
			wantHTTPCode: http.StatusBadRequest,
			want:         []byte(`{"code":"INVALID_USER_ID","message":"invalid UUID length: 7"}`),
		},
		{
			name: "invalid body",
			args: args{
				id:  "ccae37ea-d41e-4371-a3a3-89203b9e2608",
				dto: []byte(`invalid`),
			},
			setup: func() {
			},
			wantHTTPCode: http.StatusBadRequest,
			want:         []byte(`{"code":"INVALID_USER_DATA","message":"invalid character 'i' looking for beginning of value"}`),
		},
	}

	e := &Endpoint{
		logger: zap.NewNop(),
		svc:    svc,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer svc.AssertExpectations(t)

			tt.setup()

			req := httptest.NewRequest(
				http.MethodPut,
				"/v1/users/ccae37ea-d41e-4371-a3a3-89203b9e2608",
				bytes.NewReader(tt.args.dto),
			)
			req = mux.SetURLVars(req, map[string]string{"id": tt.args.id})
			w := httptest.NewRecorder()

			e.UpdateUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantHTTPCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, data)
		})
	}
}

func TestEndpoint_CreateUser(t *testing.T) {
	type args struct {
		dto []byte
	}

	svc := new(MockServer)

	setCreate := func(dto DTO, user *User, err error) {
		svc.On("CreateUser", mock.Anything, dto).Return(user, err).Once()
	}

	tests := []struct {
		name         string
		args         args
		setup        func()
		wantHTTPCode int
		want         []byte
	}{
		{
			name: "success",
			args: args{
				dto: []byte(`{"lastName":"Rogozin","firstName":"Elon","birthday":"1971-06-28"}`),
			},
			setup: func() {
				setCreate(
					DTO{
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
					},
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 11, 17, 20, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					nil,
				)
			},
			wantHTTPCode: http.StatusCreated,
			want:         []byte(`{"data":[{"id":"ccae37ea-d41e-4371-a3a3-89203b9e2608","firstName":"Elon","lastName":"Rogozin","birthday":"1971-06-28","createdAt":"2022-11-17T20:00:00Z","updatedAt":null}]}`),
		},
		{
			name: "svc error",
			args: args{
				dto: []byte(`{"lastName":"Rogozin","firstName":"Elon","birthday":"1971-06-28"}`),
			},
			setup: func() {
				setCreate(
					DTO{
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
					},
					&User{
						ID:        uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
						FirstName: "Elon",
						LastName:  "Rogozin",
						Birthday:  "1971-06-28",
						CreatedAt: time.Date(2022, 11, 17, 20, 0, 0, 0, time.UTC),
						UpdatedAt: nil,
					},
					newInternalServer(InternalServerError, "internal server error"),
				)
			},
			wantHTTPCode: http.StatusInternalServerError,
			want:         []byte(`{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`),
		},
		{
			name: "invalid body",
			args: args{
				dto: []byte(`invalid`),
			},
			setup: func() {
			},
			wantHTTPCode: http.StatusBadRequest,
			want:         []byte(`{"code":"INVALID_USER_DATA","message":"invalid character 'i' looking for beginning of value"}`),
		},
	}

	e := &Endpoint{
		logger: zap.NewNop(),
		svc:    svc,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer svc.AssertExpectations(t)

			tt.setup()

			req := httptest.NewRequest(
				http.MethodPost,
				"/v1/users",
				bytes.NewReader(tt.args.dto),
			)

			w := httptest.NewRecorder()

			e.CreateUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantHTTPCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, data)
		})
	}
}

func TestEndpoint_DeleteUser(t *testing.T) {
	type args struct {
		id string
	}

	svc := new(MockServer)

	setDelete := func(id uuid.UUID, err error) {
		svc.On("DeleteUser", mock.Anything, id).Return(err).Once()
	}

	tests := []struct {
		name         string
		args         args
		setup        func()
		wantHTTPCode int
		want         []byte
	}{
		{
			name: "success",
			args: args{
				id: "ccae37ea-d41e-4371-a3a3-89203b9e2608",
			},
			setup: func() {
				setDelete(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					nil,
				)
			},
			wantHTTPCode: http.StatusOK,
			want:         []byte(`{}`),
		},
		{
			name: "svc error",
			args: args{
				id: "ccae37ea-d41e-4371-a3a3-89203b9e2608",
			},
			setup: func() {
				setDelete(
					uuid.MustParse("ccae37ea-d41e-4371-a3a3-89203b9e2608"),
					newInternalServer(InternalServerError, "internal server error"),
				)
			},
			wantHTTPCode: http.StatusInternalServerError,
			want:         []byte(`{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`),
		},
		{
			name: "invalid id",
			args: args{
				id: "invalid",
			},
			setup: func() {
			},
			wantHTTPCode: http.StatusBadRequest,
			want:         []byte(`{"code":"INVALID_USER_ID","message":"invalid UUID length: 7"}`),
		},
	}

	e := &Endpoint{
		logger: zap.NewNop(),
		svc:    svc,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer svc.AssertExpectations(t)

			tt.setup()

			req := httptest.NewRequest(http.MethodDelete, "/v1/users/ccae37ea-d41e-4371-a3a3-89203b9e2608", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.args.id})
			w := httptest.NewRecorder()

			e.DeleteUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantHTTPCode, res.StatusCode)

			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, data)
		})
	}
}

func TestNewEndpoint(t *testing.T) {
	type args struct {
		logger *zap.Logger
		svc    service
	}

	svc := new(MockServer)

	tests := []struct {
		name string
		args args
		want *Endpoint
	}{
		{
			name: "success",
			args: args{
				logger: zap.NewNop(),
				svc:    svc,
			},
			want: &Endpoint{
				logger: zap.NewNop(),
				svc:    svc,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewEndpoint(tt.args.logger, tt.args.svc), "NewEndpoint(%v, %v)", tt.args.logger, tt.args.svc)
		})
	}
}
