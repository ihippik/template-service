package user

import "net/http"

// service error codes.
const (
	InvalidUserID       = "INVALID_USER_ID"
	InvalidUserData     = "INVALID_USER_DATA"
	InternalServerError = "INTERNAL_SERVER_ERROR"
	NotFound            = "NOT_FOUND"
	ValidationError     = "VALIDATION_ERROR"
)

// ServiceError represent service custom error.
type ServiceError struct {
	HTTPCode int    `json:"-"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
}

// Error implement Error interface.
func (e ServiceError) Error() string {
	return e.Message
}

var errInternalServer = &ServiceError{
	HTTPCode: http.StatusInternalServerError,
	Code:     "INTERNAL_SERVER",
	Message:  "internal server error",
}

func newBadRequest(code, msg string) *ServiceError {
	return &ServiceError{HTTPCode: http.StatusBadRequest, Code: code, Message: msg}
}

func newInternalServer(code, msg string) *ServiceError {
	return &ServiceError{HTTPCode: http.StatusInternalServerError, Code: code, Message: msg}
}

func newNotFoundErr(code, msg string) *ServiceError {
	return &ServiceError{HTTPCode: http.StatusNotFound, Code: code, Message: msg}
}

func newValidationErr(code, msg string) *ServiceError {
	return &ServiceError{HTTPCode: http.StatusUnprocessableEntity, Code: code, Message: msg}
}
