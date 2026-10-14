-- name: CreateItem :one
INSERT INTO items (title) VALUES (?) RETURNING *;

-- name: GetItem :one
SELECT * FROM items WHERE id = ?;

-- name: ListItems :many
SELECT * FROM items ORDER BY id DESC;
