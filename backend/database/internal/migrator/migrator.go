package migrator

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PrototypeSirius/ruglogger/apperror"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Run(log *logrus.Logger, db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return apperror.SystemError(err, 1501, "could not create database driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return apperror.SystemError(err, 1502, "could not create migrate instance")
	}
	err = m.Up()
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	version, dirty, verErr := m.Version()
	if verErr != nil {
		return apperror.SystemError(verErr, 1503, "could not get migrate version")
	}
	if dirty {
		log.Info(fmt.Sprintf("Dirty migration detected at version %d", version))
		prevVersion := int(version) - 1
		if prevVersion < 0 {
			prevVersion = 0
		}
		if forceErr := m.Force(prevVersion); forceErr != nil {
			return apperror.SystemError(forceErr, 1504, fmt.Sprintf("failed to force rollback to %d", prevVersion))
		}
		log.Info(fmt.Sprintf("Successfully forced version to %d", prevVersion))
		return apperror.SystemError(err, 1505, fmt.Sprintf("dirty migration at version %d rolled back", version))
	}
	return apperror.SystemError(err, 1506, "migration failed")
}
