package game

import "fmt"

type LoopStatus int

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

func (s Status) String() string {
    switch s.Metadata.Type.Enum() {
    case Playing.Enum():
        switch s.Metadata.Correct {
        case true:
            return fmt.Sprintf("%s is correct!", s.Metadata.Move.Words)
        case false:
            return fmt.Sprintf("Aww! %s was incorrect. Try Again", s.Metadata.Move.Words)
        }
    case Win.Enum():
        return "WINNER!!"
    case Lose.Enum():
        return "GAME OVER!"
    default:
        return "Something went wrong."
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

func (s LoopStatus) String() string {
    return []string{"Win", "Lose", "Playing", "Broken", "Inactive", "Correct", "Incorrect", "None"}[s]
}

func (s LoopStatus) Enum() int {
    return int(s)
}
