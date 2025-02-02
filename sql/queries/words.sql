-- name: GetWords :many
SELECT id, word
FROM words;

-- name: GetWordByValue :one
SELECT id, word
FROM words
WHERE value = ?
LIMIT 1;

-- name: SaveWord :exec
INSERT INTO words (value, word)
VALUES(?, ?)
RETURNING id;
