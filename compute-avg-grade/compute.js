// Get all rows of the table. Use 1 as the index for bachelor's.
const rows = document.querySelectorAll(".table-striped")[0].rows

// Total sum and total number of credits.
let sum = 0
let nr_credits = 0

header = rows[0]

let grade_pos = 0
let credit_pos = 0

for (let i = 0; i < header.cells.length; i++) {
    if (header.cells[i].innerHTML === " Nota") {
        grade_pos = i
    }
    if (header.cells[i].innerHTML === "Nr. Credite") {
        credit_pos = i
    }
}

rows.forEach((row, index) => {

    // header
    if (index === 0) {
        return
    }

    grade = row.cells[grade_pos].innerText
    credit = row.cells[credit_pos].innerText

    // Some subjects don't have grades, can be "passed" or "failed", or simply
    // not computed yet. The grade may also be set to the empty string, which
    // means that it wasn't set yet and should be ignored.
    if (isNaN(grade) || grade === "") {
        return
    }

    sum += +grade * +credit
    nr_credits += +credit

    // useful for debugging
    console.log(`${index}: Grade ${grade}, Credits ${credit}`)
})

console.log(`Avg: ${+sum / +nr_credits}`)