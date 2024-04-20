package datalayer

import (
	"context"
	"fmt"
	"time"

	sqlc "github.com/cs5224virgo/virgo/db/generated"
)

type DetailedMessage struct {
	CreatedAt time.Time        `json:"createdAt"`
	Content   string           `json:"content"`
	Type      sqlc.MessageType `json:"type"`
	User      DetailedUser     `json:"user"`
	RoomCode  string           `json:"roomCode"`
}

func (s *DataLayer) GetAllMessagesForRoom(roomCode string) ([]DetailedMessage, error) {
	if roomCode == "" {
		return nil, ErrIDZero
	}
	msgs, err := s.DB.Queries.GetMessagesByRoomCode(context.Background(), roomCode)
	if err != nil {
		return nil, fmt.Errorf("database failure: %w", err)
	}
	ret := []DetailedMessage{}
	for _, msg := range msgs {
		du, err := s.toDetailedMessage(msg)
		if err != nil {
			return nil, fmt.Errorf("database failure: %w", err)
		}
		ret = append(ret, *du)
	}
	return ret, nil
}

func (s *DataLayer) GetLastWeekMessagesForRoom(roomCode string) ([]DetailedMessage, error) {
	if roomCode == "" {
		return nil, ErrIDZero
	}
	lastWeek := time.Now().AddDate(0, 0, -7)
	msgs, err := s.DB.Queries.GetMessagesAfterTimeByRoomCode(context.Background(), sqlc.GetMessagesAfterTimeByRoomCodeParams{
		Code:      roomCode,
		CreatedAt: lastWeek,
	})
	if err != nil {
		return nil, fmt.Errorf("database failure: %w", err)
	}
	ret := []DetailedMessage{}
	for _, msg := range msgs {
		du, err := s.toDetailedMessage(msg)
		if err != nil {
			return nil, fmt.Errorf("database failure: %w", err)
		}
		ret = append(ret, *du)
	}
	return ret, nil
}

func (s *DataLayer) toDetailedMessage(msg sqlc.Message) (*DetailedMessage, error) {
	room, err := s.DB.Queries.GetRoomByID(context.Background(), msg.RoomID)
	if err != nil {
		return nil, fmt.Errorf("database failure: %w", err)
	}
	user, err := s.DB.Queries.GetUserByID(context.Background(), msg.UserID)
	if err != nil {
		return nil, fmt.Errorf("database failure: %w", err)
	}
	du := s.ToDetailedUser(user)
	dm := DetailedMessage{
		CreatedAt: msg.CreatedAt,
		Content:   msg.Content,
		Type:      msg.Type,
		User:      du,
		RoomCode:  room.Code,
	}
	return &dm, nil
}

func (s *DataLayer) NewMessage(content string, roomCode string, username string) (DetailedMessage, DetailedRoom, error) {
	if content == "" || roomCode == "" || username == "" {
		return DetailedMessage{}, DetailedRoom{}, ErrIDZero
	}
	user, err := s.DB.Queries.GetUserByUsername(context.Background(), username)
	if err != nil {
		return DetailedMessage{}, DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	roomID, err := s.DB.Queries.GetRoomIDFromRoomCode(context.Background(), roomCode)
	if err != nil {
		return DetailedMessage{}, DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	msg, err := s.DB.Queries.CreateMessage(context.Background(), sqlc.CreateMessageParams{
		Type:    sqlc.MessageTypeNormal,
		Content: content,
		UserID:  user.ID,
		RoomID:  roomID,
	})
	if err != nil {
		return DetailedMessage{}, DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	dm, err := s.toDetailedMessage(msg)
	if err != nil {
		return DetailedMessage{}, DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	room, err := s.DB.Queries.GetRoomByID(context.Background(), roomID)
	if err != nil {
		return DetailedMessage{}, DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	dr, err := s.toDetailedRoom(room)
	if err != nil {
		return DetailedMessage{}, DetailedRoom{}, fmt.Errorf("database failure: %w", err)
	}
	return *dm, *dr, nil
}
