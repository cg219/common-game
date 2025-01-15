-- name: SaveSubject :exec
INSERT INTO subjects(name, word1, word2, word3, word4)
VALUES (?, ?, ?, ?, ?);

-- name: GetSubjects :many
SELECT name, word1, word2, word3, word4 FROM subjects;

-- name: GetSubjectsForBoard :many
WITH first AS (
    SELECT id, name, word1, word2, word3, word4
    FROM subjects
    ORDER BY random()
    LIMIT 1
),
second AS (
    SELECT id, name, word1, word2, word3, word4
    FROM subjects
    WHERE NOT EXISTS (
        SELECT 1
        FROM first
        WHERE first.word1 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word2 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word3 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word4 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
    )
    ORDER BY random()
    LIMIT 1
),
third AS (
    SELECT id, name, word1, word2, word3, word4
    FROM subjects
    WHERE NOT EXISTS (
        SELECT 1
        FROM first
        WHERE first.word1 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word2 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word3 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word4 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        UNION
        SELECT 1
        FROM second
        WHERE second.word1 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR second.word2 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR second.word3 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR second.word4 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
    )
    ORDER BY random()
    LIMIT 1
),
fourth AS (
    SELECT id, name, word1, word2, word3, word4
    FROM subjects
    WHERE NOT EXISTS (
        SELECT 1
        FROM first
        WHERE first.word1 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word2 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word3 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR first.word4 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        UNION
        SELECT 1
        FROM second
        WHERE second.word1 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR second.word2 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR second.word3 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR second.word4 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        UNION
        SELECT 1
        FROM third
        WHERE third.word1 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR third.word2 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR third.word3 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
        OR third.word4 IN (subjects.word1, subjects.word2, subjects.word3, subjects.word4)
    )
    ORDER BY random()
    LIMIT 1
)
SELECT * FROM first
UNION ALL
SELECT * FROM second
UNION ALL
SELECT * FROM third
UNION ALL
SELECT * FROM fourth;

-- name: GetBoardForGame :one
WITH gids AS (
    SELECT gid
    FROM users_games
    WHERE uid = ?
)
SELECT id, subject1, subject2, subject3, subject4
FROM boards
WHERE id NOT IN (
    SELECT bid
    FROM games
    WHERE id IN gids
)
LIMIT 1;

-- name: SaveNewGame :one
INSERT INTO games(active, start)
VALUES (?, ?)
RETURNING id;

-- name: SaveBoardToGame :exec
UPDATE games
SET bid = ?
WHERE id = ?;

-- name: SaveUserToGame :exec
INSERT INTO users_games(uid, gid)
VALUES(?, ?);

-- name: GetGameUidByGameId :one
SELECT uid
FROM users_games
WHERE gid = ?
LIMIT 1;

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

