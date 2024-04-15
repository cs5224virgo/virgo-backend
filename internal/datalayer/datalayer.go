package datalayer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cs5224virgo/virgo/db"
	sqlc "github.com/cs5224virgo/virgo/db/generated"
	"github.com/cs5224virgo/virgo/internal/jwt"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/sqids/sqids-go"
	"golang.org/x/crypto/bcrypt"
)

const sqidconst = 420

var ErrLoginFailed = errors.New("authentication failed")
var ErrIDZero = errors.New("id is zero")

type DataLayer struct {
	DB    *db.DB
	sqids *sqids.Sqids
}

func NewDataLayer(db *db.DB) *DataLayer {
	sqids, err := sqids.New(sqids.Options{
		Alphabet:  "tGWpFSsDab063e1o8wTkx5i9MuOYQBJXcjynK2HI7m4UPZfEdrLlRqAhVvNzCg",
		MinLength: 7,
	})
	if err != nil {
		logger.Fatal("cannot init sqids", err)
	}
	return &DataLayer{
		DB:    db,
		sqids: sqids,
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

func (s *DataLayer) GetUserByID(id uint) (*sqlc.User, error) {
	if id == 0 {
		return nil, ErrIDZero
	}
	user, err := s.DB.Queries.GetUserByID(context.Background(), int32(id))
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	return &user, nil
}

func (s *DataLayer) AuthenticateUser(username string, pepperedPassword string) (user sqlc.User, token string, err error) {
	if username == "" {
		err = fmt.Errorf("username is blank")
		return
	}
	if pepperedPassword == "" {
		err = fmt.Errorf("password is blank")
		return
	}
	user, err = s.DB.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("login failed: user doesn't exist")
			err = ErrLoginFailed
			return
		}
		err = fmt.Errorf("database failure: %w", err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pepperedPassword))
	if err != nil {
		logger.Info("compare hash and password failed: ", err)
		err = ErrLoginFailed
		return
	}

	// expire the token in a day
	expire := time.Now().Add(time.Hour * 24)
	token, err = jwt.NewToken(uint(user.ID), user.Username, expire)
	if err != nil {
		err = fmt.Errorf("cannot generate a new jwt token: %w", err)
		return
	}
	return
}
