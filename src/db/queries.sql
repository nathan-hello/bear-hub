-- table: todos
-- name: SelectTodosByUsername :many
SELECT * FROM todos WHERE username = ?;
-- name: InsertTodo :one
INSERT INTO todos (body, username, created_at) VALUES (?, ?, ?) RETURNING *;
-- name: UpdateTodo :one
UPDATE todos SET body = ? WHERE id = ? RETURNING *;
-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;

-- table: users
-- name: InsertUser :one
INSERT INTO users (id, email, username, password_salt, encrypted_password, password_created_at)
VALUES (?, ?, ?, ?, ?, ?) RETURNING id, email, username;
-- name: SelectUserByEmail :one
SELECT * FROM users WHERE email = ?;
-- name: SelectUserByUsername :one
SELECT * FROM users WHERE username = ?;
-- name: SelectEmailOrUsernameAlreadyExists :one
SELECT email FROM users WHERE email = ? OR username = ?;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- table: tokens
-- name: SelectTokenFromId :one
SELECT * FROM tokens WHERE id = ?;
-- name: SelectTokenFromJwtString :one
SELECT * FROM tokens WHERE jwt = ?;
-- name: InsertToken :one
INSERT INTO tokens (jwt_type, jwt, valid, family) VALUES (?, ?, ?, ?) RETURNING *;
-- name: UpdateTokenValid :one
UPDATE tokens SET valid = ? WHERE jwt = ? RETURNING id;
-- name: UpdateUserTokensToInvalid :exec
UPDATE tokens SET valid = FALSE WHERE id IN (
        SELECT token_id FROM users_tokens WHERE user_id = ?
    );
-- name: UpdateTokensFamilyInvalid :exec
UPDATE tokens SET valid = FALSE WHERE family = ?;
-- name: DeleteTokensByUserId :exec
DELETE FROM tokens WHERE id IN (
        SELECT token_id FROM users_tokens WHERE user_id = ?
    );

-- table: users_tokens
-- name: SelectUsersTokens :many
SELECT * FROM users_tokens WHERE user_id = ?;
-- name: SelectUserIdFromToken :one
SELECT user_id FROM users_tokens WHERE token_id = ? LIMIT 1;
-- name: InsertUsersTokens :exec
INSERT INTO users_tokens (user_id, token_id) VALUES (?, ?);

-- table: chatrooms
-- name: SelectChatrooms :many
SELECT * FROM chatrooms ORDER BY created_at DESC LIMIT ?;
-- name: InsertChatroom :one
INSERT INTO chatrooms (name, creator, created_at) VALUES (?, ?, ?) RETURNING id;
-- name: DeleteChatroom :exec
DELETE FROM chatrooms WHERE id = ?;
-- name: UpdateChatroomName :one
UPDATE chatrooms SET name = ? WHERE id = ? RETURNING *;

-- table: messages
-- name: SelectMessagesByChatroom :many
SELECT * FROM messages WHERE room_id = ? ORDER BY created_at DESC LIMIT ?;
-- name: SelectMessagesByUser :many
SELECT * FROM messages WHERE author = ? ORDER BY created_at DESC LIMIT ?;
-- name: InsertMessage :exec
INSERT INTO messages (author, message, color, room_id, created_at) VALUES (?, ?, ?, ?, ?);
-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ?;
-- name: UpdateMessage :one
UPDATE messages SET message = ? WHERE id = ? RETURNING *;
