package game

import (
	"time"
)

func Create() (*Game, error) {
    return newGame(), nil
}

func StartWithGame(input <-chan Move, g *Game) <-chan Status {
    return start(input, g)
}

func Start(input <-chan Move) <-chan Status {
    return start(input, nil)
}

func start(input <-chan Move, g *Game) <-chan Status {
    var game *Game

    if g == nil {
        var err error

        g, err = Create()

        if err != nil {
            panic(err)
        }
    }

    game = g
    game.Reset()
    output := make(chan Status)
    status := game.Run(input)

    go func() {
        for s := range status {
            // fmt.Println(s.LoopStatus, s.Status)
            output <- s.Status
        }

        close(output)
    }()

    return output
}

func newGame() *Game {
    var subjects [4]Subject

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
    g.Metadata.StartTime = time.Now()

    defer tick.Stop()
    defer func() {
        g.Metadata.EndTime = time.Now()
    }()
    defer close(output)

    for {
        select {
        case move, ok := <-input:
            if !ok {
                return
            }

            correct, sub := g.CheckSelection(move.words)
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
