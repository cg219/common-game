-- +goose Up
-- +goose StatementBegin
CREATE TABLE subjects (
    id INTEGER PRIMARY KEY,
    subject TEXT NOT NULL,
    words TEXT NOT NULL,
    used INTEGER DEFAULT 0,
    correct INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subjects;
-- +goose StatementEnd
