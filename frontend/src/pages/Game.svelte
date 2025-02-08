<script lang="ts">
    import type { GameResponse } from "../lib/customtypes";
    import Layout from "../lib/Layout.svelte";
    import { onMount, setContext } from "svelte";
    import GamePiece from "../lib/GamePiece.svelte";

    let words = $state([])
    let gameStatus = $state(2)
    let gid = $state(0)
    let turns = $state(0)
    let subjects = $state([])

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

    onMount(() => {
        getGame()
    })

    async function getGame(evt?: Event) {
        evt?.preventDefault()

        const res = await fetch("/api/game", {
            method: "POST"
        })

        const data = (await res.json()) as GameResponse

        gid = data.id
        words = data.words
        turns = data.moveLeft
        gameStatus = data.status
        subjects = []
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

        if (data.move.correct) {
            setSubjects(data.move.subjects)
            return setCorrect(data.move.words)
        }

        return setWrong(data.move.words)
    }

    function setSubjects(val: { id: number, name: string }[]) {
        subjects = val.map((v) => {
            return [v.id, v.name]
        })
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

<Layout title="The Common Game" subtitle="Match groups of 4 words that have something in common." links={[]}>
    {#if (gameStatus == 0 || gameStatus == 1)}
    <form onsubmit={getGame} class="container" id="newgame" method="POST" action="/api/game">
        <button type="submit">Start a New Game</button>
    </form>
    {/if}

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

        <aside>
            <ul>
                {#each subjects as subject}
                    <li
                        class:s0={subject[0] == 0}
                        class:s1={subject[0] == 1}
                        class:s2={subject[0] == 2}
                        class:s3={subject[0] == 3}
                    >{subject[1]}</li>
                {/each}
            </ul>
        </aside>
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

    aside {
        ul {
            padding: 0;
        }

        li {
            position: relative;
            padding-left: 3rem;

            &:before {
                display: block;
                position: absolute;
                content: "";
                width: 2rem;
                height: 2rem;
                left: 0;
                top: 0;
            }

            &.s0:before {
                background-color: lightblue;
            }

            &.s1:before {
                background-color: lightpink;
            }

            &.s2:before {
                background-color: lightgoldenrodyellow;
            }

            &.s3:before {
                background-color: lightsteelblue;
            }
        }
    }
</style>
