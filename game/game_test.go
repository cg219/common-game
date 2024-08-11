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
            catValue string
        } {
            { [4]string{"Monday", "Tuesday", "Thursday", "Sunday"}, true, "Days of the Week"},
            { [4]string{"Monday", "Friday", "Thursday", "Sunday"}, false, ""},
        }

        for i, s := range wordSelection {
            r, cat := game.CheckSelection(s.words)

            if s.value != r {
                t.Fatalf("\nWords: %s\nexpected %t; got %t", s.words, s.value, r)
            }

            if cat != nil {
                if s.catValue != cat.Name {
                    t.Fatalf("Category: expected %s; got %s", s.catValue, cat.Name)
                }
            }

            if i + 1 != game.TurnsTaken {
                t.Fatalf("Turn: expected %d; got %d", i + 1, game.TurnsTaken)
            }
        }

    }) 

    t.Run("Reset", func(t *testing.T) {
        game.Reset()

        if game.TurnsTaken != 0 {
            t.Fatalf("Turns: expected %d; got %d", 0, game.TurnsTaken)
        }

        if len(game.CompletedSubjects) > 0 {
            t.Fatalf("Completed Subjects: expected %d; got %d", 0, len(game.CompletedSubjects))
        }
    })

    t.Run("Check Status", func (t *testing.T) {

    })

    t.Run("Simulate Games", func(t *testing.T) {
        tests := []struct {
            moves []Move
            outcomes []Status
            final Status
        } {
            {
                moves: []Move{
                    newMove("Monday", "Tuesday", "Thurdday", "Sunday"),
                    newMove("Leap", "Soar", "Float", "Fly"),
                    newMove("Black", "White", "Hispanic", "Indian"),
                    newMove("Brown", "Red", "Blue", "Orange"),
                },
                outcomes: []Status{
                    Playing,
                    Playing,
                    Playing,
                    Win,
                },
                final: Win,
            },
        }

        for _, g := range tests {
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
