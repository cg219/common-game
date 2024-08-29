-- name: SaveSubject :exec
INSERT INTO subjects(subject, words)
VALUES (?, ?);

-- name: GetSubjects :many
SELECT (wubject, words)
FROM subjects
