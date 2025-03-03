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

    const pieceStyle = "w-full aspect-square text-[70%] sm:text-lg whitespace-wrap overflow-hidden rounded-lg transition-colors duration-200 cursor-pointer text-slate-700";
</script>

<button
    class="{pieceStyle} {selected ? "bg-teal-800! text-slate-200!" : "bg-slate-300 hover:bg-slate-500 hover:text-slate-200!"}
    {correct && subject == 0 ? "bg-blue-300! text-slate-800! hover:text-slate-800!" : ""}
    {correct && subject == 1 ? "bg-pink-300! text-slate-800! hover:text-slate-800!" : ""}
    {correct && subject == 2 ? "bg-amber-300! text-slate-800! hover:text-slate-800!" : ""}
    {correct && subject == 3 ? "bg-emerald-300! text-slate-800! hover:text-slate-800!" : ""}
    {wrong ? "bg-red-500! text-slate-200! hover:text-slate-200!" : ""}"

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
