-- name: GetOffenseType :one
SELECT id, jar_id, name, description, cost_amount, cost_unit, is_active, created_at, updated_at
FROM offense_types
WHERE id = $1;

-- name: ListOffenseTypesForJar :many
SELECT id, jar_id, name, description, cost_amount, cost_unit, is_active, created_at, updated_at
FROM offense_types
WHERE jar_id = $1 AND is_active = true
ORDER BY name ASC;

-- name: ListAllOffenseTypesForJar :many
SELECT id, jar_id, name, description, cost_amount, cost_unit, is_active, created_at, updated_at
FROM offense_types
WHERE jar_id = $1
ORDER BY created_at DESC;

-- name: CreateOffenseType :one
INSERT INTO offense_types (jar_id, name, description, cost_amount, cost_unit)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, jar_id, name, description, cost_amount, cost_unit, is_active, created_at, updated_at;

-- name: UpdateOffenseType :one
UPDATE offense_types
SET name = $2, description = $3, cost_amount = $4, cost_unit = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, jar_id, name, description, cost_amount, cost_unit, is_active, created_at, updated_at;

-- name: DeactivateOffenseType :one
UPDATE offense_types
SET is_active = false, updated_at = NOW()
WHERE id = $1
RETURNING id, jar_id, name, description, cost_amount, cost_unit, is_active, created_at, updated_at;