<script lang="ts">
    import Layout from "../lib/Layout.svelte";
    import type { GameResponse } from "../lib/customtypes";

    let words = $state([])
    let gameStatus = $state(2)
    let gid = $state(0)
    let move = $state({
        words: []
    })
    let hasMove = $state(false)
    let moves = $state([])
    let turns = $state(0)

    async function newGame(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/api/game", {
            method: "POST"
        })

        const data = (await res.json()) as GameResponse

        console.log("yerr")
        console.log(data)

        gid = data.id
        words = data.words
        turns = data.moveLeft
        gameStatus = data.status

        console.log(data)
    }

    async function submitMoves() {
        console.log(moves)
        const res = await fetch("/api/game", {
            method: "PUT",
            body: JSON.stringify({
                words: moves
            })
        })

        const data = await res.json() as GameResponse
        reset(moves)
        
        gid = data.id
        words = data.words
        turns = data.moveLeft
        gameStatus = data.status
        hasMove = data.hasMove
        move = data.move

        console.log(data)
    }

    function toggleSelection(m: (string|undefined)[], word: string): number {
        var i = 0
        var next = 0
        var len = 0
        var shouldAdd = true
        var nextSet = false

        m.forEach((w) => {
            if (w == word) {
                m[i] = undefined
                shouldAdd = false
            } else if (w != undefined) {
                if (!nextSet) {
                    next++
                }

                len++
            } else {
                if (!nextSet) {
                    nextSet = true
                    next = i
                }
            }

            i++
        })

        if (shouldAdd && len < 4) {
            m[next] = word
            len++
        }

        return len
    }

    function populateSubject(correct: boolean, subject: string, word: string): string | null {
        if (hasMove) {
            if ( correct && [0, 5].includes(gameStatus) && move.words.includes(word)) {
                return subject
            } else if (correct) {
                return subject
            }
            
            return null
        }

        return subject
    }

    function reset(m: (string|undefined)[]) {
        m.forEach((_, i) => m[i] = undefined)
    }

    function removeWrong() {
        // htmx.findAll(".game-board button.wrong").forEach((b) => b.classList.remove("wrong"))
    }

    function exists(word: string) : boolean {
        console.log(word)
        return move.words.includes(word)
    }

    function selectPiece(val: string) {
        return function(evt: MouseEvent) {
            const piece = evt.target as Element

            if (piece && !piece.classList.contains("correct")) {
                removeWrong()
                const selected = toggleSelection(moves, val)
                piece.classList.toggle("selected")
                console.log(`Moves: ${moves}`)

                if (selected == 4) submitMoves()
            }
        }
    }
</script>

<Layout title="The Common Game" subtitle="Start">
    <form onsubmit={newGame} class="container" id="newgame" method="POST" action="/api/game">
        <button type="submit">Start a New Game</button>
    </form>

    <div class="game">
        {#if gameStatus == 0}
        <h3>WINNER!!</h3>
        {:else if gameStatus == 1}
        <h3>You Lost. Try Again</h3>
        {:else}
        <h3>Mistakes Left: {turns}</h3>
        {/if}

        <section
            class="game-board"
            data-status={gameStatus}
            data-game-id={gid}>
            {#each words as { word, value, correct, subject}}
                <button
                    class="game-piece"
                    class:correct={correct}
                    class:wrong={!correct && exists(value)}
                    data-value={value}
                    data-correct={correct}
                    data-subject={populateSubject(correct, subject, word)}
                    onclick={selectPiece(value)}>{word}</button>
            {/each}
        </section>
    </div>
</Layout>
