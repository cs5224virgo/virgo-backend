-- name: CreateMessage :one
INSERT INTO messages (content, type, user_id, room_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMessage :exec
UPDATE messages
SET content = $2, type = $3, updated_at = NOW()
WHERE id = $1;

-- name: DeleteMessage :exec
UPDATE messages
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetLatestMessageByRoomID :one
SELECT *
FROM messages
WHERE room_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: GetMessagesByRoomID :many
SELECT *
FROM messages
WHERE room_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2;

-- name: GetMessagesByRoomCode :many
SELECT m.*
FROM messages m
JOIN rooms r ON m.room_id = r.id
WHERE r.code = $1 AND m.deleted_at IS NULL
ORDER BY m.created_at ASC;
