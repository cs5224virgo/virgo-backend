// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rooms.sql

package sqlc

import (
	"context"
	"database/sql"
)

const addUserToRoom = `-- name: AddUserToRoom :exec
INSERT INTO rooms_users_memberships (room_id, user_id, unread)
VALUES ($1, $2, $3)
`

type AddUserToRoomParams struct {
	RoomID int32
	UserID int32
	Unread int32
}

func (q *Queries) AddUserToRoom(ctx context.Context, arg AddUserToRoomParams) error {
	_, err := q.db.ExecContext(ctx, addUserToRoom, arg.RoomID, arg.UserID, arg.Unread)
	return err
}

const addUserToRoomUsernameRoomCode = `-- name: AddUserToRoomUsernameRoomCode :exec
INSERT INTO rooms_users_memberships (room_id, user_id, unread)
SELECT r.id, u.id, 0
FROM rooms r, users u
WHERE r.code = $2 AND u.username = $1
`

type AddUserToRoomUsernameRoomCodeParams struct {
	Username string
	Code     string
}

func (q *Queries) AddUserToRoomUsernameRoomCode(ctx context.Context, arg AddUserToRoomUsernameRoomCodeParams) error {
	_, err := q.db.ExecContext(ctx, addUserToRoomUsernameRoomCode, arg.Username, arg.Code)
	return err
}

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (code, name, description)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at, deleted_at, code, name, description
`

type CreateRoomParams struct {
	Code        string
	Name        string
	Description sql.NullString
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom, arg.Code, arg.Name, arg.Description)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Code,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const deleteRoom = `-- name: DeleteRoom :exec
UPDATE rooms
SET deleted_at = NOW()
WHERE id = $1
`

func (q *Queries) DeleteRoom(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteRoom, id)
	return err
}

const getRoomByCode = `-- name: GetRoomByCode :one
SELECT id, created_at, updated_at, deleted_at, code, name, description FROM rooms
WHERE code = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetRoomByCode(ctx context.Context, code string) (Room, error) {
	row := q.db.QueryRowContext(ctx, getRoomByCode, code)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Code,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const getRoomByID = `-- name: GetRoomByID :one
SELECT id, created_at, updated_at, deleted_at, code, name, description FROM rooms
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetRoomByID(ctx context.Context, id int32) (Room, error) {
	row := q.db.QueryRowContext(ctx, getRoomByID, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Code,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const getRoomIDFromRoomCode = `-- name: GetRoomIDFromRoomCode :one
SELECT id
FROM rooms
WHERE code = $1
`

func (q *Queries) GetRoomIDFromRoomCode(ctx context.Context, code string) (int32, error) {
	row := q.db.QueryRowContext(ctx, getRoomIDFromRoomCode, code)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getRoomsByUser = `-- name: GetRoomsByUser :many
SELECT r.id, r.created_at, r.updated_at, r.deleted_at, r.code, r.name, r.description
FROM rooms r
INNER JOIN rooms_users_memberships rum ON r.id = rum.room_id
WHERE rum.user_id = $1
AND r.deleted_at IS NULL
`

func (q *Queries) GetRoomsByUser(ctx context.Context, userID int32) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getRoomsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Code,
			&i.Name,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnreadCountByUserIDRoomID = `-- name: GetUnreadCountByUserIDRoomID :one
SELECT unread
FROM rooms_users_memberships
WHERE user_id = $1 AND room_id = $2
`

type GetUnreadCountByUserIDRoomIDParams struct {
	UserID int32
	RoomID int32
}

func (q *Queries) GetUnreadCountByUserIDRoomID(ctx context.Context, arg GetUnreadCountByUserIDRoomIDParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, getUnreadCountByUserIDRoomID, arg.UserID, arg.RoomID)
	var unread int32
	err := row.Scan(&unread)
	return unread, err
}

const getUsersInARoom = `-- name: GetUsersInARoom :many
SELECT u.id, u.created_at, u.updated_at, u.deleted_at, u.username, u.password, u.display_name
FROM users u
INNER JOIN rooms_users_memberships rum ON u.id = rum.user_id
WHERE rum.room_id = $1
AND u.deleted_at IS NULL
`

func (q *Queries) GetUsersInARoom(ctx context.Context, roomID int32) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersInARoom, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Username,
			&i.Password,
			&i.DisplayName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeUserFromRoom = `-- name: RemoveUserFromRoom :exec
DELETE FROM rooms_users_memberships
WHERE user_id = (SELECT id FROM users WHERE username = $1)
  AND room_id = (SELECT id FROM rooms WHERE code = $2)
`

type RemoveUserFromRoomParams struct {
	Username string
	Code     string
}

func (q *Queries) RemoveUserFromRoom(ctx context.Context, arg RemoveUserFromRoomParams) error {
	_, err := q.db.ExecContext(ctx, removeUserFromRoom, arg.Username, arg.Code)
	return err
}

const setUnreadCountByUsernameRoomCode = `-- name: SetUnreadCountByUsernameRoomCode :exec
UPDATE rooms_users_memberships
SET unread = $3
FROM rooms, users
WHERE rooms_users_memberships.room_id = rooms.id
  AND rooms_users_memberships.user_id = users.id
  AND users.username = $1
  AND rooms.code = $2
`

type SetUnreadCountByUsernameRoomCodeParams struct {
	Username string
	Code     string
	Unread   int32
}

func (q *Queries) SetUnreadCountByUsernameRoomCode(ctx context.Context, arg SetUnreadCountByUsernameRoomCodeParams) error {
	_, err := q.db.ExecContext(ctx, setUnreadCountByUsernameRoomCode, arg.Username, arg.Code, arg.Unread)
	return err
}

const updateRoomInfo = `-- name: UpdateRoomInfo :exec
UPDATE rooms
SET name = $1,
    description = $2,
    updated_at = NOW()
WHERE id = $3
`

type UpdateRoomInfoParams struct {
	Name        string
	Description sql.NullString
	ID          int32
}

func (q *Queries) UpdateRoomInfo(ctx context.Context, arg UpdateRoomInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateRoomInfo, arg.Name, arg.Description, arg.ID)
	return err
}
