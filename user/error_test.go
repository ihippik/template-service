package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newBadRequest(t *testing.T) {
	type args struct {
		code string
		msg  string
	}
	tests := []struct {
		name string
		args args
		want *ServiceError
	}{
		{
			name: "invalid data",
			args: args{
				code: InvalidUserID,
				msg:  "invalid user id",
			},
			want: &ServiceError{
				HTTPCode: 400,
				Code:     "INVALID_USER_ID",
				Message:  "invalid user id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newBadRequest(tt.args.code, tt.args.msg), "newBadRequest(%v, %v)", tt.args.code, tt.args.msg)
		})
	}
}
func Test_newInternalServer(t *testing.T) {
	type args struct {
		code string
		msg  string
	}
	tests := []struct {
		name string
		args args
		want *ServiceError
	}{
		{
			name: "invalid data",
			args: args{
				code: InternalServerError,
				msg:  "connection refused",
			},
			want: &ServiceError{
				HTTPCode: 500,
				Code:     "INTERNAL_SERVER_ERROR",
				Message:  "connection refused",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newInternalServer(tt.args.code, tt.args.msg), "newBadRequest(%v, %v)", tt.args.code, tt.args.msg)
		})
	}
}
