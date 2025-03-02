<script lang="ts">
    import Layout from "../lib/Layout.svelte";
    import { formStyle, setStyle, inputStyle, buttonStyle, disabledStyle } from "../lib/styles";

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
    <form class={formStyle}>
        <fieldset class={setStyle}>
            <label for="reset-pass">Reset Password</label>
            <input type="text" class={inputStyle} name="username" placeholder="Username" bind:value={username} />
            <input type="button" class={buttonStyle} onclick={resetPassword} name="reset-pass" value="Reset Password"/>
        </fieldset>
        <fieldset class={setStyle}>
            <label for="new-key">New API Key</label>
            <input type="text" class={inputStyle} placeholder="Name" bind:value={apiname}>
            <input type="button" class={buttonStyle} onclick={generateKey} name="api-generate" value="Generate">
            <input type="text" class={apikey == "" ? inputStyle : disabledStyle} placeholder="Key" name="new-key" disabled value={apikey}>
        </fieldset>
    </form>
</Layout>
