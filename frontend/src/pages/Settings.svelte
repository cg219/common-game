<script lang="ts">
    import Layout from "../lib/Layout.svelte";

    let apikey = $state("")
    let apiname = $state("")
    let username = $state("")

    async function resetPassword() {
        const data = new URLSearchParams();

        data.append("username", username)

        await fetch("/api/forgot-password", {
            headers: {
                "Content-type": "application/x-www-form-urlencoded"
            },
            method: "POST",
            body: data
        })
    }

    async function generateKey() {
        const res = await fetch(`/api/generate-apikey/${apiname}`, { method: "POST" }).then((res) => res.json())

        apikey = res.apikey
    }
</script>

<Layout title="The Common Game" subtitle="My Settings" links={[]}>
    <form>
        <fieldset>
            <label for="reset-pass">Reset Password</label>
            <input type="text" name="username" placeholder="Username" bind:value={username} />
            <input type="button" onclick={resetPassword} name="reset-pass" value="Reset Password"/>
        </fieldset>
        <fieldset>
            <label for="new-key">New API Key</label>
            <input type="text" placeholder="Name" bind:value={apiname}>
            <input type="button" onclick={generateKey} name="api-generate" value="Generate">
            <input type="text" name="new-key" disabled value={apikey}>
        </fieldset>
    </form>
</Layout>
