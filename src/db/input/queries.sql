-- name: SelectTodos :many
SELECT *
FROM todos
LIMIT $1;
-- name: InsertTodo :one
INSERT INTO todos (body)
values ($1)
RETURNING *;
-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
-- name: SelectProfileById :one
SELECT *
FROM profiles
WHERE profiles.id = $1;
-- name: SelectEmailOrUsernameAlreadyExists :one
SELECT email
FROM users
WHERE users.email = $1
    OR users.username = $2;
-- name: InsertUser :one
INSERT INTO users (
        email,
        username,
        encrypted_password,
        password_created_at
    )
values ($1, $2, $3, $4)
RETURNING (email, username);
-- name: SelectUserWithEmailPassword :one
SELECT (id)
FROM users
WHERE email = $1
    AND encrypted_password = $2;