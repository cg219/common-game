-- name: SaveWord :exec
INSERT INTO words (word)
VALUES(?);

-- name: GetWords :many
SELECT word
FROM words
