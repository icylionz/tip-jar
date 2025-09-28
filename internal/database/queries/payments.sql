-- name: GetPayment :one
SELECT id, offense_id, user_id, amount, proof_type, proof_url, verified, verified_by, created_at, updated_at
FROM payments
WHERE id = $1;

-- name: ListPaymentsForOffense :many
SELECT id, offense_id, user_id, amount, proof_type, proof_url, verified, verified_by, created_at, updated_at
FROM payments
WHERE offense_id = $1
ORDER BY created_at DESC;

-- name: CreatePayment :one
INSERT INTO payments (offense_id, user_id, amount, proof_type, proof_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, offense_id, user_id, amount, proof_type, proof_url, verified, verified_by, created_at, updated_at;

-- name: VerifyPayment :one
UPDATE payments
SET verified = true, verified_by = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, offense_id, user_id, amount, proof_type, proof_url, verified, verified_by, created_at, updated_at;

-- name: ListPaymentsForUser :many
SELECT p.id, p.offense_id, p.user_id, p.amount, p.proof_type, p.proof_url, p.verified, p.verified_by, p.created_at, p.updated_at,
       o.jar_id, tj.name as jar_name, ot.name as offense_type_name
FROM payments p
INNER JOIN offenses o ON p.offense_id = o.id
INNER JOIN tip_jars tj ON o.jar_id = tj.id
INNER JOIN offense_types ot ON o.offense_type_id = ot.id
WHERE p.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;