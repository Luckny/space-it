-- name: CreateToken :one
INSERT INTO tokens (user_id, expiry, attributes)
VALUES ( $1, $2, $3)
RETURNING *;

-- name: GetToken :one
SELECT * FROM tokens
WHERE token_id = $1 LIMIT 1;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE token_id = $1;
