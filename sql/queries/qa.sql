-- name: ReportBug :exec
INSERT INTO bugreports(problem, result, steps, uid, created_at)
VALUES(?, ?, ?, ?, strftime("%s", "now"));
