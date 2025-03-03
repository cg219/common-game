<script lang="ts">
import type { Snippet } from "svelte";
import type { Link } from "./customtypes";

type Props = {
    title: string
    subtitle: string
    links?: Link[]
    children: Snippet
}

let { title, subtitle, links, children }: Props = $props();

async function logout(evt: Event) {
    evt.preventDefault()

    const res = await fetch("/auth/logout", {
        method: "POST",
        credentials: "same-origin"
    })

    const data = await res.json()

    if (data.success) location.pathname = "/"
}

const navlinkStyle = "text-zinc-200 text-md no-underline hover:underline hover:decoration-teal-500";
const selectedlinkStyle = `${navlinkStyle} text-teal-500`
</script>

<main class="mx-auto w-full font-[Lato]">
    <nav class="flex flex-col sm:flex-row w-9/10 gap-10 justify-between mx-auto my-10">
        <ul>
            <li>
                <hgroup>
                    <h1 class="text-4xl text-teal-800 hover:text-teal-700 transition-colors duration-200 font-[Bowlby_One_SC] font-normal"><a class="" href="/">{title}</a></h1>
                    <p class="text-sm sm:text-base text-slate-400">{subtitle}</p>
                </hgroup>
            </li>
        </ul>
        {#if links}
            <ul class="flex flex-row gap-5 sm:gap-8 justify-start sm:justify-evenly">
                {#each links as { current, url, name }}
                    {#if current}
                        <li class={selectedlinkStyle}><a href="{url}" aria-current="page">{name}</a></li>
                    {:else}
                        <li class={navlinkStyle}><a href="{url}">{name}</a></li>
                    {/if} 
                {/each}
                <li class={navlinkStyle}>
                    <a href="/account">Account</a>
                </li>
                <li class={navlinkStyle}>
                    <a href="/report">Report a Bug</a>
                </li>
                <li class={navlinkStyle}>
                    <a onclick={logout} href="#logout" class="contrast">Logout</a>
                </li>
            </ul>
        {/if}
    </nav>

    <section class="overflow-hidden w-9/10 m-auto mt-10">
        {@render children()}
    </section>
</main>

<style>
    @import 'tailwindcss';
</style>
