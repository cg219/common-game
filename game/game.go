package game

import "strings"

type Status int

const (
    subject1 = iota
    sujbect2
    subject3
    subject4
)

const (
    Win Status = iota
    Lose
    Playing
    Broken
    None
)

type Game struct {
    TurnsTaken int
    MaxTurns int
    Subjects [4]Category
    CompletedSubjects []int
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

    subjects[subject1] = Category{
        Name: "Days of the Week",
        Words: [4]string{"Monday", "Tuesday", "Thursday", "Sunday"},
    }

    subjects[sujbect2] = Category{
        Name: "Through the Air",
        Words: [4]string{"Leap", "Soar", "Float", "Fly"},
    }

    subjects[subject3] = Category{
        Name: "Races",
        Words: [4]string{"Black", "White", "Hispanic", "Indian"},
    }

    subjects[subject4] = Category{
        Name: "Colors",
        Words: [4]string{"Brown", "Red", "Blue", "Orange"},
    }

    return &Game{
        TurnsTaken: 0,
        MaxTurns: 4,
        Subjects: subjects,
    }
}

// Start Game Interface

func (g *Game) CheckSelection(words [4]string) (bool, *Category) {
    cat := -1
    matches := 0

    defer func() {
        g.TurnsTaken++
    }()

    for _, cw := range words { // cw = Current Word
        for csi, cc := range g.Subjects { // cc = Current Category csi = Current Subject Index
            for _, ccw := range cc.Words { // ccw = Current Category Word
                if strings.EqualFold(cw, ccw) {
                    if cat == -1 {
                        cat = csi
                        g.CompletedSubjects = append(g.CompletedSubjects, csi)
                        matches++
                        continue
                    }

                    if cat != csi {
                        return false, nil
                    }

                    g.CompletedSubjects = append(g.CompletedSubjects, csi)
                    matches++
                }
            }
        }
    }

    if matches == 4 {
        return true, &g.Subjects[cat]
    }

    return false, nil
}

func (g *Game) CheckStatus() Status {
    lookup := make(map[int]bool)

    if g.TurnsTaken < g.MaxTurns {
        return Playing
    }

    if g.TurnsTaken == g.MaxTurns {
        for _, cs := range g.CompletedSubjects {
            v, ok := lookup[cs]

            if !ok {
                lookup[cs] = true
                continue
            }

            if v {
                return Lose
            }
        }    

        return Win
    }

    if g.TurnsTaken > g.MaxTurns {
        return Broken
    }

    return None
}

// --End Game Interface
