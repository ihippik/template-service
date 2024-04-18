package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository is a database PostgreSQL repository.
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates new Repository instance.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// List receive all user from the database.
func (r *Repository) List(ctx context.Context) ([]*User, error) {
	var models []*User

	rows, err := r.db.QueryxContext(
		ctx,
		"SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users",
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	for rows.Next() {
		var model User

		if err := rows.StructScan(&model); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		models = append(models, &model)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close: %w", err)
	}

	return models, nil
}

// Get receive user form the database by her id.
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	var model User

	err := r.db.QueryRowxContext(ctx,
		"SELECT id, first_name, last_name, birthday, created_at, updated_at FROM users WHERE id=$1",
		id,
	).StructScan(&model)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, errNotExists
	case err != nil:
		return nil, fmt.Errorf("exec: %w", err)
	}

	return &model, nil
}

// Update user form the database by her id.
func (r *Repository) Update(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE users SET first_name=$1, last_name=$2, birthday=$3, updated_at=$4 WHERE id=$5",
		user.FirstName,
		user.LastName,
		user.Birthday,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

// Create new user in the database.
func (r *Repository) Create(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (id, first_name, last_name, birthday, created_at) VALUES($1, $2, $3, $4, $5)",
		user.ID,
		user.FirstName,
		user.LastName,
		user.Birthday,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

// Delete user from the database by her id.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
