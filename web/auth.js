async function getPKCreds() {
    const username = htmx.values(htmx.find("#register")).get("username")
    const res = await fetch("/auth/register", {
        method: "POST",
        body: JSON.stringify({ username })
    })

    const data = await res.json()

    const publicKey = {
        challenge: new TextEncoder().encode(atob(data.challenge)),
        rp: {
            name: "The Common Game",
            id: "localhost"
        },
        user: {
            id: new TextEncoder().encode(atob(data.userid)),
            name: data.username,
            displayName: data.displayName
        },
        pubKeyCredParams: [{alg: -7, type: "public-key"},{alg: -257, type: "public-key"}],
        authenticatorSelection: {
            authenticatiorAttachment: "platform",
            requireResidentKey: true
        }
    }

    console.log(publicKey)

    const credential = await navigator.credentials.create({ publicKey })

    console.log(credential)
}
