-- +goose Up
-- +goose StatementBegin
CREATE TABLE bugreports (
    id INTEGER PRIMARY KEY,
    problem TEXT NOT NULL,
    result TEXT NOT NULL,
    steps TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    uid INTEGER NOT NULL,
    CONSTRAINT fk_users_report
    FOREIGN KEY(uid)
    REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bugreports;
-- +goose StatementEnd
