package game

import (
	"testing"
)

func TestCreateGame(t *testing.T) {
    game, err := Create()

    if err != nil {
        t.Fatalf("Error occured during game creation: %v", err)
    }

    if game == nil {
        t.Fatalf("No Game Created. Expected Game.")
    }
}

func TestGameWordSelection(t *testing.T) {
    tests := []struct {
            words [4]string
            value bool
            catValue string
    } {
        { [4]string{"Monday", "Tuesday", "Thursday", "Sunday"}, true, "Days of the Week"},
        { [4]string{"Monday", "Friday", "Thursday", "Sunday"}, false , ""},
    }

    game, err := Create()

    if err != nil {
        t.Fatalf("Error occured during game creation: %v", err)
    }

    for _, s := range tests {
        r, cat := game.CheckSelection(s.words)

        if s.value != r {
            t.Fatalf("Words: %s - expected %t; got %t", s.words, s.value, r)
        }

        if cat != nil {
            if s.catValue != cat.Name {
                t.Fatalf("Slected Category Incorrect. expected %s; got %s", s.catValue, cat.Name)
            }
        }
    }
}
