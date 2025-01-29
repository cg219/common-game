package importer

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/cg219/common-game/internal/database"
	_ "modernc.org/sqlite"
)

func setup() (*database.Queries, *sql.DB, context.Context, context.CancelFunc) {
    cwd, _ := os.Getwd();
    db, err := sql.Open("sqlite", filepath.Join(cwd, "data/database.db"))
    if err != nil {
        log.Fatalf("err opening db: %s", err.Error())
    }

    queries := database.New(db)
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

    return queries, db, ctx, stop
}

func transformValue(value string) string {
    replacer := strings.NewReplacer(
        "'", "",
        "-", "",
        "_", "",
    )
    return strings.ToLower(
        replacer.Replace(
            strings.TrimSpace(value),
        ),
    )
}

func importSubject(ctx context.Context, row []string, query *database.Queries, conn *sql.DB) {
    fmt.Printf("importing subject: %s\n", row[0])
    tx, err := conn.BeginTx(ctx, nil)
    if err != nil {
        log.Fatalf("err creating tx: %s", err.Error())
    }

    qtx := query.WithTx(tx)
    wids := make([]int, 4)
    wvals := make([]string, 4)

    for i, word := range row {
        if i == 0 {
            continue
        } 

        val := transformValue(word)
        wvals[i-1] = val
        res, err := query.GetWordByValue(ctx, val)
        if err != nil {
            if err == sql.ErrNoRows {
                fmt.Printf("adding word: %s\n", word)
                err = qtx.SaveWord(ctx, database.SaveWordParams{ Value: val, Word: strings.TrimSpace(word) })
                if err != nil {
                    tx.Rollback()
                    log.Fatalf("err saving word value: %s", err.Error())
                }
            } else {
                log.Fatalf("err getting word value: %s", err.Error())
            }
        } else {
            wids[i-1] = int(res.ID)
        }
    }

    err = tx.Commit()
    if err != nil {
        tx.Rollback()
        log.Fatalf("err committing transaction: %s", err.Error())
    }

    for i := range row {
        if i == 0 {
            continue
        }

        if wids[i-1] == 0 {
            res, err := query.GetWordByValue(ctx, wvals[i-1])
            if err != nil {
                if err != sql.ErrNoRows {
                    log.Fatalf("No Value-tored, Fatal Error: %s", err.Error())
                } else {
                    log.Fatalf("err getting word by val: %s", err.Error())
                }
            }

            wids[i-1] = int(res.ID)
        }
    }

    sort.Ints(wids)

    err = query.SaveSubject(ctx, database.SaveSubjectParams{
        Name: row[0],
        Word1: int64(wids[0]),
        Word2: int64(wids[1]),
        Word3: int64(wids[2]),
        Word4: int64(wids[3]),
    })

    if err != nil {
        fmt.Printf("subject \"%s\" exists; skipping\n", row[0])
    } else {
        fmt.Printf("done: %s\n", row[0])
    }
}

func Run(r io.Reader) {
    q, db, ctx, stop := setup()
    defer db.Close()
    defer stop()

    csvreader := csv.NewReader(r)

    for {
        select {
        case <- ctx.Done():
            return
        default:
            row, err := csvreader.Read()
            if err ==io.EOF {
                return
            }

            if err != nil {
                log.Fatalf("err reading csv: %s", err.Error())
            }
            
            importSubject(ctx, row, q, db)
        }
    }
}
