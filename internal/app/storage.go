package app

import (
	"context"
	"log/slog"

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

func (s *Storage) GetActiveGames() (error, []ActiveGame) {
    return s.GetActiveGamesWithContext(context.Background())
}
