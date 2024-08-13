package game

import (
	"strings"
	"time"
)

type Subject struct {
    Name string
    Words [4]string
}

type Move struct {
    words [4]string
}

type Game struct {
    WrongTurns int
    MaxTurns int
    Subjects [4]Subject
    CompletedSubjects []int
    HealthTickInvteral time.Duration
    IsInactive bool
}

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
