-- +goose Up
-- +goose StatementBegin
CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    active BOOLEAN DEFAULT FALSE,
    turns INTEGER DEFAULT 0,
    wrong INTEGER DEFAULT 0,
    win BOOLEAN DEFAULT 0,
    start INTEGER,
    end INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE games;
-- +goose StatementEnd
