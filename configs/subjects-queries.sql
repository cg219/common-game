-- name: SaveSubject :exec
INSERT INTO subjects(subject, words)
VALUES (?, ?);

-- name: GetSubjects :many
SELECT (subject, words)
FROM subjects;

-- name: GetSubjectsForGame :many
WITH randomrows AS (
    SELECT id, subject, words
    FROM subjects
    ORDER BY random()
    LIMIT 4
),
single_words AS (
    SELECT randomrows.id, json_array(json_each.value) AS stringv
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
