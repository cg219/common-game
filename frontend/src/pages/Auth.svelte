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

    const cardStyle = "rounded-lg border-zinc-700/50 border-1 border-solid w-xl mx-auto bg-zinc-900 text-zinc-100 overflow-auto"
    const inputStyle = "bg-zinc-800 p-5 text-md my-2 border-zinc-600/30 border-1 border-solid focus:border-teal-800/80 rounded-lg outline-none";
    const buttonStyle = "bg-teal-800 rounded-lg p-5 my-2 cursor-pointer";
    const headerStyle = "text-lg mb-5 px-10 bg-zinc-800 py-5"
    const formStyle = "flex flex-col justify-even w-full px-10"
    const footerStyle = "w-full px-10 py-4 mt-5 bg-zinc-800"
    const linkStyle = "underline text-slate-400 hover:decoration-teal-500 hover:text-slate-300 transition-colors duration-200"

</script>

<Layout title="The Common Game" subtitle="">
    <article id="login" style:display={showlogin ? "block" : "none"} class={cardStyle}>
        <header class={headerStyle}>Login</header>
        <form onsubmit={login} class={formStyle} id="login" method="POST" action="/api/login">
            <input type="text" name="username" aria-invalid={lusererr} onfocus={() => clearValidation("luser")} placeholder="Username" bind:value={lusername} class={inputStyle} />
            <input type="password" name="password" aria-invalid={lpasserr} onfocus={() => clearValidation("lpass")} placeholder="Password" bind:value={lpassword} class={inputStyle} />
            <button type="submit" class={buttonStyle}>Login</button>
        </form>
        <footer class={footerStyle}>
            or <a class={linkStyle} href="#register">sign up here</a>
        </footer>
    </article>

    <article id="register" style:display={showlogin ? "none" : "block"} class={cardStyle}>
        <header class={headerStyle}>Sign Up</header>
        <form onsubmit={register} class={formStyle} id="register" method="POST" action="/api/register">
            <input type="text" class={inputStyle} name="email" aria-invalid={remailerr} onfocus={() => clearValidation("remail")}  placeholder="Email" bind:value={email} />
            <input type="text" class={inputStyle} name="username" aria-describedby="username-invalid" aria-invalid={rusererr} onfocus={() => clearValidation("ruser")}  placeholder="Username" bind:value={rusername} />
            <small id="username-invalid">
                {#if rusername.trim().length > 0 && rusererr == "true"}
                   Username Taken. Choose another one. 
                {/if} 
            </small>
            <input type="password" class={inputStyle} name="password" aria-invalid={rpasserr} onfocus={() => clearValidation("rpass")}  placeholder="Password" bind:value={rpassword} />
            <button type="submit" class={buttonStyle}>Register</button>
        </form>
        <footer class={footerStyle}>
            or <a class={linkStyle} href="#login">login here</a>
        </footer>
    </article>
    <dialog open={showconf}>
        <article>
            Check your email to validate your account and login
        </article>
    </dialog>
</Layout>

