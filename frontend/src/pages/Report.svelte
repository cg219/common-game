<script lang="ts">
    import Layout from "../lib/Layout.svelte";

    let problem = $state("")
    let result = $state("")
    let steps = $state("")
    let showThanks = $state(false)

    async function report(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/api/report", {
            method: "POST",
            body: JSON.stringify({ problem, result, steps })
        })

        const data = await res.json()

        if (data.success) showThanks = true
    }
</script>

<Layout title="The Common Game" subtitle="Report a Bug">
    {#if showThanks}
        <h1>Thanks!!!</h1>
        <p>Added your report. Stay tuned for any updates.</p>
        <a href="/game">Go back to playing</a>
    {:else}
        <form onsubmit={report} class="container" id="report" method="POST" action="/api/report">
            <label for="problem">What went wrong?</label>
            <textarea name="problem" placeholder="Problem" bind:value={problem}></textarea>

            <label for="result">What did you expect to happen?</label>
            <textarea name="result" placeholder="Expectation" bind:value={result}></textarea>

            <label for="steps">How can we reproduce the problem?</label>
            <textarea name="steps" placeholder="Steps to reproduce" bind:value={steps}></textarea>

            <button type="submit">Submit Report</button>
        </form>
    {/if}
</Layout>


