-- table: todos
-- insert/select/update to be done
-- name: SelectUserTodos :many
SELECT *
FROM todos
WHERE author = $1;
-- name: InsertTodo :one
INSERT INTO todos (body, author)
VALUES ($1, $2)
RETURNING *;
-- name: UpdateTodo :one
UPDATE todos
SET body = $1
WHERE id = $2
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
-- table: tokens
-- name: InsertToken :exec
INSERT INTO tokens (jwt_type, jwt, valid)
VALUES ($1, $2, $3);
-- name: UpdateTokenValid :exec
UPDATE tokens
SET valid = $1
WHERE jwt = $2;
-- name: DeleteTokensByUserId :exec
DELETE FROM tokens
WHERE tokens.id = (
        SELECT token_id
        FROM users_tokens
        WHERE users_tokens.user_id = $1
    );
-- name: SelectUsersTokens :many
SELECT *
FROM users_tokens
WHERE user_id = $1;
-- name: IsValidToken :one
SELECT valid
FROM tokens
WHERE jwt = $1;
-- table: users_tokens
-- name: InsertUsersTokens :exec
INSERT INTO users_tokens (user_id, token_id)
VALUES ($1, $2);