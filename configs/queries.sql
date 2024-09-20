-- name: GetWords :many
SELECT word
FROM words;

-- name: SaveWord :exec
INSERT INTO words (word)
VALUES(?);

-- name: SaveSubject :exec
INSERT INTO subjects(subject, words)
VALUES (?, ?);

-- name: GetSubjects :many
SELECT subject, words FROM subjects;

-- name: GetSubjectsForGame :many
WITH randomrows AS (
    SELECT id, subject, words
    FROM subjects
    ORDER BY random()
    LIMIT 4
),
single_words AS (
    SELECT randomrows.id, json_extract(json_each.value) AS stringv
    FROM randomrows, json_each(randomrows.words)
)
SELECT id, subject, words
FROM randomrows
WHERE NOT EXISTS (
    SELECT 1
    FROM single_words l1
    JOIN single_words l2 ON l1.stringv = l2.stringv AND l1.id != l2.id
    WHERE l1.id = randomrows.id
);

-- name: GetPlayerByName :one
SELECT id, name
FROM players
WHERE name = ?;

-- name: SavePlayer :exec
INSERT INTO players (name)
VALUES (?);

-- name: GetPlayerById :one
SELECT id, name
FROM players
WHERE id = ?;

-- name: SaveNewGame :one
INSERT INTO games(active, start, player_id)
VALUES (?, ?, ?)
RETURNING id;

-- name: UpdateGame :exec
UPDATE games
SET active = ?,
    turns = ?,
    wrong = ?,
    win = ?,
    end = ?
WHERE id = ?;

-- name: UpdateGameTurns :exec
UPDATE games
SET turns = ?,
    wrong = ?
WHERE id = ?;

-- name: UpdateGameStatus :exec
UPDATE games
SET win = ?,
    end = ?,
    active = ?
WHERE id = ?;

-- name: GetUserByKey :one
SELECT uid, key_id, keys
FROM users
WHERE key_id = ?
LIMIT 1;

-- name: GetUserById :one
SELECT uid, key_id, GROUP_CONCAT(keys, '|:|') AS key_string
FROM users
WHERE uid = ?
GROUP BY uid, key_id;

-- name: SaveUser :exec
INSERT INTO users(uid, key_id, keys)
VALUES(?, ?, ?);

-- name: RemoveUserById :exec
DELETE FROM users
WHERE uid = ?;

-- name: RemoveKey :exec
DELETE FROM users
WHERE uid = ? AND key_id = ?;

-- name: UpdateUserKey :exec
UPDATE users
SET keys = ?
WHERE uid = ? AND key_id = ?;
