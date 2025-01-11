-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT NOT NULL,
    reset TEXT,
    reset_time INTEGER
);

CREATE TABLE sessions (
    accessToken TEXT NOT NULL,
    refreshToken TEXT NOT NULL,
    valid INTEGER DEFAULT 1,
    UNIQUE(accessToken, refreshToken),
    PRIMARY KEY(accessToken, refreshToken)
);

CREATE TABLE apikeys (
    key TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    uid INTEGER,
    CONSTRAINT fk_users
    FOREIGN KEY(uid)
    REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
