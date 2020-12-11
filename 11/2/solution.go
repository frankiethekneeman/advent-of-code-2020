package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

type RESULT_TYPE = int

/*
Begin Solution
*/

type SeatingChart = [][]rune

const (
    FREE = 'L'
    FLOOR = '.'
    FILLED = '#'
    DEBUG = true
)

func getFirstVisible(r int, c int, dr int, dc int, chart SeatingChart) rune {
    if dr == 0 && dc == 0 {
        //This is a non-extensible vector.  Garbage in, garbage out.
        return FLOOR
    }
    for m := 1; true; m++ {
        if r + m * dr < 0 || r + m * dr >= len(chart) {
            // if we've run off the row...
            return FLOOR
        }
        if c + m * dc < 0 || c + m * dc >= len(chart[r + m * dr]) {
            // if we've run off the column...
            return FLOOR
        }
        space := chart[r + m * dr][c + m * dc]
        if space != FLOOR {
            return space
        }
    }
    panic("this line can't be executed")
}
func countNeighbors(r int, c int, chart SeatingChart) int {
    count := 0
    for dr := -1; dr < 2; dr++ {
        for dc := -1; dc < 2; dc++ {
            if getFirstVisible(r, c, dr, dc, chart) == FILLED {
                count++
            }
        }
    }
    return count
}

func increment(in SeatingChart) SeatingChart {
    output := make(SeatingChart, len(in))
    for r, row := range in {
        output[r] = make ([]rune, len(row))
        for c, seat := range row {
            neighbors := countNeighbors(r, c, in)
            if seat == FREE && neighbors == 0 {
                output[r][c] = FILLED
            } else if seat == FILLED && neighbors >= 5 {
                output[r][c] = FREE
            } else {
                output[r][c] = seat
            }
        }
    }
    return output
}

func isSame(left SeatingChart, right SeatingChart) bool {
    if left == nil && right == nil {
        return true
    }
    if left == nil || right == nil {
        return false
    }
    if len(left) != len(right) {
        return false
    }
    for r, row := range left {
        if len(row) != len(right[r]) {
            return false
        }
        for c, seat := range row {
            if right[r][c] != seat {
                return false
            }
        }
    }
    return true
}
func debug(in SeatingChart) {
    if DEBUG {
        for _, row := range in {
            println(string(row))
        }
        println("")
    }
}

func getStable(in SeatingChart) SeatingChart {
    var prev SeatingChart = nil
    curr := in
    for !isSame(prev, curr) {
        prev = curr
        curr = increment(curr)
    }
    return curr
}


func solution(lines []string) RESULT_TYPE {
    in := make(SeatingChart, len(lines))
    for i, line := range lines {
        in[i] = []rune(line)
    }
    stable := getStable(in)
    countFilled := 0
    for _, row := range stable {
        for _, seat := range row {
            if seat == FILLED {
                countFilled++
            }
        }
    }
    return countFilled
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        26,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 11
    passedTests := true
    for i, expected := range TEST_CASES(){
        result := solution(getTest(DAY, i + 1))
        if result != expected {
            fmt.Printf("Test %d: expected %v, but got %v instead!\n", i+1, expected, result)
            passedTests = false
        }
    }
    if passedTests {
        fmt.Printf("%v\n", solution(getInput(DAY)))
    }
}

func getInput(day int) []string {
    return getFile(strconv.Itoa(day) + "/input")
}
func getTest(day int, test int) []string {
    return getFile(strconv.Itoa(day) + "/test" + strconv.Itoa(test))
}

func getFile(location string) []string {
    file, err := os.Open(location)
    checkErr(err)
    scanner := bufio.NewScanner(file)
    var lines []string
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    checkErr(scanner.Err())
    return lines
}
