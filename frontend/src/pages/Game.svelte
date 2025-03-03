<script lang="ts">
    import type { GameResponse } from "../lib/customtypes";
    import Layout from "../lib/Layout.svelte";
    import { onMount, setContext } from "svelte";
    import GamePiece from "../lib/GamePiece.svelte";
    import { buttonStyle, formStyle } from "../lib/styles";

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

    const legendStyle = "relative text-zinc-100 text-base sm:text-lg leading-[2rem] pl-[3rem] before:block before:absolute before:content-[''] before:w-[2rem] before:h-[2rem] before:left-[0] before:top-[0]";
    const headerStyle = "text-xl text-zinc-100 mb-5";
    const localformStyle = "flex flex-col w-full sm:w-[600px] mx-auto text-slate-100";
    const localbuttonStyle = `${buttonStyle} hover:bg-teal-700 transition-colors duration-200 mb-4`;
</script>

<Layout title="The Common Game" subtitle="Match groups of 4 words that have something in common." links={[]}>
    {#if (gameStatus == 0 || gameStatus == 1)}
    <form onsubmit={getGame} class={localformStyle} id="newgame" method="POST" action="/api/game">
        <button class={localbuttonStyle} type="submit">Start a New Game</button>
    </form>
    {/if}

    <div class="w-full sm:w-[600px] aspect-square mx-auto">
        {#if gameStatus == 0}
        <h3 class={headerStyle}>WINNER!!</h3>
        {:else if gameStatus == 1}
        <h3 class={headerStyle}>You Lost. Try Again</h3>
        {:else}
        <h3 class={headerStyle}>Mistakes Left: {turns}</h3>
        {/if}

        <section class="w-full aspect-square grid grid-cols-[repeat(4,_25%)] grid-rows-4">
            {#each words as word}
                <GamePiece value={word.word} correct={word.correct} wrong={word.wrong} subject={word.subject} selected={word.selected} />
            {/each}
        </section>

        <aside class="py-4">
            <ul class="p-0 flex flex-col gap-4">
                {#each subjects as subject}
                    <li class="{legendStyle}
                        {subject[0] == 0 ? 'before:bg-blue-300' : ''}
                        {subject[0] == 1 ? 'before:bg-pink-300' : ''}
                        {subject[0] == 2 ? 'before:bg-amber-300' : ''}
                        {subject[0] == 3 ? 'before:bg-emerald-300' : ''}"
                    >{subject[1]}</li>
                {/each}
            </ul>
        </aside>
    </div>
</Layout>
