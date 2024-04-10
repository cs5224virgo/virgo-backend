package db

import (
	"errors"
	"fmt"
	"os"

	"database/sql"

	"github.com/cs5224virgo/virgo/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	perrors "github.com/pkg/errors"
	"github.com/spf13/viper"
)

func rollbackMigrateVersion(db *sql.DB, m *migrate.Migrate) error {
	// if currently ver 1 and dirty, rollback by doing drop table schema_migrations (ver 0); else force previous version
	version, dirty, err := m.Version()
	if err != nil {
		return perrors.Wrap(err, "unable to run migration")
	}

	if version == 1 && dirty {
		_, err := db.Exec("DROP TABLE schema_migrations")
		if err != nil {
			return perrors.Wrap(err, "unable to drop schema_migrations")
		}
	} else {
		err := m.Force(int(version) - 1)
		if err != nil {
			return perrors.Wrap(err, "cannot force to previous version")
		}
	}

	return nil
}

func Migrate() error {
	cwd, err := os.Getwd()
	if err != nil {
		return perrors.Wrap(err, "cant get current directory")
	}
	migrateDirectory := fmt.Sprintf("file://%s/db/migrations", cwd)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", viper.GetString("db.username"), viper.GetString("db.password"), viper.GetString("db.hostname"), viper.GetString("db.port"), viper.GetString("db.name"))

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return perrors.Wrap(err, "cannot connect to db")
	}

	tx, err := db.Begin()
	if err != nil {
		return perrors.Wrap(err, "cannot start db transaction")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		tx.Rollback()
		return perrors.Wrap(err, "cannot init postgres driver")
	}
	m, err := migrate.NewWithDatabaseInstance(migrateDirectory, "postgres", driver)
	if err != nil {
		tx.Rollback()
		return perrors.Wrap(err, "cannot open migrations folder")
	}

	logger.Info("connected to db")

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			tx.Rollback()
			logger.Info("already up to date.")
			return nil
		}
		err2 := tx.Rollback()
		if err2 != nil {
			logger.Warn("unable to rollback, need manual fixing", err2)
		}
		err2 = rollbackMigrateVersion(db, m)
		if err2 != nil {
			logger.Warn("unable to rollback #2, need manual fixing", err2)
		}
		return perrors.Wrap(err, "unable to run migration")
	}

	err = tx.Commit()
	if err != nil {
		return perrors.Wrap(err, "cannot commit transaction")
	}

	logger.Info("migrate up successful")
	return nil
}

func PrintVersion() error {
	cwd, err := os.Getwd()
	if err != nil {
		return perrors.Wrap(err, "cant get current directory")
	}
	migrateDirectory := fmt.Sprintf("file://%s/db/migrations", cwd)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", viper.GetString("db.username"), viper.GetString("db.password"), viper.GetString("db.hostname"), viper.GetString("db.port"), viper.GetString("db.name"))
	m, err := migrate.New(migrateDirectory, dbURL)
	if err != nil {
		return perrors.Wrap(err, "error while initiating db migration")
	}
	version, dirty, err := m.Version()
	if err != nil {
		return perrors.Wrap(err, "unable to query for version")
	}
	logger.Infof("Current migration version=%d dirty=%t\n", version, dirty)
	return nil
}
