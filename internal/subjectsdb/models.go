// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package storage

import (
	"database/sql"
)

type Subject struct {
	ID      int64
	Subject string
	Words   string
	Used    sql.NullInt64
	Correct sql.NullInt64
}
