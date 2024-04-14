package datalayer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DetailedUser struct {
	Username    string
	DisplayName *string
}

type DetailedRoom struct {
	ID                 int32          `json:"id"`
	CreatedAt          time.Time      `json:"createdAt"`
	Code               string         `json:"code"`
	Name               string         `json:"name"`
	Description        *string        `json:"description"`
	LastActivity       time.Time      `json:"lastActivity"`
	LastMessagePreview string         `json:"lastMessagePreview"`
	Users              []DetailedUser `json:"users"`
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

			} else {
				return nil, fmt.Errorf("database failure: %w", err)
			}
		}
		for _, user := range users {
			detailedUser := DetailedUser{
				Username: user.Username,
			}
			if user.DisplayName.Valid {
				detailedUser.DisplayName = &user.DisplayName.String
			}
			detailedRoom.Users = append(detailedRoom.Users, detailedUser)
		}
		detailedRooms = append(detailedRooms, detailedRoom)
	}
	return detailedRooms, nil
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
