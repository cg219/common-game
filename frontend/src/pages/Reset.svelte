<script lang="ts">
    import { onMount } from "svelte";
    import Layout from "../lib/Layout.svelte";
    import { buttonStyle, formStyle, inputStyle, linkStyle, setStyle } from "../lib/styles";

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

<Layout title="The Common Game" subtitle="Reset your password">
    {#if valid}
        <form class={formStyle} id="reset" onsubmit={resetPassword} method="POST" action="/api/reset-password">
            <input type="hidden" name="username" bind:value={username} />
            <input type="hidden" name="reset" bind:value={reset}/>
            <fieldset class={setStyle}>
                <input class={inputStyle} type="password" name="password" placeholder="Password" bind:value={password} />
                <input class={inputStyle} type="password" name="password-confirm" placeholder="Confirm Password" bind:value={passwordConfirm} />
                <button class={buttonStyle} type="submit">Reset Password</button>
            </fieldset>
        </form>
    {:else}
        <div class="mx-auto text-center">
            <p class="text-zinc-100 text-lg">Invalid Reset Link</p>
            <a class={linkStyle} href="/">Go Back Home and Login</a>
        </div>
    {/if}
</Layout>
