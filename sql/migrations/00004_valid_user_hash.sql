-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN valid_token TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN valid_token;
-- +goose StatementEnd
