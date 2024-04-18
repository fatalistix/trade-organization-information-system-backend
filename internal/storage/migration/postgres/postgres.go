package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"io/fs"
)

type Migrator struct {
	srcDriver source.Driver
}

func NewMigrator(sqlFiles fs.FS, dirName string) (*Migrator, error) {
	const op = "storage.migration.postgres.NewMigrator"

	driver, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Migrator{srcDriver: driver}, nil
}

func (m *Migrator) ApplyMigrations(db *sql.DB, dbName string) error {
	const op = "storage.migration.postgres.ApplyMigrations"

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	migrator, err := migrate.NewWithInstance("migration_embedded_sql_files", m.srcDriver, dbName, driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_, _ = migrator.Close()
	}()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: unable to apply migration: %w", op, err)
	}

	return nil
}

//func Migrate(db *sql.DB) error {
//	const op = "storage.migration.postgres.Migrate"
//
//	driver, err := postgres.WithInstance(db, &postgres.Config{})
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	m, err := migrate.NewWithDatabaseInstance(
//		"file:///db/migrations",
//		"postgres", driver,
//	)
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	err = m.Up()
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//
//	return nil
//}