-- +goose Up
-- +goose StatementBegin
CREATE TABLE players (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE players;
-- +goose StatementEnd
