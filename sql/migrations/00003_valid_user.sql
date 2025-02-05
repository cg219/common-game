-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN valid INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN valid;
-- +goose StatementEnd
