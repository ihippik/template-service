package user

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// User server domain struct.
type User struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `db:"first_name" json:"firstName"`
	LastName  string     `db:"last_name" json:"lastName"`
	Birthday  string     `json:"birthday"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
}

// DTO represent data transfer object for creating and updating a new entity.
type DTO struct {
	FirstName string `validate:"required" json:"firstName,omitempty"`
	LastName  string `validate:"required" json:"lastName,omitempty"`
	Birthday  string `validate:"required" json:"birthday,omitempty"`
}

// Validate check mandatory fields.
func (d DTO) Validate() error {
	validate := validator.New()

	return validate.Struct(d)
}
