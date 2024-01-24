-- table: todos
-- name: SelectTodosByUsername :many
SELECT * FROM todos WHERE username = $1;
-- name: InsertTodo :one
INSERT INTO todos (body, username) VALUES ($1, $2) RETURNING *;
-- name: UpdateTodo :one
UPDATE todos SET body = $1 WHERE id = $2 RETURNING *;
-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;


-- table: users
-- name: InsertUser :one
INSERT INTO users ( email, username, encrypted_password, password_created_at)
values ($1, $2, $3, $4)
RETURNING id, email, username;
-- name: SelectUserByEmail :one
SELECT * FROM users WHERE email = $1;
-- name: SelectUserByUsername :one
SELECT * FROM users WHERE username = $1;
-- name: SelectEmailOrUsernameAlreadyExists :one
SELECT email FROM users WHERE users.email = $1 OR users.username = $2;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1; 

-- table: tokens
-- name: SelectTokenFromId :one
SELECT * FROM tokens WHERE id = $1;
-- name: SelectTokenFromJwtString :one
SELECT * FROM tokens WHERE jwt = $1;
-- name: InsertToken :one
INSERT INTO tokens (jwt_type, jwt, valid, family) VALUES ($1, $2, $3, $4) RETURNING id;
-- name: UpdateTokenValid :one
UPDATE tokens SET valid = $1 WHERE jwt = $2 RETURNING id;
-- name: UpdateUserTokensToInvalid :exec
UPDATE tokens SET valid = FALSE FROM users_tokens
INNER JOIN tokens AS t ON users_tokens.token_id = t.id
    WHERE users_tokens.user_id = $1
    AND tokens.id = t.id;
-- name: UpdateTokensFamilyInvalid :exec
UPDATE tokens 
SET valid = FALSE 
WHERE family = $1;
-- name: DeleteTokensByUserId :exec
DELETE FROM tokens
WHERE tokens.id IN (
        SELECT token_id FROM users_tokens WHERE users_tokens.user_id = $1
    );


-- table: users_tokens
-- name: SelectUsersTokens :many
SELECT * FROM users_tokens WHERE user_id = $1;
-- name: SelectUserIdFromToken :one
SELECT user_id FROM users_tokens WHERE token_id = $1 LIMIT 1;
-- name: InsertUsersTokens :exec
INSERT INTO users_tokens (user_id, token_id) VALUES ($1, $2); 
