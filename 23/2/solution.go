package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type RESULT_TYPE = int64

/*
Begin Solution
*/
const N_STEPS = 10_000_000
const MAX_CUP = 1_000_000
const MIN_CUP = 1

func solution(lines []string) RESULT_TYPE {
    currentCup := parse(strings.Join(lines, ""))
    fill(10, MAX_CUP, currentCup)

    index := makeIndex(currentCup)

    for i := 0; i < N_STEPS; i++ {
        currentCup = step(currentCup, index)
    }
    star1 := index[1].clockwise.label
    star2 := index[1].clockwise.clockwise.label
    return int64(star1) * int64(star2)
}

func step(currentCup *Cup, index CupIndex) *Cup {
    pickup := currentCup.clockwise
    currentCup.clockwise = pickup.clockwise.clockwise.clockwise //pickup 3 cups.
    pickup.clockwise.clockwise.clockwise = nil
    destinationCup := getDestination(currentCup.label, index, makeIndex(pickup))
    pickup.clockwise.clockwise.clockwise = destinationCup.clockwise
    destinationCup.clockwise = pickup
    return currentCup.clockwise
}

func parse(line string) *Cup {
    runes := []rune(line)
    firstCup := cup(int(runes[0] - '0'))
    curr := firstCup
    for _, r := range runes[1:] {
        curr.clockwise = cup(int(r - '0'))
        curr = curr.clockwise
    }
    curr.clockwise = firstCup
    return firstCup
}

func fill(firstToAdd int, lastToAdd int, firstCup *Cup) {
    curr := firstCup.clockwise
    for curr.clockwise.label != firstCup.label {
        curr = curr.clockwise
    }
    for add := firstToAdd; add <= lastToAdd; add ++ {
        curr.clockwise = cup(add)
        curr = curr.clockwise
    }
    curr.clockwise = firstCup
}

func makeIndex(current *Cup) CupIndex {
    index := make(CupIndex)
    for current != nil && index[current.label] == nil {
        index[current.label] = current
        current = current.clockwise
    }
    return index
}

func getDestination(currentLabel int, index CupIndex, removed CupIndex) *Cup {
    destinationLabel := currentLabel - 1
    if destinationLabel < MIN_CUP {
        destinationLabel = MAX_CUP
    }
    for removed[destinationLabel] != nil {
        destinationLabel --
        if destinationLabel < MIN_CUP {
            destinationLabel = MAX_CUP
        }
    }
    return index[destinationLabel]
}

type Cup struct {
    label int
    clockwise *Cup
}

type CupIndex = map[int]*Cup

func cup(label int) *Cup {
    return &Cup{label, nil}
}
/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        149245887792,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 23
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
