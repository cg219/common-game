-- name: GetUserSession :one
SELECT accessToken, refreshToken, valid
FROM sessions
WHERE accessToken = ? AND refreshToken = ?
LIMIT 1;

-- name: SaveUserSession :exec
INSERT INTO sessions(accessToken, refreshToken)
VALUES(?, ?);

-- name: InvalidateUserSession :exec
UPDATE sessions
SET valid = 0
WHERE accessToken = ? AND refreshToken = ?;

-- name: SaveUser :exec
INSERT INTO users(username, email, password, valid_token)
VALUES(?, ?, ?, ?);

-- name: GetUser :one
SELECT id, username, email, valid
FROM users
WHERE username = ?;

-- name: GetUserByValidToken :one
SELECT id, username
FROM users
WHERE valid_token = ?
LIMIT 1;

-- name: InvalidateUser :exec
UPDATE users
SET valid = NULL
WHERE username = ?;

-- name: ValidateUser :exec
UPDATE users
SET valid = strftime("%s", "now"),
    valid_token = NULL
WHERE username = ?;

-- name: GetUserWithPassword :one
SELECT username, password, valid
FROM users
WHERE username = ?;

-- name: SetPasswordReset :exec
UPDATE users
SET reset = ?,
    reset_time = ?
WHERE username = ?;

-- name: ResetPassword :exec
UPDATE users
SET reset = NULL,
    reset_time = NULL,
    password = ?
WHERE reset = ? AND reset_time > ?;

-- name: CanResetPassword :one
SELECT reset_time > ? AS valid, username
FROM users
WHERE reset = ?;

-- name: SaveApiKey :exec
INSERT INTO apikeys(name, key, uid)
VALUES(?, ?, ?);

-- name: GetApiKeysForUid :many
SELECT key, name
FROM apikeys
WHERE uid = ?;

-- name: CheckValidApiKey :one
SELECT EXISTS (
    SELECT 1
    FROM apikeys
    WHERE key = ?
) as valid;

-- name: GetUserFromApiKey :one
SELECT username
FROM users
JOIN apikeys
ON users.id = apikeys.uid
WHERE apikeys.key = ?;
