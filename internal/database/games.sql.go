// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: games.sql

package database

import (
	"context"
	"database/sql"
)

const getActiveGidForUid = `-- name: GetActiveGidForUid :one
WITH lgid AS (
    SELECT ug.gid
    FROM users_games ug
    WHERE ug.uid = ?
    ORDER BY ug.gid DESC
    LIMIT 1
),
agid AS (
    SELECT g.id
    FROM games g
    WHERE g.id = (SELECT l.gid FROM lgid l) AND g.active = 1
    LIMIT 1
)
SELECT CAST(COALESCE((SELECT id FROM agid), 0) AS INTEGER)
`

func (q *Queries) GetActiveGidForUid(ctx context.Context, uid int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getActiveGidForUid, uid)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getBoardForGame = `-- name: GetBoardForGame :one
WITH gids AS (
    SELECT gid
    FROM users_games ug
    WHERE ug.uid = ?
),
sids AS (
    SELECT sub.id
    FROM subjects sub
    WHERE sub.id IN (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    LIMIT 16
),
bids AS (
    SELECT bid
    FROM games g
    WHERE g.id IN gids
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
AND b.id NOT IN bids
ORDER BY random()
LIMIT 1
`

type GetBoardForGameParams struct {
	Uid   int64
	ID    int64
	ID_2  int64
	ID_3  int64
	ID_4  int64
	ID_5  int64
	ID_6  int64
	ID_7  int64
	ID_8  int64
	ID_9  int64
	ID_10 int64
	ID_11 int64
	ID_12 int64
	ID_13 int64
	ID_14 int64
	ID_15 int64
	ID_16 int64
}

type GetBoardForGameRow struct {
	ID       int64
	Subject1 sql.NullInt64
	Subject2 sql.NullInt64
	Subject3 sql.NullInt64
	Subject4 sql.NullInt64
}

func (q *Queries) GetBoardForGame(ctx context.Context, arg GetBoardForGameParams) (GetBoardForGameRow, error) {
	row := q.db.QueryRowContext(ctx, getBoardForGame,
		arg.Uid,
		arg.ID,
		arg.ID_2,
		arg.ID_3,
		arg.ID_4,
		arg.ID_5,
		arg.ID_6,
		arg.ID_7,
		arg.ID_8,
		arg.ID_9,
		arg.ID_10,
		arg.ID_11,
		arg.ID_12,
		arg.ID_13,
		arg.ID_14,
		arg.ID_15,
		arg.ID_16,
	)
	var i GetBoardForGameRow
	err := row.Scan(
		&i.ID,
		&i.Subject1,
		&i.Subject2,
		&i.Subject3,
		&i.Subject4,
	)
	return i, err
}

const getGameUidByGameId = `-- name: GetGameUidByGameId :one
SELECT uid
FROM users_games
WHERE gid = ?
LIMIT 1
`

func (q *Queries) GetGameUidByGameId(ctx context.Context, gid int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getGameUidByGameId, gid)
	var uid int64
	err := row.Scan(&uid)
	return uid, err
}

const getRecentlyPlayedSubjects = `-- name: GetRecentlyPlayedSubjects :many
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
LIMIT 5
`

type GetRecentlyPlayedSubjectsRow struct {
	Subject1 sql.NullInt64
	Subject2 sql.NullInt64
	Subject3 sql.NullInt64
	Subject4 sql.NullInt64
}

func (q *Queries) GetRecentlyPlayedSubjects(ctx context.Context, uid int64) ([]GetRecentlyPlayedSubjectsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRecentlyPlayedSubjects, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRecentlyPlayedSubjectsRow
	for rows.Next() {
		var i GetRecentlyPlayedSubjectsRow
		if err := rows.Scan(
			&i.Subject1,
			&i.Subject2,
			&i.Subject3,
			&i.Subject4,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubjectName = `-- name: GetSubjectName :one
SELECT name
FROM subjects
WHERE id = ?
`

func (q *Queries) GetSubjectName(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRowContext(ctx, getSubjectName, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const getSubjects = `-- name: GetSubjects :many
SELECT name, word1, word2, word3, word4 FROM subjects
`

type GetSubjectsRow struct {
	Name  string
	Word1 int64
	Word2 int64
	Word3 int64
	Word4 int64
}

func (q *Queries) GetSubjects(ctx context.Context) ([]GetSubjectsRow, error) {
	rows, err := q.db.QueryContext(ctx, getSubjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubjectsRow
	for rows.Next() {
		var i GetSubjectsRow
		if err := rows.Scan(
			&i.Name,
			&i.Word1,
			&i.Word2,
			&i.Word3,
			&i.Word4,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubjectsForBoard = `-- name: GetSubjectsForBoard :many
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
SELECT id FROM fourth
`

func (q *Queries) GetSubjectsForBoard(ctx context.Context) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getSubjectsForBoard)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const populateSubjects = `-- name: PopulateSubjects :many
SELECT s.name, w.word, w2.word, w3.word, w4.word
FROM subjects s
LEFT JOIN words w ON s.word1 = w.id
LEFT JOIN words w2 ON s.word2 = w2.id
LEFT JOIN words w3 ON s.word3 = w3.id
LEFT JOIN words w4 ON s.word4 = w4.id
WHERE s.id IN (?, ?, ?, ?)
`

type PopulateSubjectsParams struct {
	ID   int64
	ID_2 int64
	ID_3 int64
	ID_4 int64
}

type PopulateSubjectsRow struct {
	Name   string
	Word   sql.NullString
	Word_2 sql.NullString
	Word_3 sql.NullString
	Word_4 sql.NullString
}

func (q *Queries) PopulateSubjects(ctx context.Context, arg PopulateSubjectsParams) ([]PopulateSubjectsRow, error) {
	rows, err := q.db.QueryContext(ctx, populateSubjects,
		arg.ID,
		arg.ID_2,
		arg.ID_3,
		arg.ID_4,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PopulateSubjectsRow
	for rows.Next() {
		var i PopulateSubjectsRow
		if err := rows.Scan(
			&i.Name,
			&i.Word,
			&i.Word_2,
			&i.Word_3,
			&i.Word_4,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveBoardToGame = `-- name: SaveBoardToGame :exec
UPDATE games
SET bid = ?
WHERE id = ?
`

type SaveBoardToGameParams struct {
	Bid sql.NullInt64
	ID  int64
}

func (q *Queries) SaveBoardToGame(ctx context.Context, arg SaveBoardToGameParams) error {
	_, err := q.db.ExecContext(ctx, saveBoardToGame, arg.Bid, arg.ID)
	return err
}

const saveNewBoard = `-- name: SaveNewBoard :exec
INSERT INTO boards(subject1, subject2, subject3, subject4)
VALUES(?, ?, ?, ?)
`

type SaveNewBoardParams struct {
	Subject1 sql.NullInt64
	Subject2 sql.NullInt64
	Subject3 sql.NullInt64
	Subject4 sql.NullInt64
}

func (q *Queries) SaveNewBoard(ctx context.Context, arg SaveNewBoardParams) error {
	_, err := q.db.ExecContext(ctx, saveNewBoard,
		arg.Subject1,
		arg.Subject2,
		arg.Subject3,
		arg.Subject4,
	)
	return err
}

const saveNewGame = `-- name: SaveNewGame :one
INSERT INTO games(active, start)
VALUES (?, ?)
RETURNING id
`

type SaveNewGameParams struct {
	Active sql.NullBool
	Start  sql.NullInt64
}

func (q *Queries) SaveNewGame(ctx context.Context, arg SaveNewGameParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, saveNewGame, arg.Active, arg.Start)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const saveSubject = `-- name: SaveSubject :exec
INSERT INTO subjects(name, word1, word2, word3, word4)
VALUES (?, ?, ?, ?, ?)
`

type SaveSubjectParams struct {
	Name  string
	Word1 int64
	Word2 int64
	Word3 int64
	Word4 int64
}

func (q *Queries) SaveSubject(ctx context.Context, arg SaveSubjectParams) error {
	_, err := q.db.ExecContext(ctx, saveSubject,
		arg.Name,
		arg.Word1,
		arg.Word2,
		arg.Word3,
		arg.Word4,
	)
	return err
}

const saveUserToGame = `-- name: SaveUserToGame :exec
INSERT INTO users_games(uid, gid)
VALUES(?, ?)
`

type SaveUserToGameParams struct {
	Uid int64
	Gid int64
}

func (q *Queries) SaveUserToGame(ctx context.Context, arg SaveUserToGameParams) error {
	_, err := q.db.ExecContext(ctx, saveUserToGame, arg.Uid, arg.Gid)
	return err
}

const updateBoard = `-- name: UpdateBoard :exec
UPDATE boards
SET played = ?,
    wins = ?
WHERE id = ?
`

type UpdateBoardParams struct {
	Played sql.NullInt64
	Wins   sql.NullInt64
	ID     int64
}

func (q *Queries) UpdateBoard(ctx context.Context, arg UpdateBoardParams) error {
	_, err := q.db.ExecContext(ctx, updateBoard, arg.Played, arg.Wins, arg.ID)
	return err
}

const updateGame = `-- name: UpdateGame :exec
UPDATE games
SET active = ?,
    turns = ?,
    wrong = ?,
    win = ?,
    end = ?
WHERE id = ?
`

type UpdateGameParams struct {
	Active sql.NullBool
	Turns  sql.NullInt64
	Wrong  sql.NullInt64
	Win    sql.NullBool
	End    sql.NullInt64
	ID     int64
}

func (q *Queries) UpdateGame(ctx context.Context, arg UpdateGameParams) error {
	_, err := q.db.ExecContext(ctx, updateGame,
		arg.Active,
		arg.Turns,
		arg.Wrong,
		arg.Win,
		arg.End,
		arg.ID,
	)
	return err
}

const updateGameStatus = `-- name: UpdateGameStatus :exec
UPDATE games
SET win = ?,
    end = ?,
    active = ?
WHERE id = ?
`

type UpdateGameStatusParams struct {
	Win    sql.NullBool
	End    sql.NullInt64
	Active sql.NullBool
	ID     int64
}

func (q *Queries) UpdateGameStatus(ctx context.Context, arg UpdateGameStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateGameStatus,
		arg.Win,
		arg.End,
		arg.Active,
		arg.ID,
	)
	return err
}

const updateGameTurns = `-- name: UpdateGameTurns :exec
UPDATE games
SET turns = ?,
    wrong = ?
WHERE id = ?
`

type UpdateGameTurnsParams struct {
	Turns sql.NullInt64
	Wrong sql.NullInt64
	ID    int64
}

func (q *Queries) UpdateGameTurns(ctx context.Context, arg UpdateGameTurnsParams) error {
	_, err := q.db.ExecContext(ctx, updateGameTurns, arg.Turns, arg.Wrong, arg.ID)
	return err
}
