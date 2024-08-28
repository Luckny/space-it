-- name: CreateSpace :one
INSERT INTO spaces (id, name, owner)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSpaceByID :one
SELECT * FROM spaces
WHERE id = $1 LIMIT 1;

-- name: GetSpaceByName :one
SELECT * FROM spaces
WHERE name = $1 LIMIT 1;

-- name: ListSpaces :many
SELECT * FROM spaces
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateSpace :one
UPDATE spaces
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSpace :exec
DELETE FROM spaces
WHERE id = $1;
