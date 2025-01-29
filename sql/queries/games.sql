-- name: SaveSubject :exec
INSERT INTO subjects(name, word1, word2, word3, word4)
VALUES (?, ?, ?, ?, ?);

-- name: GetSubjectName :one
SELECT name
FROM subjects
WHERE id = ?;

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
SELECT id FROM first
UNION ALL
SELECT id FROM second
UNION ALL
SELECT id FROM third
UNION ALL
SELECT id FROM fourth;

-- name: GetRecentlyPlayedSubjects :many
WITH gids AS (
    SELECT gid
    FROM users_games ug
    WHERE ug.uid = ?
),
bids AS (
    SELECT bid
    FROM games g
    WHERE g.id IN gids
)
SELECT b.subject1, b.subject2, b.subject3, b.subject4
FROM boards b
WHERE b.id IN bids
ORDER BY b.id DESC
LIMIT 5;

-- name: GetBoardForGame :one
WITH sids AS (
    SELECT s.id
    FROM subjects s
    WHERE s.id IN (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    LIMIT 16
)
SELECT b.id, b.subject1, b.subject2, b.subject3, b.subject4
FROM boards b
WHERE NOT EXISTS (
    SELECT 1
    FROM sids s
    WHERE s.id = b.subject1
    OR s.id = b.subject2
    OR s.id = b.subject3
    OR s.id = b.subject4
)
LIMIT 1;

-- name: PopulateSubjects :many
SELECT s.name, w.word, w2.word, w3.word, w4.word
FROM subjects s
LEFT JOIN words w ON s.word1 = w.id
LEFT JOIN words w2 ON s.word2 = w2.id
LEFT JOIN words w3 ON s.word3 = w3.id
LEFT JOIN words w4 ON s.word4 = w4.id
WHERE s.id IN (?, ?, ?, ?);

-- name: SaveNewBoard :exec
INSERT INTO boards(subject1, subject2, subject3, subject4)
VALUES(?, ?, ?, ?);

-- name: UpdateBoard :exec
UPDATE boards
SET played = ?,
    wins = ?
WHERE id = ?;

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

