-- name: SelectTodos :many
SELECT * FROM todo LIMIT $1;

-- name: InsertTodo :one
INSERT INTO todo (body) values ($1) RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todo WHERE id = $1;
