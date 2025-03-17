package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/cg219/common-game/internal/database"
	"github.com/cg219/common-game/internal/game"
)

type Storage struct {
    q *database.Queries
    log *slog.Logger
}

type ActiveGame struct {
    Id int
    Data *LiveGameData
}

func NewStorage(q *database.Queries, log *slog.Logger) *Storage {
    return &Storage{ q: q }
}

func (s *Storage) GetNewGame(conn *sql.DB, uid int) (error, int64, *game.Game) {
    return s.GetNewGameWithContext(context.Background(), conn, uid)
}

func (s *Storage) GetActiveGames() (error, []ActiveGame) {
    return s.GetActiveGamesWithContext(context.Background())
}

func (s *Storage) GetNewGameWithContext(ctx context.Context, conn *sql.DB, uid int) (error, int64, *game.Game) {
    avoid, err := s.q.GetRecentlyPlayedSubjects(ctx, int64(uid))
    if err != nil {
        s.log.Error("Error retreiving recent subjects", "err",err)
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    avoidList := make([]int64, 0)

    for _, v := range avoid {
        avoidList = append(avoidList, v.Subject1.Int64)
        avoidList = append(avoidList, v.Subject2.Int64)
        avoidList = append(avoidList, v.Subject3.Int64)
        avoidList = append(avoidList, v.Subject4.Int64)
    }

    pad := 16 - len(avoidList)

    for range pad {
        avoidList = append(avoidList, int64(0))
    }

    board, err := s.q.GetBoardForGame(ctx, database.GetBoardForGameParams{
        Uid: int64(uid),
        ID: avoidList[0],
        ID_2: avoidList[1],
        ID_3: avoidList[2],
        ID_4: avoidList[3],
        ID_5: avoidList[4],
        ID_6: avoidList[5],
        ID_7: avoidList[6],
        ID_8: avoidList[7],
        ID_9: avoidList[8],
        ID_10: avoidList[9],
        ID_11: avoidList[10],
        ID_12: avoidList[11],
        ID_13: avoidList[12],
        ID_14: avoidList[13],
        ID_15: avoidList[14],
        ID_16: avoidList[15],
    })
    if err != nil {
        s.log.Error("Error retreiving subjects", "err", err)
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    populatedBoard, err := s.q.PopulateSubjects(ctx, database.PopulateSubjectsParams{
        ID: board.Subject1.Int64,
        ID_2: board.Subject2.Int64,
        ID_3: board.Subject3.Int64,
        ID_4: board.Subject4.Int64,
    })

    tx, err := conn.BeginTx(ctx, nil)
    if err != nil {
        s.log.Error("Error creating tx", "err", err)
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    qtx := s.q.WithTx(tx)

    id, err := qtx.SaveNewGame(ctx, database.SaveNewGameParams{
        Active: sql.NullBool{ Bool: true, Valid: true },
        Start: sql.NullInt64{ Int64: time.Now().UTC().UnixMilli(), Valid: true },
    })

    if err != nil {
        s.log.Error("Error creating game", "err", err)
        tx.Rollback()
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    err = qtx.SaveUserToGame(ctx, database.SaveUserToGameParams{ Uid: int64(uid), Gid: id })
    if err != nil {
        s.log.Error("Error saving user to game", "err", err)
        tx.Rollback()
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    err = qtx.SaveBoardToGame(ctx, database.SaveBoardToGameParams{
        Bid: sql.NullInt64{ Int64: board.ID, Valid: true },
        ID: id,
    })
    if err != nil {
        s.log.Error("Error saving board to game", "err", err)
        tx.Rollback()
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    err = tx.Commit()
    if err != nil {
        s.log.Error("Error committing tx", "tx", tx, "err", err)
        tx.Rollback()
        return fmt.Errorf(INTERNAL_ERROR), 0, nil
    }

    return nil, id, game.Create(populatedBoard)
}

func (s *Storage) GetActiveGamesWithContext(ctx context.Context) (error, []ActiveGame) {
    games, err := s.q.GetActiveGames(ctx)

    if err != nil {
        s.log.Error("loading active games", "err", err)
        return err, []ActiveGame{}
    }

    list := []ActiveGame{}

    for _, cg := range games {
        populatedBoard, err := s.q.PopulateSubjects(context.Background(), database.PopulateSubjectsParams{
            ID: cg.Subject1.Int64,
            ID_2: cg.Subject2.Int64,
            ID_3: cg.Subject3.Int64,
            ID_4: cg.Subject4.Int64,
        })

        if err != nil {
            s.log.Error("populating active games", "err", err)
            return err, []ActiveGame{}
        }

        // TODO: Update game logic to maintain state or be able to export an encoding to load from
        g := game.Create(populatedBoard)

        statusCh, moveCh := g.Run()

        list = append(list, ActiveGame{
            Id: int(cg.ID),
            Data: &LiveGameData{
                game: g,
                mch: moveCh,
                sch: statusCh,
            },
        })
    }

    return nil, list
}
