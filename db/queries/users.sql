-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
    username, password, display_name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET (
    username, password, display_name
) = (
    $2, $3, $4
)
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;