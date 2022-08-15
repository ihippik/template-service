package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type service interface {
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	ListUser(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, dto DTO) (*User, error)
	CreateUser(ctx context.Context, dto DTO) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type Endpoint struct {
	logger *zap.Logger
	svc    service
}

func NewEndpoint(logger *zap.Logger, svc service) *Endpoint {
	return &Endpoint{logger: logger, svc: svc}
}

type response struct {
	Data []*User `json:"data,omitempty"`
}

// ListUsers http list users handler.
// @Title List
// @Tags User
// @Accept json
// @Produce json
// @Description list user
// @Summary fetch user
// @Success 200 {object} response
// @Failure 500 {object} ServiceError
// @Router /v1/users [GET]
func (e *Endpoint) ListUsers(w http.ResponseWriter, r *http.Request) {
	var resp response

	w.Header().Set("Content-Type", "application/json")

	models, err := e.svc.ListUser(r.Context())
	if err != nil {
		e.writeErr(w, err)
		return
	}

	resp.Data = models
	e.writeResp(w, resp)
}

// CreateUser http create user handler.
// @Title Create
// @Tags User
// @Accept json
// @Produce json
// @Description create user by id
// @Summary create user
// @Success 200 {object} response
// @Failure 400 {object} ServiceError
// @Failure 422 {object} ServiceError
// @Failure 404 {object} ServiceError
// @Failure 500 {object} ServiceError
// @Param model body DTO true "New model"
// @Router /v1/users [POST]
func (e *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		dto  DTO
		resp response
	)

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		e.logger.Warn("decode user data", zap.Error(err))
		e.writeErr(w, newBadRequest(InvalidUserData, err.Error()))

		return
	}

	model, err := e.svc.CreateUser(r.Context(), dto)
	if err != nil {
		e.writeErr(w, err)
		return
	}

	resp.Data = append(resp.Data, model)

	w.WriteHeader(http.StatusCreated)

	e.writeResp(w, resp)
}

// UpdateUser http update user handler.
// @Title Update
// @Tags User
// @Accept json
// @Produce json
// @Description update user by id
// @Summary update user
// @Success 200 {object} response
// @Failure 400 {object} ServiceError
// @Failure 422 {object} ServiceError
// @Failure 404 {object} ServiceError
// @Failure 500 {object} ServiceError
// @Param id path string true "User ID"
// @Param model body DTO true "New model"
// @Router /v1/users/{id} [PUT]
func (e *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		dto  DTO
		resp response
	)

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		e.logger.Warn("could not parse user id", zap.Error(err))
		e.writeErr(w, newBadRequest(InvalidUserID, err.Error()))

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		e.logger.Warn("decode user data", zap.Error(err))
		e.writeErr(w, newBadRequest(InvalidUserData, err.Error()))

		return
	}

	model, err := e.svc.UpdateUser(r.Context(), id, dto)
	if err != nil {
		e.writeErr(w, err)
		return
	}

	resp.Data = append(resp.Data, model)

	e.writeResp(w, resp)
}

// GetUser http get user handler.
// @Title Get
// @Tags User
// @Accept json
// @Produce json
// @Description get user by id
// @Summary get user
// @Success 200 {object} response
// @Failure 400 {object} ServiceError
// @Failure 404 {object} ServiceError
// @Failure 500 {object} ServiceError
// @Param id path string true "User ID"
// @Router /v1/users/{id} [GET]
func (e *Endpoint) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		e.logger.Warn("could not parse user id", zap.Error(err))
		e.writeErr(w, newBadRequest(InvalidUserID, err.Error()))

		return
	}

	model, err := e.svc.GetUser(r.Context(), id)
	if err != nil {
		e.writeErr(w, err)
		return
	}

	var resp response

	resp.Data = append(resp.Data, model)

	e.writeResp(w, resp)
}

// DeleteUser http delete user handler.
// @Title Delete
// @Tags User
// @Accept json
// @Produce json
// @Description delete user by id
// @Summary delete user
// @Success 200 {object} response
// @Failure 500 {object} ServiceError
// @Param id path string true "User ID"
// @Router /v1/users/{id} [DELETE]
func (e *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		e.logger.Warn("could not parse user id", zap.Error(err))
		e.writeErr(w, newBadRequest(InvalidUserID, err.Error()))

		return
	}

	if err := e.svc.DeleteUser(r.Context(), id); err != nil {
		e.writeErr(w, err)
		return
	}

	var resp response

	resp.Data = []*User{}

	e.writeResp(w, resp)
}

func (e *Endpoint) writeResp(w http.ResponseWriter, uData any) {
	data, err := json.Marshal(uData)
	if err != nil {
		e.writeErr(w, newInternalServer(InternalServerError, err.Error()))
		e.logger.Warn("marshal data", zap.Error(err))
		return
	}

	if _, err = w.Write(data); err != nil {
		e.logger.Error("write error", zap.Error(err))
	}
}

func (e *Endpoint) writeErr(w http.ResponseWriter, err error) {
	var svcErr *ServiceError

	if errors.As(err, &svcErr) {
		data, err := json.Marshal(svcErr)
		if err != nil {
			e.logger.Error("marshal server err", zap.Error(err))
			return
		}

		w.WriteHeader(svcErr.HTTPCode)

		if _, err = w.Write(data); err != nil {
			e.logger.Error("write error", zap.Error(err))
		}

		return
	}

	data, err := json.Marshal(errInternalServer)
	if err != nil {
		e.logger.Error("marshal internal err", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(data)
}
