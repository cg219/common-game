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
}
