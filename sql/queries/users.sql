-- name: CreateUser :one
INSERT INTO users (id, name, email, api_key)
VALUES ($1, $2, $3, 
    encode(sha256(random()::text::bytea), 'hex')
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;