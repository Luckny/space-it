-- name: CreatePermission :one
INSERT INTO permissions (user_id, space_id, read_permission, write_permission, delete_permission)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateReadPermission :one
INSERT INTO permissions (user_id, space_id, read_permission)
VALUES ($1, $2, true)
RETURNING *;

-- name: CreateWritePermission :one
INSERT INTO permissions (user_id, space_id, write_permission)
VALUES ($1, $2, true)
RETURNING *;

-- name: CreateDeletePermission :one
INSERT INTO permissions (user_id, space_id, delete_permission)
VALUES ($1, $2, true)
RETURNING *;

-- name: CreateAllPermission :one
INSERT INTO permissions (user_id, space_id, read_permission, write_permission, delete_permission)
VALUES ($1, $2, true, true, true)
RETURNING *;

-- name: GetPermissionsByUserAndSpaceID :one
SELECT * FROM permissions
WHERE user_id = $1
AND space_id = $2
LIMIT 1;

