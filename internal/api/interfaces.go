package api

import (
	sqlc "github.com/cs5224virgo/virgo/db/generated"
	"github.com/cs5224virgo/virgo/internal/datalayer"
	"github.com/gin-gonic/gin"
)

type APIDataLayer interface {
	IsUsernameAvailable(username string) (bool, error)
	CreateUser(params sqlc.CreateUserParams) error
	AuthenticateUser(username string, hashedPassword string) (user sqlc.User, token string, err error)
	GetUserByID(id uint) (*sqlc.User, error)
	GetRoomsByUserID(userID int32) ([]datalayer.DetailedRoom, error)
	GetRoomCodesByUserID(userID int32) ([]string, error)
	CreateRoom(userID int32, roomName string, roomDescription *string) (datalayer.DetailedRoom, error)
	LeaveRoom(username string, roomCode string) error
	AddUsersToRoom(usernames []string, roomCode string) (datalayer.DetailedRoom, error)
	JoinRoom(username string, roomCode string) (datalayer.DetailedRoom, error)
	ToDetailedUser(user sqlc.User) datalayer.DetailedUser
	GetAllMessagesForRoom(roomCode string) ([]datalayer.DetailedMessage, error)
}

type WebSocketHub interface {
	ServeWs(c *gin.Context, username string)
	AnnounceAddUserToRoom(user datalayer.DetailedUser, room datalayer.DetailedRoom, isNewMember bool)
	AnnounceLeaveRoom(user datalayer.DetailedUser, roomCode string)
}
