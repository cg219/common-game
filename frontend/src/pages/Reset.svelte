<script lang="ts">
    import { onMount } from "svelte";
    import Layout from "../lib/Layout.svelte";

    type Props = {
        valid: boolean
        reset: string
        username: string
    }

    let reset = $state("")
    let username = $state("")
    let password = $state("")
    let passwordConfirm = $state("")
    let valid = $state(false)

    onMount(async () => {
        const url = new URL(location.href)
        const res = await fetch(`/reset/${url.pathname.split("/").at(-1)}`, {
            method: "POST",
            credentials: "same-origin"
        })

        const data = await res.json() as Props

        valid = data.valid
        reset = data.reset
        valid = data.valid
    })

    async function resetPassword(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/api/reset-password", {
            method: "POST",
            credentials: "same-origin",
            body: JSON.stringify({
                username,
                reset,
                password,
                passwordConfirm
            })
        })

        const data = await res.json()

        if (data.success) location.pathname = "/"
    }
</script>

<Layout title="The Common Game" subtitle="Reset you password">
    <h1>Reset Password</h1>
    {#if valid}
        <form class="container" id="reset" onsubmit={resetPassword} method="POST" action="/api/reset-password">
            <input type="hidden" name="username" bind:value={username} />
            <input type="hidden" name="reset" bind:value={reset}/>
            <input type="password" name="password" placeholder="Password" bind:value={password} />
            <input type="password" name="password-confirm" placeholder="Confirm Password" bind:value={passwordConfirm} />
            <button type="submit">Reset Password</button>
        </form>
    {:else}
        <p>Invalid Reset Link</p>
        <a href="/">Go Back Home and Login</a>
    {/if}
</Layout>
