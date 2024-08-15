package game

import (
	"math/rand"
	"strings"
	"time"
)

type Subject struct {
    Name string
    Words [4]string
}

type Move struct {
    Words [4]string
}

type Game struct {
    MaxTurns int
    Subjects [4]Subject
    CompletedSubjects []int
    HealthTickInvteral time.Duration
    IsInactive bool
    Metadata struct {
        Player
        Stats
    }
}

type Player struct {
    Username string
    Wins int
    Losses int
}

type Stats struct {
    TotalTurns int
    WrongTurns int
    Correct int
    StartTime time.Time
    EndTime time.Time
}

func (g *Game) Words() []string {
    words := make([]string, 16)
    seed := rand.Perm(16)

    var merged []string

    for _, s := range g.Subjects {
        for _, w := range s.Words {
            merged = append(merged, w)
        }
    }

    for i, w := range seed {
        words[w] = merged[i]
    }

    return words
}

func (g *Game) CheckSelection(words [4]string) (bool, *Subject) {
    cat := -1
    matches := 0

    if g.IsInactive {
        return false, nil
    }

    g.Metadata.TotalTurns++

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
                        g.Metadata.WrongTurns++
                        return false, nil
                    }

                    matches++
                }
            }
        }
    }

    if matches == 4 && !g.IsInactive {
        g.CompletedSubjects = append(g.CompletedSubjects, cat)
        g.Metadata.Correct++
        return true, &g.Subjects[cat]
    }

    g.Metadata.WrongTurns++
    return false, nil
}

func (g *Game) Reset() {
    g.Metadata.Stats = Stats{}
    g.CompletedSubjects = make([]int, 0)
    g.IsInactive = false
}

func (g *Game) Run() (<-chan StatusGroup, chan<- Move) {
    statusCh := make(chan StatusGroup)
    moveCh := make(chan Move)

    go loop(moveCh, statusCh, g)

    return statusCh, moveCh
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

    if g.Metadata.WrongTurns > g.MaxTurns {
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
            if g.Metadata.WrongTurns == g.MaxTurns {
                g.IsInactive = true
                return Lose
            }

            return Playing
        }
    } 

    if len(lookup) == 4 {
        g.IsInactive = true
        return Win
    }

    if g.Metadata.WrongTurns == g.MaxTurns {
        g.IsInactive = true
        return Lose
    }

    return Playing
}
