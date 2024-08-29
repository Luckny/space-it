-- name: CreateRequestLog :one
INSERT INTO request_log (id, method, path, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CreateResponseLog :exec
INSERT INTO response_log (id, status)
VALUES ($1, $2);

