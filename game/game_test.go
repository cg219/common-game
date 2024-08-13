package game

import (
	"testing"
	"time"
)

func TestGame(t *testing.T) {
    game, err := Create()

    if err != nil {
        t.Fatalf("error: %s", err)
    }

    if game == nil {
        t.Fatal("expected Game, got nil")
    }

    newMove := func (words ...string) Move {
        var a [4]string

        copy(a[:], words)

        move := &Move{
            words: a,
        }

        return *move
    }

    t.Run("Words Selection", func(t *testing.T) {
        wordSelection := []struct {
            words [4]string
            value bool
            turns int
            catValue string
        } {
            { [4]string{"Monday", "Tuesday", "Thursday", "Sunday"}, true, 0, "Days of the Week"},
            { [4]string{"Monday", "Friday", "Thursday", "Sunday"}, false, 1, ""},
        }

        for _, s := range wordSelection {
            r, cat := game.CheckSelection(s.words)

            if s.value != r {
                t.Fatalf("\nWords: %s\nexpected %t; got %t", s.words, s.value, r)
            }

            if cat != nil {
                if s.catValue != cat.Name {
                    t.Fatalf("Category: expected %s; got %s", s.catValue, cat.Name)
                }
            }

            if s.turns != game.Metadata.WrongTurns {
                t.Fatalf("Turn: expected %d; got %d", s.turns, game.Metadata.WrongTurns)
            }
        }

    }) 

    t.Run("Reset", func(t *testing.T) {
        game.Reset()

        if game.Metadata.WrongTurns != 0 {
            t.Fatalf("Turns: expected %d; got %d", 0, game.Metadata.WrongTurns)
        }

        if len(game.CompletedSubjects) > 0 {
            t.Fatalf("Completed Subjects: expected %d; got %d", 0, len(game.CompletedSubjects))
        }
    })

    t.Run("Check LoopStatus", func (t *testing.T) {

    })

    t.Run("Test Game Cancellation", func(t *testing.T) {
        tests := []struct {
            moves []Move
            outcome LoopStatus
        } {
            {
                moves: []Move{
                    newMove("Monday", "Tuesday", "Thursday", "Sunday"),
                    newMove("Leap", "Soar", "Float", "Fly"),
                },
                outcome: Inactive,
            },
        }

        for _, g := range tests {
            game.Reset()
            game.HealthTickInvteral = 100 * time.Millisecond
            input := make(chan Move)
            output := game.Run(input)

            go func() {
                for _, m := range g.moves {
                    input <- m
                    time.Sleep(200 * time.Millisecond)
                }

                close(input)
            }()

            var status LoopStatus

            for s := range output {
                status = s.LoopStatus
            }

            if status.Enum() != g.outcome.Enum() {
                i := len(g.moves) - 1
                t.Fatalf("\nMove: %s\nexpected %s; got %s", g.moves[i].words, g.outcome, status)
            }
        }
    })

    t.Run("Test StartGame", func(t *testing.T) {
        test := struct {
            moves []Move
            outcomes []LoopStatus
        } {
            moves: []Move{
                newMove("Monday", "Tuesday", "Thursday", "Sunday"),
                newMove("Leap", "Soar", "Float", "Fly"),
                newMove("Black", "Soar", "Hispanic", "Indian"),
                newMove("Black", "White", "Blue", "Indian"),
                newMove("Brown", "Red", "Blue", "Orange"),
                newMove("Black", "Red", "Blue", "Orange"),
                newMove("Black", "Hispanic", "Blue", "Orange"),
            },
            outcomes: []LoopStatus{
                Correct,
                Correct,
                Incorrect,
                Incorrect,
                Correct,
                Incorrect,
                Lose,
            },
        }

        input := make(chan Move)
        statuses := Start(input)
        step := 0

        go func() {
            for _, m := range test.moves {
                input <- m
            }

            close(input)
        }()

        for s := range statuses {
            if s.Status() != test.outcomes[step] {
                t.Fatalf("\nMove: %s\nexpected %s; got %s", test.moves[step].words, test.outcomes[step], s.Status())
            }
            step++
        }
    })

    t.Run("Simulate Games", func(t *testing.T) {
        tests := []struct {
            moves []Move
            outcomes []LoopStatus
            statuses []LoopStatus
            final LoopStatus
        } {
            {
                moves: []Move{
                    newMove("Monday", "Tuesday", "Thursday", "Sunday"),
                    newMove("Leap", "Soar", "Float", "Fly"),
                    newMove("Black", "White", "Hispanic", "Indian"),
                    newMove("Brown", "Red", "Blue", "Orange"),
                },
                outcomes: []LoopStatus{
                    Playing,
                    Playing,
                    Playing,
                    Win,
                },
                statuses: []LoopStatus{
                    Correct,
                    Correct,
                    Correct,
                    Win,
                },
                final: Win,
            },
            {
                moves: []Move{
                    newMove("Monday", "Tuesday", "Thursday", "Sunday"),
                    newMove("Leap", "Soar", "Float", "Fly"),
                    newMove("Black", "Soar", "Hispanic", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Brown", "Red", "Blue", "Orange"),
                    newMove("Black", "Red", "Blue", "Orange"),
                    newMove("Black", "Hispanic", "Blue", "Orange"),
                },
                outcomes: []LoopStatus{
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Lose,
                },
                statuses: []LoopStatus{
                    Correct,
                    Correct,
                    Incorrect,
                    Incorrect,
                    Correct,
                    Incorrect,
                    Lose,
                },
                final: Lose,
            },
            {
                moves: []Move{
                    newMove("Monday", "Tuesday", "Thursday", "Sunday"),
                    newMove("Leap", "Soar", "Float", "Fly"),
                    newMove("Black", "Soar", "Hispanic", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Brown", "Red", "Blue", "Orange"),
                    newMove("Black", "Red", "Blue", "Orange"),
                    newMove("Black", "White", "Hispanic", "Indian"),
                },
                outcomes: []LoopStatus{
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Win,
                },
                statuses: []LoopStatus{
                    Correct,
                    Correct,
                    Incorrect,
                    Incorrect,
                    Correct,
                    Incorrect,
                    Win,
                },
                final: Win,
            },
            {
                moves: []Move{
                    newMove("Tuesday", "Sunday", "Thursday", "Monday"),
                    newMove("Leap", "Soar", "Float", "Fly"),
                    newMove("Black", "Soar", "Hispanic", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Brown", "Red", "Blue", "Orange"),
                    newMove("Black", "Hispanic", "White", "Indian"),
                },
                outcomes: []LoopStatus{
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Win,
                },
                statuses: []LoopStatus{
                    Correct,
                    Correct,
                    Incorrect,
                    Incorrect,
                    Correct,
                    Win,
                },
                final: Win,
            },
            {
                moves: []Move{
                    newMove("Tuesday", "Sunday", "thursdaY", "Monday"),
                    newMove("Leap", "Soar", "float", "Fly"),
                    newMove("Black", "Soar", "hispanic", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Brown", "red", "Blue", "Orange"),
                    newMove("black", "Hispanic", "White", "Indian"),
                },
                outcomes: []LoopStatus{
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Win,
                },
                statuses: []LoopStatus{
                    Correct,
                    Correct,
                    Incorrect,
                    Incorrect,
                    Correct,
                    Win,
                },
                final: Win,
            },
            {
                moves: []Move{
                    newMove("Tuesday", "Sunday", "thursdaY", "Monday"),
                    newMove("Leap", "Soar", "float", "Fly"),
                    newMove("Black", "Soar", "hispanic", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Black", "White", "Blue", "Indian"),
                    newMove("Brown", "red", "Blue", "Orange"),
                    newMove("Brown", "red", "Blue", "Orange"),
                },
                outcomes: []LoopStatus{
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Playing,
                    Lose,
                    Broken,
                    Broken,
                },
                statuses: []LoopStatus{
                    Correct,
                    Correct,
                    Incorrect,
                    Incorrect,
                    Incorrect,
                    Lose,
                    None,
                    None,
                },
                final: Broken,
            },
        }

        for _, g := range tests {
            game.Reset()
            input := make(chan Move)
            output := game.Run(input)

            go func() {
                for _, m := range g.moves {
                    input <- m
                }

                close(input)
            }()

            step := 0

            for s := range output {
                if s.LoopStatus != g.outcomes[step] {
                    t.Fatalf("\nLoop Status Error\nMove: %s\nexpected %s; got %s", g.moves[step].words, g.outcomes[step], s.LoopStatus)
                }

                if s.Status.Status() != g.statuses[step] {
                    t.Fatalf("\nStatus Error\nMove: %s\nexpected %s; got %s", g.moves[step].words, g.statuses[step], s.Status)
                }

                step++
            }
        }
    })
}
