const moves = []

function toggleSelection(m, word) {
    var i = 0
    var next = 0
    var len = 0
    var shouldAdd = true
    var nextSet = false

    m.forEach((w) => {
        if (w == word) {
            m[i] = undefined
            shouldAdd = false
        } else if (w != undefined) {
            if (!nextSet) {
                next++
            }

            len++
        } else {
            if (!nextSet) {
                nextSet = true
                next = i
            }
        }

        i++
    })

    if (shouldAdd && len < 4) {
        m[next] = word
        len++
    }

    return len
}

function reset(m) {
    m.forEach((_, i) => m[i] = undefined)
}

function removeWrong() {
    htmx.findAll(".game-board button.wrong").forEach((b) => b.classList.remove("wrong"))
}

htmx.find(".game").addEventListener("click", function (evt) {
    const piece = evt.target.closest(".game-board button")

    if (piece && !piece.classList.contains("correct")) {
        removeWrong()
        const word = piece.attributes.getNamedItem("data-value").value

        const selected = toggleSelection(moves, word)
        piece.classList.toggle("selected")
        console.log(`Moves: ${moves}`)

        if (selected == 4) {
            console.log(moves)
            htmx.trigger(".game-board", "moveSelected", { words: moves })
            reset(moves)
        }
    }
})
