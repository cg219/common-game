async function getRegPKCreds() {
    const client = new webauthn.WebAuthnClient()
    const username = htmx.values(htmx.find("#register")).get("username")
    const res = await fetch("/auth/register", {
        method: "POST",
        body: JSON.stringify({ username })
    })

    const data = await res.json()
    console.log(data)

    const pub = await client.register(data)

    console.log(pub)

    const res2 = await fetch("/auth/verify", {
        method: "POST",
        body: JSON.stringify(pub)
    })

    console.log(res2)
}

async function getAuthPKCreds() {
    const client = new webauthn.WebAuthnClient()
    const username = htmx.values(htmx.find("#login")).get("username")
    const res = await fetch("/auth", {
        method: "POST",
        body: JSON.stringify({ username })
    })

    const data = await res.json()
    console.log(data)

    const pub = await client.authenticate(data)

    console.log(pub)

    const res2 = await fetch("/auth/auth-verify", {
        method: "POST",
        body: JSON.stringify(pub)
    })

    console.log(res2)
}

window.getRegPKCreds = getRegPKCreds
window.getAuthPKCreds = getAuthPKCreds
