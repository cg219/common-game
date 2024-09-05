package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/cg219/common-game/internal/subjectsdb"
	_ "modernc.org/sqlite"
)

func TestGameSubjects(t *testing.T) {
    ctx := context.Background()
    ddl, err := os.ReadFile("./configs/subjects-schema.sql")
    if err != nil {
        t.Fatalf("Error: %s", err)
    }

    db, err := sql.Open("sqlite", "subjects.db")
    if err != nil {
        t.Fatalf("Error: %s", err)
    }

    defer db.Close()

    if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
        t.Fatalf("Error: %s", err)
    }

    sq := subjectsdb.New(db)
    res, err := sq.GetSubjectsForGame(ctx)

    if err != nil {
        t.Fatalf("Error: %s", err)
    }

    fmt.Println(res)
}
