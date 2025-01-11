/// <reference types="svelte" />
/// <reference types="vite/client" />

import { mount } from 'svelte'
import Reset from "./pages/Reset.svelte"
import Auth from "./pages/Auth.svelte"

const pages = new Map()

pages.set("reset", Reset)
pages.set("auth", Auth)

const app = (page: string) => {
    let p = pages.get("*");

    if (pages.has(page)) p = pages.get(page)

    mount(p, {
        target: document.getElementById('app')!,
    })
}

export default app

