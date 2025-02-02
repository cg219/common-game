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
</script>

<button
    class:selected={selected}
    class:correct={correct}
    class:wrong={wrong}
    class:s0={subject == 0}
    class:s1={subject == 1}
    class:s2={subject == 2}
    class:s3={subject == 3}

    onclick={() => {
        if (!correct) {
            select(value, selected)}
        }
    }>
    {value}
</button>

<style>
    button {
        width: 100%;
        height: 100%;

        &.selected {
            background-color: #CFCFCF;
        }

        &.wrong {
            background-color: #EF0038;
        }

        &.correct.s0 {
            background-color: lightblue;
            color: black;
        }

        &.correct.s1 {
            background-color: lightpink;
            color: black;
        }

        &.correct.s2 {
            background-color: lightgoldenrodyellow;
            color: black;
        }

        &.correct.s3 {
            background-color: lightsteelblue;
            color: black;
        }
    }
</style>
