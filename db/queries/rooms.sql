-- name: CreateRoom :one
INSERT INTO rooms (code, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoomByID :one
SELECT * FROM rooms
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetRoomByCode :one
SELECT * FROM rooms
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: UpdateRoomInfo :exec
UPDATE rooms
SET name = $1,
    description = $2,
    updated_at = NOW()
WHERE id = $3;

-- name: GetRoomsByUser :many
SELECT r.*
FROM rooms r
INNER JOIN rooms_users_memberships rum ON r.id = rum.room_id
WHERE rum.user_id = $1
AND r.deleted_at IS NULL;

-- name: GetUsersInARoom :many
SELECT u.*
FROM users u
INNER JOIN rooms_users_memberships rum ON u.id = rum.user_id
WHERE rum.room_id = $1
AND u.deleted_at IS NULL;

-- name: GetUnreadCountByUserIDRoomID :one
SELECT unread
FROM rooms_users_memberships
WHERE user_id = $1 AND room_id = $2;

-- name: SetUnreadCountByUsernameRoomCode :exec
UPDATE rooms_users_memberships
SET unread = $3
FROM rooms, users
WHERE rooms_users_memberships.room_id = rooms.id
  AND rooms_users_memberships.user_id = users.id
  AND users.username = $1
  AND rooms.code = $2;

-- name: AddUserToRoom :exec
INSERT INTO rooms_users_memberships (room_id, user_id, unread)
VALUES ($1, $2, $3);

-- name: DeleteRoom :exec
UPDATE rooms
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetRoomIDFromRoomCode :one
SELECT id
FROM rooms
WHERE code = $1;