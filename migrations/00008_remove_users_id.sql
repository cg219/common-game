-- +goose Up
-- +goose StatementBegin
DROP TABLE users;

CREATE TABLE users (
    uid TEXT,
    key_id TEXT UNIQUE,
    keys TEXT UNIQUE,
    PRIMARY KEY(uid, key_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

CREATE TABLE users (
    id INTEGER,
    key_id TEXT UNIQUE,
    keys TEXT UNIQUE,
    PRIMARY KEY(id, key_id)
);
-- +goose StatementEnd
