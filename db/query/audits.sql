-- name: CreateAuthenticatedRequestLog :one
INSERT INTO request_log (method, path, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateUnauthenticatedRequestLog :one
INSERT INTO request_log (method, path)
VALUES ($1, $2)
RETURNING *;

-- name: CreateResponseLog :one
INSERT INTO response_log (id, status)
VALUES ($1, $2)
RETURNING *;

