-- name: SelectTodos :many
SELECT * FROM todo LIMIT $1;

-- name: InsertTodo :one
INSERT INTO todo (body) values ($1) RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todo WHERE id = $1;

-- name: SelectProfileById :one
SELECT * FROM profiles WHERE profiles.id = $1;

-- name: SelectUsernameFromProfileById :one
SELECT username FROM profiles WHERE profiles.id = $1;

-- name: SelectProfileByUsername :one
SELECT * FROM profiles WHERE profiles.username = $1;

-- name: UpdateProfileUsername :one
UPDATE profiles SET username = $1 WHERE profiles.id = $2 RETURNING *;

-- name: SelectEmailAlreadyExists :one
SELECT email FROM auth.users WHERE auth.users.email = $1;
