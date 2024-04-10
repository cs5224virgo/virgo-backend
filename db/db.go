package db

import (
	"fmt"

	"github.com/cs5224virgo/virgo/logger"
	perrors "github.com/pkg/errors"
	"github.com/spf13/viper"

	"database/sql"

	sqlc "github.com/cs5224virgo/virgo/db/generated"

	_ "github.com/lib/pq"
)

type DB struct {
	DB      *sql.DB
	Queries *sqlc.Queries
}

// InitDB : connect to the db, and return the controller variable to be passed around
func InitDB() (*DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", viper.GetString("db.username"), viper.GetString("db.password"), viper.GetString("db.hostname"), viper.GetString("db.port"), viper.GetString("db.name"))
	logger.Info("dburl is:", dbURL)
	logger.Info("Opening a connection to the db...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, perrors.Wrap(err, "cannot connect to db")
	}
	sqlcdb := sqlc.New(db)

	return &DB{
		DB:      db,
		Queries: sqlcdb,
	}, nil
}
