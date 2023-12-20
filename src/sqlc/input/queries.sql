-- name: GetTodosWithLimit :many
SELECT * FROM todo LIMIT $1;


