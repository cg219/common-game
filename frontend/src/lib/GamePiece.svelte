<script lang="ts">
    import { getContext } from "svelte";

    type Props = {
        value: string
        correct: boolean
        wrong: boolean
        selected: boolean
        subject?: number
    }

    let { value, correct, wrong, subject, selected = false } : Props = $props()
    let select = getContext("select") as Function

    const pieceStyle = "w-full aspect-square text-sm whitespace-nowrap overflow-hidden rounded-lg transition-colors duration-200 cursor-pointer text-zinc-100";
</script>

<button
    class="{pieceStyle} {selected ? "bg-teal-100! text-slate-800!" : "bg-teal-800 hover:bg-teal-700"}
    {correct && subject == 0 ? "bg-blue-300! text-slate-800!" : ""}
    {correct && subject == 1 ? "bg-pink-300! text-slate-800!" : ""}
    {correct && subject == 2 ? "bg-amber-300! text-slate-800!" : ""}
    {correct && subject == 3 ? "bg-emerald-300! text-slate-800!" : ""}
    {wrong ? "bg-red-300!" : ""}"

    onclick={() => {
        if (!correct) {
            select(value, selected)}
        }
    }>
    {value}
</button>

<style>
    @import 'tailwindcss';
</style>
