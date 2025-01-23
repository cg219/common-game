export type Link = {
    url: string
    name: string
    current: boolean
}

export type GameResponse = {
    words: {
        correct: boolean
        word: string
    }[]
    id: number
    moveLeft: number
    status: number
    hasMove: boolean
    move: {
        correct: boolean
        words: string[]
        subjects: {
            id: number
            name: string
        }[]
    }
}
