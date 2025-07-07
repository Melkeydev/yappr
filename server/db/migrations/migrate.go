package migration

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationsFS embed.FS

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func MigrateDown(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(db, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rollback completed successfully")
	return nil
}

func MigrateToVersion(db *sql.DB, version int64) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.UpTo(db, "migrations", version); err != nil {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	log.Printf("Migration to version %d completed successfully\n", version)
	return nil
}

func MigrateReset(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Reset(db, "migrations"); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	log.Println("Migration reset completed successfully")
	return nil
}

func MigrateStatus(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Status(db, "migrations"); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}
