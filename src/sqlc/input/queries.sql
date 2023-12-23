-- name: GetTodosWithLimit :many
SELECT * FROM todo LIMIT $1;

-- name: InsertTodo :one
INSERT INTO todo () values ($1) RETURNING *

