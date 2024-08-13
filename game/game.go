package game

import (
	"fmt"
	"strings"
	"time"
)

type LoopStatus int
type Subject struct {
    Name string
    Words [4]string
}

type Move struct {
    words [4]string
}

type StatusMetadata struct {
    Move Move
    Subject Subject
    Type LoopStatus
    Correct bool
}

type Status struct {
    Metadata StatusMetadata
}

type StatusGroup struct {
    LoopStatus LoopStatus
    Status Status
}

type Game struct {
    WrongTurns int
    MaxTurns int
    Subjects [4]Subject
    CompletedSubjects []int
    HealthTickInvteral time.Duration
    IsInactive bool
}

const (
    Win LoopStatus = iota
    Lose
    Playing
    Broken
    Inactive
    Correct
    Incorrect
    None
)

func Create() (*Game, error) {
    return newGame(), nil
}

func StartWithGame(input <-chan Move, g *Game) <-chan Status {
    return start(input, g)
}

func Start(input <-chan Move) <-chan Status {
    return start(input, nil)
}

func start(input <-chan Move, g *Game) <-chan Status {
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
    output := make(chan Status)
    status := game.Run(input)

    go func() {
        for s := range status {
            // fmt.Println(s.LoopStatus, s.Status)
            output <- s.Status
        }

        close(output)
    }()

    return output
}

func newGame() *Game {
    var subjects [4]Subject

    subjects[0] = Subject{
        Name: "Days of the Week",
        Words: [4]string{"Monday", "Tuesday", "Thursday", "Sunday"},
    }

    subjects[1] = Subject{
        Name: "Through the Air",
        Words: [4]string{"Leap", "Soar", "Float", "Fly"},
    }

    subjects[2] = Subject{
        Name: "Races",
        Words: [4]string{"Black", "White", "Hispanic", "Indian"},
    }

    subjects[3] = Subject{
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

func loop(input <-chan Move, output chan<- StatusGroup, g *Game) {
    tick := time.NewTicker(g.HealthTickInvteral)

    defer tick.Stop()
    defer close(output)

    for {
        select {
        case move, ok := <-input:
            if !ok {
                return
            }

            correct, sub := g.CheckSelection(move.words)
            loopStatus := g.CheckStatus()

            var status Status

            if sub == nil {
                status = Status{
                    Metadata: StatusMetadata{
                        Move: move,
                        Type: loopStatus,
                        Correct: correct,
                    },
                }
            } else {
                status = Status{
                    Metadata: StatusMetadata{
                        Move: move,
                        Subject: *sub,
                        Type: loopStatus,
                        Correct: correct,
                    },
                }
            }


            tick.Reset(g.HealthTickInvteral)

            output <- StatusGroup{
                LoopStatus: loopStatus,
                Status: status,
            }
        case <-tick.C:
            g.IsInactive = true
            status := Status{
                Metadata: StatusMetadata{
                    Type: Inactive,
                },
            }

            output <- StatusGroup{
                LoopStatus: Inactive,
                Status: status,
            } 

            return
        }
    }
}

func (s Status) String() string {
    switch s.Metadata.Type.Enum() {
    case Playing.Enum():
        switch s.Metadata.Correct {
        case true:
            return fmt.Sprintf("%s is correct!\n", s.Metadata.Move.words)
        case false:
            return fmt.Sprintf("Aww! %s was incorrect. Try Again\n", s.Metadata.Move.words)
        }
    case Win.Enum():
        return fmt.Sprintf("WINNER!!\n")
    case Lose.Enum():
        return fmt.Sprintf("GAME OVER!\n")
    default:
        return fmt.Sprintf("Something went wrong.\n")
    }

    return ""
}

func (s Status) Status() LoopStatus {
    switch s.Metadata.Type.Enum() {
    case Win.Enum():
        return Win
    case Lose.Enum():
        return Lose
    case Playing.Enum():
        if s.Metadata.Correct {
            return Correct
        }
        
        return Incorrect
    default:
        return None
    }
}

// Start LoopStatus Interface

func (s LoopStatus) String() string {
    return []string{"Win", "Lose", "Playing", "Broken", "Inactive", "Correct", "Incorrect", "None"}[s]
}

func (s LoopStatus) Enum() int {
    return int(s)
}

// --End LoopStatus Interface

// Start Game Interface

func (g *Game) CheckSelection(words [4]string) (bool, *Subject) {
    cat := -1
    matches := 0

    for _, cw := range words { // cw = Current Word
        for csi, cc := range g.Subjects { // cc = Current Subject csi = Current Subject Index
            for _, ccw := range cc.Words { // ccw = Current Subject Word
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

func (g *Game) Run(ch <-chan Move) <-chan StatusGroup {
    statusCh := make(chan StatusGroup)

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
