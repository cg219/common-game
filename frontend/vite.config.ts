import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { resolve } from "@std/path/resolve"

export default defineConfig({
  plugins: [svelte()],
    build: {
        rollupOptions: {
            input: {
                game: resolve(import.meta.dirname!,  "entrypoints/game.html"),
                auth: resolve(import.meta.dirname!,  "entrypoints/auth.html"),
                reset: resolve(import.meta.dirname!,  "entrypoints/reset.html"),
                report: resolve(import.meta.dirname!,  "entrypoints/report.html"),
                settings: resolve(import.meta.dirname!,  "entrypoints/settings.html")
            },
            output: {
                dir: "../static-app"
            }
        }
    }
})
