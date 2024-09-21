-- +goose Up
-- +goose StatementBegin
DROP TABLE games;

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    active BOOLEAN DEFAULT FALSE,
    turns INTEGER DEFAULT 0,
    wrong INTEGER DEFAULT 0,
    win BOOLEAN DEFAULT 0,
    start INTEGER,
    end INTEGER
);

CREATE TABLE users_games (
    uid TEXT NOT NULL,
    gid INTEGER NOT NULL,
    PRIMARY KEY(uid, gid)
);
DROP TABLE games;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_games;
DROP TABLE games;

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    active BOOLEAN DEFAULT FALSE,
    turns INTEGER DEFAULT 0,
    wrong INTEGER DEFAULT 0,
    win BOOLEAN DEFAULT 0,
    start INTEGER,
    end INTEGER,
    uid TEXT,
    CONSTRAINT fk_users
    FOREIGN KEY (uid)
    REFERENCES users(uid)
);
-- +goose StatementEnd
