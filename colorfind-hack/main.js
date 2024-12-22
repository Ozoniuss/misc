function distance(s1, s2) {
    const l1 = s1.indexOf('(')
    const r1 = s1.indexOf(')')
    const ss1 = s1.slice(l1+1, r1)
    
    const l2 = s2.indexOf('(')
    const r2 = s2.indexOf(')')
    const ss2 = s2.slice(l2+1, r2)



    const p1 = ss1.split(", ")
    const p2 = ss2.split(", ")
    
    return Math.sqrt(
        (parseFloat(p1[0]) - parseFloat(p2[0])) * (parseFloat(p1[0]) - parseFloat(p2[0])) +
        (parseFloat(p1[1]) - parseFloat(p2[1])) * (parseFloat(p1[1]) - parseFloat(p2[1])) +
        (parseFloat(p1[2]) - parseFloat(p2[2])) * (parseFloat(p1[2]) - parseFloat(p2[2]))
    )
}
async function findUnique() {
    const cells = document.getElementsByClassName("grid-cell")
    const colorCount = {}

    for (const c of cells) {
        const bgColor = c.style["background-color"]
        colorCount[bgColor] = (colorCount[bgColor] || 0) + 1
    }

    let same = null
    let different = null
    for (const color in colorCount) {
        if (colorCount[color] === 1) {
            different = color
        } else {
            same = color
        }
    }
    console.log("distance", distance(same, different))

    let toclick = null
    for (const c of cells) {
        if (c.style["background-color"] === different) {
            toclick = c
            break
        }
    }

    toclick.click()
    await new Promise(resolve => setTimeout(resolve, 100))
}

async function run() {
    const n = 1_000_000
    for (let i = 0; i < n; i++) {
        await findUnique()
    }
}

// run()
