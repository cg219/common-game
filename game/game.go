package game

import "strings"

type Game struct {
    Turns int
    Subjects [4]Category
}

type Category struct {
    Name string
    Words [4]string
}

func Create() (*Game, error) {
    return newGame(), nil
}

func newGame() *Game {
    var subjects [4]Category

    subjects[0] = Category{
        Name: "Days of the Week",
        Words: [4]string{"Monday", "Tuesday", "Thursday", "Sunday"},
    }

    subjects[1] = Category{
        Name: "Through the Air",
        Words: [4]string{"Leap", "Soar", "Float", "Fly"},
    }

    subjects[2] = Category{
        Name: "Races",
        Words: [4]string{"Black", "White", "Hispanic", "Indian"},
    }

    subjects[3] = Category{
        Name: "Colors",
        Words: [4]string{"Brown", "Red", "Blue", "Orange"},
    }

    return &Game{
        Turns: 1,
        Subjects: subjects,
    }
}

// Start Game Interface

func (g *Game) CheckSelection(words [4]string) bool {
    var cat string
    matches := 0

    for _, cw := range words { // cw = Current Word
        for _, cc := range g.Subjects { // cc = Current Category
            for _, ccw := range cc.Words { // ccw = Current Category Word
                if strings.EqualFold(cw, ccw) {
                    if cat == "" {
                        cat = cc.Name
                        matches++
                        continue
                    }

                    if cat != cc.Name {
                        return false
                    }

                    matches++
                }
            }
        }
    }

    return matches == 4
}

// --End Game Interface
