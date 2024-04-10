package api

import (
	sqlc "github.com/cs5224virgo/virgo/db/generated"
)

type APIDataLayer interface {
	IsUsernameAvailable(username string) (bool, error)
	CreateUser(params sqlc.CreateUserParams) error
}
