// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username, password, display_name
) VALUES (
    $1, $2, $3
)
RETURNING id, created_at, updated_at, deleted_at, username, password, display_name
`

type CreateUserParams struct {
	Username    string
	Password    string
	DisplayName sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password, arg.DisplayName)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Username,
		&i.Password,
		&i.DisplayName,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, created_at, updated_at, deleted_at, username, password, display_name FROM users
WHERE deleted_at IS NULL
ORDER BY id
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
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

const getUserByID = `-- name: GetUserByID :one
SELECT id, created_at, updated_at, deleted_at, username, password, display_name FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Username,
		&i.Password,
		&i.DisplayName,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, created_at, updated_at, deleted_at, username, password, display_name FROM users
WHERE username = $1 AND deleted_at IS NULL LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Username,
		&i.Password,
		&i.DisplayName,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET (
    username, password, display_name
) = (
    $2, $3, $4
)
WHERE id = $1 AND deleted_at IS NULL
`

type UpdateUserParams struct {
	ID          int32
	Username    string
	Password    string
	DisplayName sql.NullString
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.ID,
		arg.Username,
		arg.Password,
		arg.DisplayName,
	)
	return err
}
