package jeeves

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
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

func GenerateBoards(amount int) {
    q, db, ctx, stop := setup()
    defer db.Close()
    defer stop()

    fmt.Printf("generating %d boards\n", amount)

    for range amount {
        subjects, err := q.GetSubjectsForBoard(ctx)
        if err != nil {
            log.Fatalf("err getting subjects: %s\n", err.Error())
        }

        ids := make([]int, 0)

        for _, i := range subjects {
            ids = append(ids, int(i))
        }

        sort.Ints(ids)

        err = q.SaveNewBoard(ctx, database.SaveNewBoardParams{
            Subject1: sql.NullInt64{ Int64: int64(ids[0]), Valid: true },
            Subject2: sql.NullInt64{ Int64: int64(ids[1]), Valid: true },
            Subject3: sql.NullInt64{ Int64: int64(ids[2]), Valid: true },
            Subject4: sql.NullInt64{ Int64: int64(ids[3]), Valid: true },
        })

        if err != nil {
            log.Fatalf("err saving board: %s\n", err.Error())
        }
    }
}

func GetNewBoard(uid int) {
    q, db, ctx, stop := setup()
    defer db.Close()
    defer stop()

    fmt.Printf("get new board for uid: %d\n", uid)

    board, err := q.GetBoardForGame(ctx, int64(uid));
    if err != nil {
        log.Fatalf("err getting board: %s\n", err.Error())
    }

    res, err := q.PopulateSubjects(ctx, database.PopulateSubjectsParams{
        ID: board.Subject1.Int64,
        ID_2: board.Subject2.Int64,
        ID_3: board.Subject3.Int64,
        ID_4: board.Subject4.Int64,
    })

    if err != nil {
        log.Fatalf("err populating subjects: %s\n", err.Error())
    }

    fmt.Println(board)
    fmt.Println(res)
}
