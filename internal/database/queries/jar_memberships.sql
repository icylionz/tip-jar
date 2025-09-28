-- name: GetJarMembership :one
SELECT id, jar_id, user_id, role, joined_at
FROM jar_memberships
WHERE jar_id = $1 AND user_id = $2;

-- name: CreateJarMembership :one
INSERT INTO jar_memberships (jar_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING id, jar_id, user_id, role, joined_at;

-- name: ListJarMembers :many
SELECT jm.id, jm.jar_id, jm.user_id, jm.role, jm.joined_at,
       u.email, u.name, u.avatar
FROM jar_memberships jm
INNER JOIN users u ON jm.user_id = u.id
WHERE jm.jar_id = $1
ORDER BY jm.joined_at ASC;

-- name: UpdateMemberRole :one
UPDATE jar_memberships
SET role = $3
WHERE jar_id = $1 AND user_id = $2
RETURNING id, jar_id, user_id, role, joined_at;

-- name: DeleteJarMembership :exec
DELETE FROM jar_memberships
WHERE jar_id = $1 AND user_id = $2;

-- name: IsUserJarMember :one
SELECT EXISTS(
    SELECT 1 FROM jar_memberships
    WHERE jar_id = $1 AND user_id = $2
);

-- name: IsUserJarAdmin :one
SELECT EXISTS(
    SELECT 1 FROM jar_memberships
    WHERE jar_id = $1 AND user_id = $2 AND role = 'admin'
);