<script lang="ts">
    import Layout from "../lib/Layout.svelte";

    let rusername = $state("")
    let rpassword = $state("")
    let lusername = $state("")
    let lpassword = $state("")
    let email = $state("")

    async function login(evt: Event) {
        evt.preventDefault()

        const res = await fetch("/auth/login", {
            method: "POST",
            body: JSON.stringify({
                username: lusername,
                password: lpassword
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
                username: rusername,
                password: rpassword,
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
        <input type="text" name="username" placeholder="Username" bind:value={lusername} />
        <input type="password" name="password" placeholder="Password" bind:value={lpassword} />
        <button type="submit">Login</button>
    </form>

    <h1>Sign Up</h1>
    <form onsubmit={register} class="container" id="register" method="POST" action="/api/register">
        <input type="text" name="email" placeholder="Email" bind:value={email} />
        <input type="text" name="username" placeholder="Username" bind:value={rusername} />
        <input type="password" name="password" placeholder="Password" bind:value={rpassword} />
        <button type="submit">Register</button>
    </form>
</Layout>

