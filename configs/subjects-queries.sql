-- name: SaveSubject :exec
INSERT INTO subjects(subject, words)
VALUES (?, ?);

-- name: GetSubjects :many
SELECT (subject, words)
FROM subjects
