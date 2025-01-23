<script lang="ts">
    import type { GameResponse } from "../lib/customtypes";
    import Layout from "../lib/Layout.svelte";
    import { setContext } from "svelte";
    import GamePiece from "../lib/GamePiece.svelte";

    let words = $state([])
    let gameStatus = $state(2)
    let gid = $state(0)
    let hasMove = $state(false)
    let turns = $state(0)

    setContext("select", (value: string, wasSelected: string) => {
        if (![2, 5, 6].includes(gameStatus)) return;
        const idx = words.findIndex((w) => w.word == value)

        words[idx].selected = wasSelected ? false : true;

        const selected = words
            .map((w) => {
                w.wrong = false;
                return w;
            })
            .filter((w) => w.selected)
            .map((w) => w.word)

        if (selected.length == 4) {
            submitMoves(selected)
        }
    })

    async function newGame(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/api/game", {
            method: "POST"
        })

        const data = (await res.json()) as GameResponse

        gid = data.id
        words = data.words
        turns = data.moveLeft
        gameStatus = data.status
    }

    async function submitMoves(moves: string[]) {
        const res = await fetch("/api/game", {
            method: "PUT",
            body: JSON.stringify({
                words: moves,
                gid
            })
        })

        const data = await res.json() as GameResponse

        gid = data.id
        words = data.words
        turns = data.moveLeft
        gameStatus = data.status
        hasMove = data.hasMove

        if (data.move.correct) return setCorrect(data.move.words)
        return setWrong(data.move.words)
    }

    function setCorrect(val: string[]) {
       val.forEach((w) => {
            const i = words.findIndex((x) => x.word == w)
            words[i].correct = true;
        }) 
    }

    function setWrong(val: string[]) {
       val.forEach((w) => {
            const i = words.findIndex((x) => x.word == w)
            words[i].wrong = true;
        }) 
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

        <section>
            {#each words as word}
                <GamePiece value={word.word} correct={word.correct} wrong={word.wrong} subject={word.subject} selected={word.selected} />
            {/each}
        </section>
    </div>
</Layout>

<style>
    section {
        width: 100%;
        height: 100%;
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        grid-template-rows: repeat(4, 1fr);
    }
</style>
