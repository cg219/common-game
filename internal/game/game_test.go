package game

import (
	"testing"
)

func TestGame(t *testing.T) {
    boardData := []GameData{
        GameData{ Name: "Days of the Week", Word: "Monday", Word2: "Tuesday", Word3: "Thursday", Word4: "Sunday" },
        GameData{ Name: "Through the Air", Word: "Leap", Word2: "Soar", Word3: "Float", Word4: "Fly" },
        GameData{ Name: "Races", Word: "Black", Word2: "White", Word3: "Hispanic", Word4: "Indian" },
        GameData{ Name: "Colors", Word: "Brown", Word2: "Red", Word3: "Blue", Word4: "Orange" },
    }

    game := Create(boardData)

    if game == nil {
        t.Fatal("expected Game, got nil")
    }

    newMove := func (words ...string) Move {
        var a [4]string

        copy(a[:], words)

        move := &Move{
            Words: a,
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

    // t.Run("Test Game Cancellation", func(t *testing.T) {
    //     tests := []struct {
    //         moves []Move
    //         outcome LoopStatus
    //     } {
    //         {
    //             moves: []Move{
    //                 newMove("Monday", "Tuesday", "Thursday", "Sunday"),
    //                 newMove("Leap", "Soar", "Float", "Fly"),
    //             },
    //             outcome: Inactive,
    //         },
    //     }
    //
    //     for _, g := range tests {
    //         game.Reset()
    //         game.HealthTickInvteral = 100 * time.Millisecond
    //         output, input := game.Run()
    //
    //         go func() {
    //             for _, m := range g.moves {
    //                 input <- m
    //                 time.Sleep(200 * time.Millisecond)
    //             }
    //
    //             close(input)
    //         }()
    //
    //         var status LoopStatus
    //
    //         for s := range output {
    //             status = s.LoopStatus
    //         }
    //
    //         if status.Enum() != g.outcome.Enum() {
    //             i := len(g.moves) - 1
    //             t.Fatalf("\nMove: %s\nexpected %s; got %s", g.moves[i].Words, g.outcome, status)
    //         }
    //     }
    // })

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

        game.Reset()
        statuses, input := game.Run()
        step := 0

        go func() {
            for _, m := range test.moves {
                input <- m
            }

            close(input)
        }()

        for s := range statuses {
            if s.Status.Status() != test.outcomes[step] {
                t.Fatalf("\nMove: %s\nexpected %s; got %s", test.moves[step].Words, test.outcomes[step], s.Status.Status())
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
            stats struct {
                correct int
                wrong int
                total int
            }
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
                stats: struct{correct int; wrong int; total int}{
                    correct: 4,
                    wrong: 0,
                    total: 4,
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
                stats: struct{correct int; wrong int; total int}{
                    correct: 3,
                    wrong: 4,
                    total: 7,
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
                stats: struct{correct int; wrong int; total int}{
                    correct: 4,
                    wrong: 3,
                    total: 7,
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
                stats: struct{correct int; wrong int; total int}{
                    correct: 4,
                    wrong: 2,
                    total: 6,
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
                stats: struct{correct int; wrong int; total int}{
                    correct: 4,
                    wrong: 2,
                    total: 6,
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
                    Inactive,
                    Inactive,
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
                stats: struct{correct int; wrong int; total int}{
                    correct: 2,
                    wrong: 4,
                    total: 6,
                },
                final: Broken,
            },
        }

        for i, g := range tests {
            game.Reset()
            output, input := game.Run()

            go func() {
                for _, m := range g.moves {
                    input <- m
                }

                close(input)
            }()

            step := 0

            for s := range output {
                if s.LoopStatus != g.outcomes[step] {
                    t.Fatalf("\nLoop Status Error\nMove: %s\nexpected %s; got %s", g.moves[step].Words, g.outcomes[step], s.LoopStatus)
                }

                if s.Status.Status() != g.statuses[step] {
                    t.Fatalf("\nStatus Error\nMove: %s\nexpected %s; got %s", g.moves[step].Words, g.statuses[step], s.Status)
                }

                step++
            }

            if g.stats.correct != game.Metadata.Correct {
                t.Fatalf("Stats Correct Error: expected %d; got %d; Test: %d", g.stats.correct, game.Metadata.Correct, i)
            }

            if g.stats.wrong != game.Metadata.WrongTurns {
                t.Fatalf("Stats Wrong Error: expected %d; got %d; Test: %d", g.stats.wrong, game.Metadata.WrongTurns, i)
            }

            if g.stats.total != game.Metadata.TotalTurns {
                t.Fatalf("Stats Total Error: expected %d; got %d; Test: %d", g.stats.total, game.Metadata.TotalTurns, i)
            }
        }
    })
}
