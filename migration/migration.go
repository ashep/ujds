package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var fs embed.FS

func setup(db *sql.DB) (*migrate.Migrate, error) {
	srcDrv, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to load migration scripts: %w", err)
	}

	dbDrv, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", srcDrv, "ujds", dbDrv)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func Up(db *sql.DB) error {
	m, err := setup(db)
	if err != nil {
		return err
	}

	if err = m.Up(); errors.Is(err, migrate.ErrNoChange) {
	} else if err != nil {
		return err
	}

	return nil
}

func Down(db *sql.DB) error {
	m, err := setup(db)
	if err != nil {
		return err
	}

	if err = m.Down(); errors.Is(err, migrate.ErrNoChange) {
	} else if err != nil {
		return err
	}

	return nil
}
