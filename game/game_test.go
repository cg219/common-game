package game

import (
	"fmt"
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

    fmt.Printf("Game: %v", game)
}

func TestGameWordSelection(t *testing.T) {
    tests := []struct {
            words [4]string
            value bool
    } {
        { [4]string{"Monday", "Tuesday", "Thursday", "Sunday"}, true },
        { [4]string{"Monday", "Friday", "Thursday", "Sunday"}, false },
    }

    game, err := Create()

    if err != nil {
        t.Fatalf("Error occured during game creation: %v", err)
    }

    for _, s := range tests {
        r := game.CheckSelection(s.words)

        if s.value != r {
            t.Fatalf("Words: %s - expected %t; got %t", s.words, s.value, r)
        }
    }
}
