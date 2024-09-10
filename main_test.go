package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/cg219/common-game/internal/data"
    _ "github.com/tursodatabase/go-libsql"
)

func TestGameSubjects(t *testing.T) {
    ctx := context.Background()
    ddl, err := os.ReadFile("./configs/schema.sql")
    if err != nil {
        t.Fatalf("Error: %s", err)
    }

    db, err := sql.Open("libsql", "file:./database.db")
    if err != nil {
        t.Fatalf("Error: %s", err)
    }

    defer db.Close()

    if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
        t.Fatalf("Error: %s", err)
    }

    sq := data.New(db)
    res, err := sq.GetSubjectsForGame(ctx)

    if err != nil {
        t.Fatalf("Error: %s", err)
    }

    fmt.Println(res)
}
