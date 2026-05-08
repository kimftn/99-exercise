package migration

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	postgresdriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const defaultMigrationsPath = "file://db/migrations"

type Runner struct {
	migrator *migrate.Migrate
}

func NewRunner(databaseURL string) (*Runner, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open migration database: %w", err)
	}

	driver, err := postgresdriver.WithInstance(db, &postgresdriver.Config{})
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create migration driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(defaultMigrationsPath, "postgres", driver)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return &Runner{migrator: migrator}, nil
}

func (r *Runner) Up() error {
	if err := r.migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (r *Runner) Down(steps int) error {
	var err error
	if steps <= 0 {
		err = r.migrator.Down()
	} else {
		err = r.migrator.Steps(-steps)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (r *Runner) Version() (uint, bool, error) {
	version, dirty, err := r.migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}

	return version, dirty, err
}

func (r *Runner) Close() error {
	sourceErr, databaseErr := r.migrator.Close()
	if sourceErr != nil {
		return sourceErr
	}

	if databaseErr != nil {
		return databaseErr
	}

	return nil
}
