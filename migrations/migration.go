package migrations

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Up(conn string) error {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return fmt.Errorf("db connection: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("embeded: %w", err)
	}

	if err := goose.Up(db.DB, "."); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func Down(conn string) error {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return fmt.Errorf("db connection: %w", err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("embeded: %w", err)
	}

	if err := goose.Down(db.DB, "."); err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}

	return nil
}
