-- name: GetOffense :one
SELECT id, jar_id, offense_type_id, reporter_id, offender_id, notes, cost_override, status, created_at, updated_at
FROM offenses
WHERE id = $1;

-- name: ListOffensesForJar :many
SELECT o.id, o.jar_id, o.offense_type_id, o.reporter_id, o.offender_id, o.notes, o.cost_override, o.status, o.created_at, o.updated_at,
       ot.name as offense_type_name, ot.cost_type, ot.cost_amount, ot.cost_action,
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
       ot.name as offense_type_name, ot.cost_type, ot.cost_amount, ot.cost_action,
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
            WHEN ot.cost_type = 'monetary' THEN 
                CASE 
                    WHEN o.cost_override IS NOT NULL THEN o.cost_override
                    ELSE ot.cost_amount
                END
            ELSE 0
        END
    ), 0) as total_owed
FROM offenses o
INNER JOIN offense_types ot ON o.offense_type_id = ot.id
WHERE o.jar_id = $1 AND o.offender_id = $2 AND o.status = 'pending';