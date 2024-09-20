-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER,
    key_id TEXT UNIQUE,
    keys TEXT UNIQUE,
    PRIMARY KEY(id, key_id)
);

CREATE TABLE users_players (
    user_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL UNIQUE,
    UNIQUE(user_id, player_id),
    PRIMARY KEY(user_id, player_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE users_players;
-- +goose StatementEnd
