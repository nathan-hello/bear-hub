-- table: todos
-- name: SelectAllTodos :many
SELECT *
FROM todos;
-- name: SelectTodosByIds :many
SELECT *
FROM todos
WHERE id = ANY($1::int [])
LIMIT $2;
-- name: InsertTodo :one
INSERT INTO todos (body, author)
values ($1, $2)
RETURNING *;
-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
-- table: users
-- name: InsertUser :one
INSERT INTO users (
        email,
        username,
        encrypted_password,
        password_created_at
    )
values ($1, $2, $3, $4)
RETURNING id,
    email,
    username;
-- name: SelectUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: SelectUserByUsername :one
SELECT *
FROM users
WHERE username = $1;
-- name: SelectEmailOrUsernameAlreadyExists :one
SELECT email
FROM users
WHERE users.email = $1
    OR users.username = $2;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
-- table: profiles
-- name: SelectProfileById :one
SELECT *
FROM profiles
WHERE profiles.id = $1;
-- name: SelectProfileByUsername :one
SELECT *
FROM profiles
    INNER JOIN users ON profiles.id = users.id
WHERE users.username = $1;
-- name: InsertProfile :one
INSERT INTO profiles (id)
values ($1)
returning (id);
-- name: DeleteProfile :exec
DELETE from profiles
where id = $1;