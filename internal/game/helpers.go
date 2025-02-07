package game

import (
	"log"
	"time"

	"github.com/cg219/common-game/internal/database"
)

type wordsResponse struct {
    Words []string `json:"words"`
}

func Create(rows []database.PopulateSubjectsRow) *Game {
    return newGame(rows)
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
            output <- s.Status
        }

        close(output)
    }()

    return output
}

func newGame(rows []database.PopulateSubjectsRow) *Game {
    var subjects [4]Subject

    if rows != nil {

        for i, v := range rows {
            var words [4]string

            words[0] = v.Word.String
            words[1] = v.Word_2.String
            words[2] = v.Word_3.String
            words[3] = v.Word_4.String

            subjects[i] = Subject{ Name: v.Name, Words: words }
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
        // select {
        // case move, ok := <-input:
            move, ok := <-input
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
        // case <-tick.C:
        //     g.IsInactive = true
        //     status := Status{
        //         Metadata: StatusMetadata{
        //             Type: Inactive,
        //         },
        //     }
        //
        //     output <- StatusGroup{
        //         LoopStatus: Inactive,
        //         Status: status,
        //     } 
        //
        //     return
        // }
    }
}
