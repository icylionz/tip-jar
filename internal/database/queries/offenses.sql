-- name: GetOffense :one
SELECT id, jar_id, offense_type_id, reporter_id, offender_id, notes, cost_override, status, created_at, updated_at
FROM offenses
WHERE id = $1;

-- name: ListOffensesForJar :many
SELECT o.id, o.jar_id, o.offense_type_id, o.reporter_id, o.offender_id, o.notes, o.cost_override, o.status, o.created_at, o.updated_at,
       ot.name as offense_type_name, ot.cost_amount, ot.cost_unit,
       reporter.name as reporter_name, offender.name as offender_name
FROM offenses o
INNER JOIN offense_types ot ON o.offense_type_id = ot.id
INNER JOIN users reporter ON o.reporter_id = reporter.id
INNER JOIN users offender ON o.offender_id = offender.id
WHERE o.jar_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPendingOffensesForUser :many
SELECT o.id, o.jar_id, o.offense_type_id, o.reporter_id, o.offender_id, o.notes, o.cost_override, o.status, o.created_at, o.updated_at,
       ot.name as offense_type_name, ot.cost_amount, ot.cost_unit,
       tj.name as jar_name
FROM offenses o
INNER JOIN offense_types ot ON o.offense_type_id = ot.id
INNER JOIN tip_jars tj ON o.jar_id = tj.id
WHERE o.offender_id = $1 AND o.status = 'pending'
ORDER BY o.created_at DESC;

-- name: CreateOffense :one
INSERT INTO offenses (jar_id, offense_type_id, reporter_id, offender_id, notes, cost_override)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, jar_id, offense_type_id, reporter_id, offender_id, notes, cost_override, status, created_at, updated_at;

-- name: UpdateOffenseStatus :one
UPDATE offenses
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, jar_id, offense_type_id, reporter_id, offender_id, notes, cost_override, status, created_at, updated_at;

-- name: GetUserBalanceInJar :one
SELECT 
    COALESCE(SUM(
        CASE 
            WHEN o.cost_override IS NOT NULL THEN o.cost_override
            ELSE ot.cost_amount
        END
    ), 0) as total_owed
FROM offenses o
INNER JOIN offense_types ot ON o.offense_type_id = ot.id
WHERE o.jar_id = $1 AND o.offender_id = $2 AND o.status = 'pending';

-- name: GetUserBalancesByUnitInJar :many
SELECT 
    COALESCE(ot.cost_unit, 'items') as unit,
    COALESCE(SUM(
        CASE 
            WHEN o.cost_override IS NOT NULL THEN o.cost_override
            ELSE ot.cost_amount
        END
    ), 0) as total_owed,
    COUNT(*) as offense_count
FROM offenses o
INNER JOIN offense_types ot ON o.offense_type_id = ot.id
WHERE o.jar_id = $1 AND o.offender_id = $2 AND o.status = 'pending'
GROUP BY ot.cost_unit
ORDER BY total_owed DESC;

-- name: GetJarBalancesByUnit :many
SELECT 
    u.id as user_id,
    u.name as user_name,
    u.avatar,
    COALESCE(ot.cost_unit, 'items') as unit,
    COALESCE(SUM(
        CASE 
            WHEN o.cost_override IS NOT NULL THEN o.cost_override
            ELSE ot.cost_amount
        END
    ), 0) as total_owed,
    COUNT(*) as offense_count
FROM users u
INNER JOIN jar_memberships jm ON u.id = jm.user_id
LEFT JOIN offenses o ON u.id = o.offender_id AND o.jar_id = $1 AND o.status = 'pending'
LEFT JOIN offense_types ot ON o.offense_type_id = ot.id
WHERE jm.jar_id = $1
GROUP BY u.id, u.name, u.avatar, ot.cost_unit
HAVING COUNT(o.id) > 0 OR ot.cost_unit IS NULL
ORDER BY u.name, total_owed DESC;