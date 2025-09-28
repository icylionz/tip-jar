-- name: GetUserByGoogleID :one
SELECT id, email, name, avatar, google_id, created_at, updated_at 
FROM users 
WHERE google_id = $1;

-- name: GetUserByID :one
SELECT id, email, name, avatar, google_id, created_at, updated_at 
FROM users 
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, name, avatar, google_id)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name, avatar, google_id, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users 
SET name = $2, avatar = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, email, name, avatar, google_id, created_at, updated_at;

-- name: ListUsers :many
SELECT id, email, name, avatar, google_id, created_at, updated_at 
FROM users
ORDER BY created_at DESC;