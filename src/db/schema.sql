CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY NOT NULL,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_salt TEXT NOT NULL,
    encrypted_password TEXT NOT NULL,
    password_created_at TIMESTAMP NOT NULL,
    global_chat_color TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    jwt_type TEXT NOT NULL,
    jwt TEXT NOT NULL,
    valid BOOLEAN NOT NULL,
    family TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users_tokens (
    user_id TEXT NOT NULL,
    token_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (token_id) REFERENCES tokens(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, token_id)
);

CREATE TABLE IF NOT EXISTS chatrooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    creator TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
     
);

CREATE TABLE IF NOT EXISTS chatroom_members (
    chatroom_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    chatroom_color TEXT NOT NULL,
    PRIMARY KEY (chatroom_id, user_id),
    FOREIGN KEY (chatroom_id) REFERENCES chatrooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id TEXT, -- nullable for anon messages
    author_username TEXT NOT NULL,
    message TEXT NOT NULL,
    room_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES chatrooms(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
    -- No foreign key constraint for author_username in case they change their username, we don't have to rewrite the messages they sent under a different username with a join.
);

CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    body TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
