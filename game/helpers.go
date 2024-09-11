package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cg219/common-game/internal/data"
)

type GameConfig struct {
    Q *data.Queries
    Ctx context.Context
}

type wordsResponse struct {
    Words []string `json:"words"`
}

func Create(ng GameConfig) (*Game, error) {
    return newGame(ng), nil
}

func StartWithGame(g *Game) <-chan Status {
    return start(g)
}

func Start() <-chan Status {
    return start(nil)
}

func start(g *Game) <-chan Status {
    var game *Game

    if g == nil {
        log.Fatalf("Oops: Missing Game")
    }

    game = g
    game.Reset()
    output := make(chan Status)
    status, _ := game.Run()

    go func() {
        for s := range status {
            // fmt.Println(s.LoopStatus, s.Status)
            output <- s.Status
        }

        close(output)
    }()

    return output
}

func newGame(ng GameConfig) *Game {
    var subjects [4]Subject

    if ng.Q != nil {
        res, err := ng.Q.GetSubjectsForGame(ng.Ctx)

        if err != nil {
            log.Fatalf("Oops: %s", err)
        }

        fmt.Println(res)

        for i, v := range res {
            var words [4]string
            var tmp []interface{}

            if err := json.Unmarshal([]byte(v.Words), &tmp); err != nil {
                log.Fatalf("Oops: Missing Game")
            }

            for j, lv := range tmp { 
                words[j] = lv.(string)
            }

            subjects[i] = Subject{ Name: v.Subject, Words: words }
        }

        return &Game{
            Metadata: struct{Player; Stats}{
                Player: Player{},
                Stats: Stats{},
            },
            MaxTurns: 4,
            Subjects: subjects,
            HealthTickInvteral: 2 * time.Minute,
        }
    }

    subjects[0] = Subject{
        Name: "Days of the Week",
        Words: [4]string{"Monday", "Tuesday", "Thursday", "Sunday"},
    }

    subjects[1] = Subject{
        Name: "Through the Air",
        Words: [4]string{"Leap", "Soar", "Float", "Fly"},
    }

    subjects[2] = Subject{
        Name: "Races",
        Words: [4]string{"Black", "White", "Hispanic", "Indian"},
    }

    subjects[3] = Subject{
        Name: "Colors",
        Words: [4]string{"Brown", "Red", "Blue", "Orange"},
    }

    return &Game{
        Metadata: struct{Player; Stats}{
            Player: Player{},
            Stats: Stats{},
        },
        MaxTurns: 4,
        Subjects: subjects,
        HealthTickInvteral: 2 * time.Minute,
    }
}

func loop(input <-chan Move, output chan<- StatusGroup, g *Game) {
    tick := time.NewTicker(g.HealthTickInvteral)
    g.Metadata.StartTime = time.Now().UTC()

    defer tick.Stop()
    defer func() {
        g.Metadata.EndTime = time.Now().UTC()
    }()
    defer close(output)

    for {
        select {
        case move, ok := <-input:
            if !ok {
                return
            }

            correct, sub := g.CheckSelection(move.Words)
            loopStatus := g.CheckStatus()

            var status Status

            if sub == nil {
                status = Status{
                    Metadata: StatusMetadata{
                        Move: move,
                        Type: loopStatus,
                        Correct: correct,
                    },
                }
            } else {
                status = Status{
                    Metadata: StatusMetadata{
                        Move: move,
                        Subject: *sub,
                        Type: loopStatus,
                        Correct: correct,
                    },
                }
            }


            tick.Reset(g.HealthTickInvteral)

            output <- StatusGroup{
                LoopStatus: loopStatus,
                Status: status,
            }
        case <-tick.C:
            g.IsInactive = true
            status := Status{
                Metadata: StatusMetadata{
                    Type: Inactive,
                },
            }

            output <- StatusGroup{
                LoopStatus: Inactive,
                Status: status,
            } 

            return
        }
    }
}
