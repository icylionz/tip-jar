-- name: GetTipJar :one
SELECT id, name, description, invite_code, created_by, created_at, updated_at
FROM tip_jars
WHERE id = $1;

-- name: GetTipJarByInviteCode :one
SELECT id, name, description, invite_code, created_by, created_at, updated_at
FROM tip_jars
WHERE invite_code = $1;

-- name: CreateTipJar :one
INSERT INTO tip_jars (name, description, invite_code, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, invite_code, created_by, created_at, updated_at;

-- name: ListTipJarsForUser :many
SELECT tj.id, tj.name, tj.description, tj.invite_code, tj.created_by, tj.created_at, tj.updated_at
FROM tip_jars tj
INNER JOIN jar_memberships jm ON tj.id = jm.jar_id
WHERE jm.user_id = $1
ORDER BY tj.created_at DESC;

-- name: UpdateTipJar :one
UPDATE tip_jars
SET name = $2, description = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, invite_code, created_by, created_at, updated_at;

-- name: DeleteTipJar :exec
DELETE FROM tip_jars
WHERE id = $1;