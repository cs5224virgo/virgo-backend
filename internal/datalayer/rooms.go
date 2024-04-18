package datalayer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"

	sqlc "github.com/cs5224virgo/virgo/db/generated"
	"github.com/lib/pq"
)

type DetailedUser struct {
	Username    string  `json:"username"`
	DisplayName *string `json:"displayName"`
}

type DetailedUserUnread struct {
	User   DetailedUser `json:"user"`
	Unread int32        `json:"unread"`
}

type DetailedRoom struct {
	ID                 int32                `json:"id"`
	CreatedAt          time.Time            `json:"createdAt"`
	Code               string               `json:"code"`
	Name               string               `json:"name"`
	Description        *string              `json:"description"`
	LastActivity       time.Time            `json:"lastActivity"`
	LastMessagePreview string               `json:"lastMessagePreview"`
	Users              []DetailedUserUnread `json:"users"`
}

func (s *DataLayer) GetRoomsByUserID(userID int32) ([]DetailedRoom, error) {
	if userID == 0 {
		return nil, ErrIDZero
	}
	rooms, err := s.DB.Queries.GetRoomsByUser(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("database failure: %w", err)
	}

	detailedRooms := []DetailedRoom{}
	for _, room := range rooms {
		detailedRoom, err := s.toDetailedRoom(room)
		if err != nil {
			return nil, err
		}
		detailedRooms = append(detailedRooms, *detailedRoom)
	}

	slices.SortStableFunc(detailedRooms, func(a, b DetailedRoom) int {
		return b.LastActivity.Compare(a.LastActivity)
	})
	return detailedRooms, nil
}

func (s *DataLayer) toDetailedRoom(room sqlc.Room) (*DetailedRoom, error) {
	detailedRoom := DetailedRoom{
		ID:        room.ID,
		CreatedAt: room.CreatedAt,
		Code:      room.Code,
		Name:      room.Name,
	}
	if room.Description.Valid {
		detailedRoom.Description = &room.Description.String
	}
	lastMsg, err := s.DB.Queries.GetLatestMessageByRoomID(context.Background(), room.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			detailedRoom.LastActivity = room.CreatedAt
		} else {
			return nil, fmt.Errorf("database failure: %w", err)
		}
	} else {
		detailedRoom.LastActivity = lastMsg.UpdatedAt
		detailedRoom.LastMessagePreview = lastMsg.Content
	}
	users, err := s.DB.Queries.GetUsersInARoom(context.Background(), room.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// do nothing
		} else {
			return nil, fmt.Errorf("database failure: %w", err)
		}
	}
	for _, user := range users {
		detailedUser := s.ToDetailedUser(user)
		unread, err := s.DB.Queries.GetUnreadCountByUserIDRoomID(context.Background(), sqlc.GetUnreadCountByUserIDRoomIDParams{
			UserID: user.ID,
			RoomID: room.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("database failure: %w", err)
		}
		detailedUserUnread := DetailedUserUnread{
			User:   detailedUser,
			Unread: unread,
		}
		detailedRoom.Users = append(detailedRoom.Users, detailedUserUnread)
	}
	return &detailedRoom, nil
}

func (s *DataLayer) ToDetailedUser(user sqlc.User) DetailedUser {
	detailedUser := DetailedUser{
		Username: user.Username,
	}
	if user.DisplayName.Valid {
		detailedUser.DisplayName = &user.DisplayName.String
	}
	return detailedUser
}

func (s *DataLayer) GetRoomCodesByUserID(userID int32) ([]string, error) {
	if userID == 0 {
		return nil, ErrIDZero
	}
	rooms, err := s.DB.Queries.GetRoomsByUser(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("database failure: %w", err)
	}

	roomCodes := []string{}
	for _, room := range rooms {
		roomCodes = append(roomCodes, room.Code)
	}

	return roomCodes, nil
}

func (s *DataLayer) SetUnreadCount(username string, roomCode string, unread int) error {
	if username == "" || roomCode == "" {
		return ErrIDZero
	}
	params := sqlc.SetUnreadCountByUsernameRoomCodeParams{
		Username: username,
		Code:     roomCode,
		Unread:   int32(unread),
	}
	err := s.DB.Queries.SetUnreadCountByUsernameRoomCode(context.Background(), params)
	if err != nil {
		return fmt.Errorf("database failure: %w", err)
	}
	return nil
}

func (s *DataLayer) CreateRoom(userID int32, roomName string, roomDescription *string) (DetailedRoom, error) {
	var err error
	var room sqlc.Room
	var roomCode string
	params := sqlc.CreateRoomParams{
		Name: roomName,
	}
	if roomDescription != nil {
		params.Description.Valid = true
		params.Description.String = *roomDescription
	}
	for {
		randomNumber := rand.Intn(1000)
		roomCode, err = s.sqids.Encode([]uint64{sqidconst, uint64(userID), uint64(randomNumber)})
		if err != nil {
			return DetailedRoom{}, fmt.Errorf("cannot generate room code: %w", err)
		}
		params.Code = roomCode
		if params.Name == "" {
			params.Name = roomCode
		}
		room, err = s.DB.Queries.CreateRoom(context.Background(), params)
		if err != nil {
			pgErr, ok := err.(*pq.Error)
			if ok {
				if pgErr.Code == "23505" { // detect duplication of room code
					continue
				}
			}
		}
		break
	}
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	err = s.DB.Queries.AddUserToRoom(context.Background(), sqlc.AddUserToRoomParams{
		RoomID: room.ID,
		UserID: userID,
		Unread: 0,
	})
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	detailedRoom, err := s.toDetailedRoom(room)
	if err != nil {
		return DetailedRoom{}, err
	}
	return *detailedRoom, nil
}

func (s *DataLayer) LeaveRoom(username string, roomCode string) error {
	if username == "" || roomCode == "" {
		return ErrIDZero
	}
	err := s.DB.Queries.RemoveUserFromRoom(context.Background(), sqlc.RemoveUserFromRoomParams{
		Username: username,
		Code:     roomCode,
	})
	if err != nil {
		return fmt.Errorf("database failure: %w", err)
	}
	return nil
}

func (s *DataLayer) AddUsersToRoom(usernames []string, roomCode string) (DetailedRoom, error) {
	if roomCode == "" || len(usernames) == 0 {
		return DetailedRoom{}, ErrIDZero
	}
	for _, username := range usernames {
		err := s.DB.Queries.AddUserToRoomUsernameRoomCode(context.Background(), sqlc.AddUserToRoomUsernameRoomCodeParams{
			Username: username,
			Code:     roomCode,
		})
		if err != nil {
			return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
		}
	}
	room, err := s.DB.Queries.GetRoomByCode(context.Background(), roomCode)
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	du, err := s.toDetailedRoom(room)
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	return *du, err
}

func (s *DataLayer) JoinRoom(username string, roomCode string) (DetailedRoom, error) {
	if username == "" || roomCode == "" {
		return DetailedRoom{}, ErrIDZero
	}
	err := s.DB.Queries.AddUserToRoomUsernameRoomCode(context.Background(), sqlc.AddUserToRoomUsernameRoomCodeParams{
		Username: username,
		Code:     roomCode,
	})
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	room, err := s.DB.Queries.GetRoomByCode(context.Background(), roomCode)
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	du, err := s.toDetailedRoom(room)
	if err != nil {
		return DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	return *du, err
}
