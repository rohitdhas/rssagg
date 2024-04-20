-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeedByUser :one
SELECT * FROM feeds WHERE user_id = $1;