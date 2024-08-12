package game

import (
	"fmt"
	"strings"
	"time"
)

type LoopStatus int

type Game struct {
    WrongTurns int
    MaxTurns int
    Subjects [4]Category
    CompletedSubjects []int
    HealthTickInvteral time.Duration
    IsInactive bool
}

type Category struct {
    Name string
    Words [4]string
}

type Move struct {
    words [4]string
}

const (
    Win LoopStatus = iota
    Lose
    Playing
    Broken
    Inactive
    None
)

func Create() (*Game, error) {
    return newGame(), nil
}

func Start(g *Game) {
    var game *Game

    if g == nil {
        var err error

        g, err = Create()

        if err != nil {
            panic(err)
        }
    }

    game = g
    game.Reset()
    input := make(chan Move)
    output := game.Run(input)

    for s := range output {
        fmt.Println(s)
    }
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
        HealthTickInvteral: 2 * time.Minute,
    }
}

func loop(input <-chan Move, output chan<- LoopStatus, g *Game) {
    tick := time.NewTicker(g.HealthTickInvteral)

    defer tick.Stop()
    defer close(output)

    for {
        select {
        case move, ok := <-input:
            if !ok {
                return
            }

            g.CheckSelection(move.words)
            s := g.CheckStatus()

            tick.Reset(g.HealthTickInvteral)
            output <- s
        case <-tick.C:
            g.IsInactive = true
            output <- Inactive 
            return
        }
    }
}

// Start LoopStatus Interface

func (s LoopStatus) String() string {
    return []string{"Win", "Lose", "Playing", "Broken", "Inactive", "None"}[s]
}

func (s LoopStatus) Enum() int {
    return int(s)
}

// --End LoopStatus Interface

// Start Game Interface

func (g *Game) CheckSelection(words [4]string) (bool, *Category) {
    cat := -1
    matches := 0

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
    g.IsInactive = false
}

func (g *Game) Run(ch <-chan Move) <-chan LoopStatus {
    statusCh := make(chan LoopStatus)

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

func (g *Game) CheckStatus() LoopStatus {
    if g.IsInactive {
        return Inactive
    }

    if g.WrongTurns > g.MaxTurns {
        return Broken
    }

    lookup := make(map[int]bool)

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
