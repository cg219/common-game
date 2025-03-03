<script lang="ts">
    import Layout from "../lib/Layout.svelte";
    import { formStyle, textareaStyle, buttonStyle, linkStyle, setStyle } from "../lib/styles";

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
        <a class={linkStyle} href="/game">Go back to playing</a>
    {:else}
        <form onsubmit={report} class={`${formStyle} text-zinc-100`} id="report" method="POST" action="/api/report">
            <fieldset class={setStyle}>
                <label for="problem">What went wrong?</label>
                <textarea class={textareaStyle} name="problem" placeholder="Problem" bind:value={problem}></textarea>
            </fieldset>

            <fieldset class={setStyle}>
                <label for="result">What did you expect to happen?</label>
                <textarea class={textareaStyle} name="result" placeholder="Expectation" bind:value={result}></textarea>
            </fieldset>

            <fieldset class={setStyle}>
                <label for="steps">How can we reproduce the problem?</label>
                <textarea class={textareaStyle} name="steps" placeholder="Steps to reproduce" bind:value={steps}></textarea>
            </fieldset>

            <fieldset class={setStyle}>
                <button class={buttonStyle} type="submit">Submit Report</button>
            </fieldset>
        </form>
    {/if}
</Layout>


