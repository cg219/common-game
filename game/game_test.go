package game

import (
	"testing"
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

            if s.turns != game.WrongTurns {
                t.Fatalf("Turn: expected %d; got %d", s.turns, game.WrongTurns)
            }
        }

    }) 

    t.Run("Reset", func(t *testing.T) {
        game.Reset()

        if game.WrongTurns != 0 {
            t.Fatalf("Turns: expected %d; got %d", 0, game.WrongTurns)
        }

        if len(game.CompletedSubjects) > 0 {
            t.Fatalf("Completed Subjects: expected %d; got %d", 0, len(game.CompletedSubjects))
        }
    })

    t.Run("Check LoopStatus", func (t *testing.T) {

    })

    t.Run("Simulate Games", func(t *testing.T) {
        tests := []struct {
            moves []Move
            outcomes []LoopStatus
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
                final: Lose,
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
                final: Lose,
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
                final: Lose,
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
                if s != g.outcomes[step] {
                    t.Fatalf("\nMove: %s\nexpected %s; got %s", g.moves[step].words, g.outcomes[step], s)
                }

                step++
            }
        }
    })
}
