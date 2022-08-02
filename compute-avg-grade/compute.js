// Get all rows of the table.
const rows = document.querySelectorAll(".table-striped")[1].rows

// Total sum and total number of credits.
let sum = 0
let nr_credits = 0

rows.forEach( (row, index) => {

    grade = row.cells[5].innerText
    credit = row.cells[6].innerText

    // Some subjects don't have grades, can be "passed" or "failed", or simply
    // not computed yet.
    if (isNaN(grade)){
        return
    }

    sum += +grade * +credit
    nr_credits += +credit

    // useful for debugging
    // console.log(`${index}: Grade ${grade}, Credits ${credit}`)
})

console.log(`Avg: ${+sum / +nr_credits}`)