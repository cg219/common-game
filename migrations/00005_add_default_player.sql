-- +goose Up
-- +goose StatementBegin
INSERT INTO players(id, name)
VALUES(0, "maintester");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM players
WHERE id = 0;
-- +goose StatementEnd
