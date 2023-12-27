-- name: SelectTodos :many
SELECT * FROM todo LIMIT $1;

-- name: InsertTodo :one
INSERT INTO todo (body) values ($1) RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todo WHERE id = $1;

-- name: InsertProfile :one
INSERT INTO profile (user_id, username) values ($1, $2) RETURNING *;

-- name: SelectProfileById :one
SELECT * FROM profile WHERE profile.id = $1;

-- name: SelectProfileByUsername :one
SELECT * FROM profile WHERE profile.username = $1;

-- name: SelectEmailAlreadyExists :one
SELECT email FROM auth.users WHERE auth.users.email = $1;
