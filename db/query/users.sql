-- name: RegisterUser :one
INSERT INTO users (id, email, password)
VALUES ( $1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;
