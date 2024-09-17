async function getPKCreds() {
    const username = htmx.values(htmx.find("#register")).get("username")

    const res = await fetch("/auth/register", {
        method: "POST",
        body: JSON.stringify({ username })
    })

    console.log(res)
}
