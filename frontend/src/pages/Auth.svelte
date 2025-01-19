<script lang="ts">
    import Layout from "../lib/Layout.svelte";

    let username = $state("")
    let password = $state("")
    let email = $state("")

    async function login(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/auth/login", {
            method: "POST",
            body: JSON.stringify({
                username,
                password
            })
        })

        const data = await res.json()

        if (data.success) location.pathname = "/game"
    }

    async function register(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/auth/register", {
            method: "POST",
            body: JSON.stringify({
                username,
                password,
                email
            })
        })

        const data = await res.json()

        if (data.success) location.pathname = "/game"
    }
</script>

<Layout title="The Common Game" subtitle="Login">
    <h1>Sign In</h1>
    <form onsubmit={login} class="container" id="login" method="POST" action="/api/login">
        <input type="text" name="username" placeholder="Username" bind:value={username} />
        <input type="password" name="password" placeholder="Password" bind:value={password} />
        <button type="submit">Login</button>
    </form>

    <h1>Sign Up</h1>
    <form onsubmit={register} class="container" id="register" method="POST" action="/api/register">
        <input type="text" name="email" placeholder="Email" bind:value={email} />
        <input type="text" name="username" placeholder="Username" bind:value={username} />
        <input type="password" name="password" placeholder="Password" bind:value={password} />
        <button type="submit">Register</button>
    </form>
</Layout>

