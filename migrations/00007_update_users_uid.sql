-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN uid TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN uid;
-- +goose StatementEnd
