-- +goose Up
-- +goose StatementBegin
DROP TABLE games;
DROP TABLE players;
DROP TABLE users_players;

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

-- +goose Down
-- +goose StatementBegin
DROP TABLE games;

CREATE TABLE players (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    active BOOLEAN DEFAULT FALSE,
    turns INTEGER DEFAULT 0,
    wrong INTEGER DEFAULT 0,
    win BOOLEAN DEFAULT 0,
    start INTEGER,
    end INTEGER,
    player_id INTEGER,
    CONSTRAINT fk_players
    FOREIGN KEY (player_id)
    REFERENCES players(id)
);

CREATE TABLE users_players (
    user_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL UNIQUE,
    UNIQUE(user_id, player_id),
    PRIMARY KEY(user_id, player_id)
);
-- +goose StatementEnd
