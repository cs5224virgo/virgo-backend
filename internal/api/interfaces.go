package api

import (
	sqlc "github.com/cs5224virgo/virgo/db/generated"
	"github.com/cs5224virgo/virgo/internal/datalayer"
)

type APIDataLayer interface {
	IsUsernameAvailable(username string) (bool, error)
	CreateUser(params sqlc.CreateUserParams) error
	AuthenticateUser(username string, hashedPassword string) (user sqlc.User, token string, err error)
	GetUserByID(id uint) (*sqlc.User, error)
	GetRoomsByUserID(userID int32) ([]datalayer.DetailedRoom, error)
}
