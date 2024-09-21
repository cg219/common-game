async function register() {
    const client = new webauthn.WebAuthnClient()
    const username = htmx.values(htmx.find("#register")).get("username")
    const res = await fetch("/auth/register", {
        method: "POST",
        body: JSON.stringify({ username })
    })

    const data = await res.json()
    const pub = await client.register(data)
    await fetch("/auth/verify", {
        method: "POST",
        body: JSON.stringify({
            username,
            response: pub
        })
    })
    window.location.href = "/"
}

async function login() {
    const client = new webauthn.WebAuthnClient()
    const username = htmx.values(htmx.find("#login")).get("username")
    const res = await fetch("/auth", {
        method: "POST",
        body: JSON.stringify({ username })
    })

    const data = await res.json()
    const pub = await client.authenticate(data)

    await fetch("/auth/auth-verify", {
        method: "POST",
        body: JSON.stringify({
            username,
            response: pub
        })
    })
    
    window.location.href = "/"
}

window.commongame = {
    register, login
}
