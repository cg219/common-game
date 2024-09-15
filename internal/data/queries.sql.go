// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package data

import (
	"context"
	"database/sql"
)

const getPlayerById = `-- name: GetPlayerById :one
SELECT id, name
FROM players
WHERE id = ?
`

func (q *Queries) GetPlayerById(ctx context.Context, id int64) (Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerById, id)
	var i Player
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getPlayerByName = `-- name: GetPlayerByName :one
SELECT id, name
FROM players
WHERE name = ?
`

func (q *Queries) GetPlayerByName(ctx context.Context, name sql.NullString) (Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerByName, name)
	var i Player
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getSubjects = `-- name: GetSubjects :many
SELECT subject, words FROM subjects
`

type GetSubjectsRow struct {
	Subject string
	Words   string
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
		if err := rows.Scan(&i.Subject, &i.Words); err != nil {
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

const getSubjectsForGame = `-- name: GetSubjectsForGame :many
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
)
`

type GetSubjectsForGameRow struct {
	ID      int64
	Subject string
	Words   string
}

func (q *Queries) GetSubjectsForGame(ctx context.Context) ([]GetSubjectsForGameRow, error) {
	rows, err := q.db.QueryContext(ctx, getSubjectsForGame)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubjectsForGameRow
	for rows.Next() {
		var i GetSubjectsForGameRow
		if err := rows.Scan(&i.ID, &i.Subject, &i.Words); err != nil {
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

const getWords = `-- name: GetWords :many
SELECT word
FROM words
`

func (q *Queries) GetWords(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getWords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		items = append(items, word)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveNewGame = `-- name: SaveNewGame :one
INSERT INTO games(active, start, player_id)
VALUES (?, ?, ?)
RETURNING id
`

type SaveNewGameParams struct {
	Active   sql.NullBool
	Start    sql.NullInt64
	PlayerID sql.NullInt64
}

func (q *Queries) SaveNewGame(ctx context.Context, arg SaveNewGameParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, saveNewGame, arg.Active, arg.Start, arg.PlayerID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const savePlayer = `-- name: SavePlayer :exec
INSERT INTO players (name)
VALUES (?)
`

func (q *Queries) SavePlayer(ctx context.Context, name sql.NullString) error {
	_, err := q.db.ExecContext(ctx, savePlayer, name)
	return err
}

const saveSubject = `-- name: SaveSubject :exec
INSERT INTO subjects(subject, words)
VALUES (?, ?)
`

type SaveSubjectParams struct {
	Subject string
	Words   string
}

func (q *Queries) SaveSubject(ctx context.Context, arg SaveSubjectParams) error {
	_, err := q.db.ExecContext(ctx, saveSubject, arg.Subject, arg.Words)
	return err
}

const saveWord = `-- name: SaveWord :exec
INSERT INTO words (word)
VALUES(?)
`

func (q *Queries) SaveWord(ctx context.Context, word string) error {
	_, err := q.db.ExecContext(ctx, saveWord, word)
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
