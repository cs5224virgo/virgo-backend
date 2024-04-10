package datalayer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cs5224virgo/virgo/db"
	sqlc "github.com/cs5224virgo/virgo/db/generated"
)

type DataLayer struct {
	DB *db.DB
}

func NewDataLayer(db *db.DB) *DataLayer {
	return &DataLayer{
		DB: db,
	}
}

func (s *DataLayer) IsUsernameAvailable(username string) (bool, error) {
	if username == "" {
		return false, nil
	}
	user, err := s.DB.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("query to get user by username failed: %w", err)
	}

	if user.ID == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *DataLayer) CreateUser(params sqlc.CreateUserParams) error {
	if params.Username == "" {
		return fmt.Errorf("username is blank")
	}
	if params.Password == "" {
		return fmt.Errorf("password is blank")
	}
	_, err := s.DB.Queries.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("query to create user failed: %w", err)
	}
	return nil
}
