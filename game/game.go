package game

import (
	"strings"
)

type Status int

type Game struct {
    WrongTurns int
    MaxTurns int
    Subjects [4]Category
    CompletedSubjects []int
}

type Category struct {
    Name string
    Words [4]string
}

type Move struct {
    words [4]string
}

const (
    Win Status = iota
    Lose
    Playing
    Broken
    None
)

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
        WrongTurns: 0,
        MaxTurns: 4,
        Subjects: subjects,
    }
}

func loop(input <-chan Move, output chan<- Status, g *Game) {
    for move := range input {
        g.CheckSelection(move.words)
        s := g.CheckStatus()

        output <- s
    }

    close(output)
}

// Start Status Interface

func (s Status) String() string {
    return []string{"Win", "Lose", "Playing", "Broken", "None"}[s]
}

func (s Status) Enum() int {
    return int(s)
}

// --End Status Interface

// Start Game Interface

func (g *Game) CheckSelection(words [4]string) (bool, *Category) {
    cat := -1
    matches := 0

    defer func() {
    }()

    for _, cw := range words { // cw = Current Word
        for csi, cc := range g.Subjects { // cc = Current Category csi = Current Subject Index
            for _, ccw := range cc.Words { // ccw = Current Category Word
                if strings.EqualFold(cw, ccw) {
                    if cat == -1 {
                        cat = csi
                        matches++
                        continue
                    }

                    if cat != csi {
                        g.WrongTurns++
                        return false, nil
                    }

                    matches++
                }
            }
        }
    }

    if matches == 4 {
        g.CompletedSubjects = append(g.CompletedSubjects, cat)
        return true, &g.Subjects[cat]
    }

    g.WrongTurns++
    return false, nil
}

func (g *Game) Reset() {
    g.WrongTurns = 0
    g.CompletedSubjects = make([]int, 0)
}

func (g *Game) Run(ch <-chan Move) <-chan Status {
    statusCh := make(chan Status)

    go loop(ch, statusCh, g)

    return statusCh
}

func (g *Game) CheckSubjectStatus(s int) bool {
    for _, cs := range g.CompletedSubjects {
        if cs == s {
            return true
        }
    }
    return false
}

func (g *Game) CheckStatus() Status {
    lookup := make(map[int]bool)

    if g.WrongTurns > g.MaxTurns {
        return Broken
    }

    for _, cs := range g.CompletedSubjects {
        v, ok := lookup[cs]

        if !ok {
            lookup[cs] = true
            continue
        }

        if v {
            if g.WrongTurns == g.MaxTurns {
                g.WrongTurns++
                return Lose
            }

            return Playing
        }
    } 

    if len(lookup) == 4 {
        return Win
    }

    if g.WrongTurns == g.MaxTurns {
        g.WrongTurns++
        return Lose
    }

    return Playing
}

// --End Game Interface
