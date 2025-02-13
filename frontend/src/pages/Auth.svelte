<script lang="ts">
    import { onMount } from "svelte";
    import Layout from "../lib/Layout.svelte";

    let rusername = $state("")
    let rpassword = $state("")
    let lusername = $state("")
    let lpassword = $state("")
    let email = $state("")
    let showlogin = $state(true)
    let lusererr = $state("")
    let lpasserr = $state("")
    let rusererr = $state("")
    let rpasserr = $state("")
    let remailerr = $state("")
    let showconf = $state(false)

    async function login(evt: Event) {
        evt.preventDefault()

        if (lusername.trim().length == 0) lusererr = "true"
        if (lpassword.trim().length == 0) lpasserr = "true"
        if ([lusererr, lpasserr].includes("true")) return;

        const res = await fetch("/auth/login", {
            method: "POST",
            body: JSON.stringify({
                username: lusername,
                password: lpassword
            })
        })

        const data = await res.json()

        if (data.success) {
            location.pathname = "/game"
        } else {
            lusererr = "true"
            lpasserr = "true"
        }
    }

    async function register(evt: Event) {
        evt.preventDefault()

        if (rusername.trim().length == 0) rusererr = "true"
        if (rpassword.trim().length == 0) rpasserr = "true"
        if (email.trim().length == 0) remailerr = "true"
        if ([rusererr, rpasserr, remailerr].includes("true")) return;

        const res = await fetch("/auth/register", {
            method: "POST",
            body: JSON.stringify({
                username: rusername,
                password: rpassword,
                email
            })
        })

        const data = await res.json()

        if (data.success) {
            showconf = true;
            setTimeout(() => showconf = false, 4000)
        } else {
            if (data.code == 0) {
                rusererr = "true"
                return
            }

            remailerr = "true"
            rusererr = "true"
            rpasserr = "true"
        }
    }

    function toggleLogin(val: string) {
        switch (val) {
            case "#login":
                showlogin = true
                break;
            case "#register":
                showlogin = false
                break;
        }
    }

    function clearValidation(val: string) {
        switch (val) {
            case "luser":
                lusererr = ""
                break;
            case "lpass":
                lpasserr = ""
                break;
            case "ruser":
                rusererr = ""
                break;
            case "rpass":
                rpasserr = ""
                break;
            case "remail":
                remailerr = ""
                break;
        }
    }

    onMount(() => {
        toggleLogin(window.location.hash)

        window.addEventListener("hashchange", (_) => {
            toggleLogin(window.location.hash)
        })
    })

</script>

<Layout title="The Common Game" subtitle="">

    <article id="login" style:display={showlogin ? "block" : "none"}>
        <header>Login</header>
        <form onsubmit={login} class="container" id="login" method="POST" action="/api/login">
            <input type="text" name="username" aria-invalid={lusererr} onfocus={() => clearValidation("luser")} placeholder="Username" bind:value={lusername} />
            <input type="password" name="password" aria-invalid={lpasserr} onfocus={() => clearValidation("lpass")} placeholder="Password" bind:value={lpassword} />
            <button type="submit">Login</button>
        </form>
        <footer>
            or <a class="secondary" href="#register">sign up here</a>
        </footer>
    </article>

    <article id="register" style:display={showlogin ? "none" : "block"}>
        <header>Sign Up</header>
        <form onsubmit={register} class="container" id="register" method="POST" action="/api/register">
            <input type="text" name="email" aria-invalid={remailerr} onfocus={() => clearValidation("remail")}  placeholder="Email" bind:value={email} />
            <input type="text" name="username" aria-describedby="username-invalid" aria-invalid={rusererr} onfocus={() => clearValidation("ruser")}  placeholder="Username" bind:value={rusername} />
            <small id="username-invalid">
                {#if rusername.trim().length > 0 && rusererr == "true"}
                   Username Taken. Choose another one. 
                {/if} 
            </small>
            <input type="password" name="password" aria-invalid={rpasserr} onfocus={() => clearValidation("rpass")}  placeholder="Password" bind:value={rpassword} />
            <button type="submit">Register</button>
        </form>
        <footer>
            or <a class="secondary" href="#login">login here</a>
        </footer>
    </article>
    <dialog open={showconf}>
        <article>
            Check your email to validate your account and login
        </article>
    </dialog>
</Layout>

