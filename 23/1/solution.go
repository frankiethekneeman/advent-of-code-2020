package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type RESULT_TYPE = int

/*
Begin Solution
*/
const N_STEPS = 100

func solution(lines []string) RESULT_TYPE {
    currentCup := parse(strings.Join(lines, ""))
    index := makeIndex(currentCup)
    for i := 0; i < N_STEPS; i++ {
        currentCup = step(currentCup, index)
    }
    cup := index[1].clockwise
    result := 0
    multiplier := 10000000
    for cup.label != 1 {
        result += multiplier * cup.label
        cup = cup.clockwise
        multiplier /= 10

    }
    return result;
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
    if destinationLabel < 1 {
        destinationLabel = 9 // This is definitely going to change in part 2.
    }
    for removed[destinationLabel] != nil {
        destinationLabel --
        if destinationLabel < 1 {
            destinationLabel = 9 // This is definitely going to change in part 2.
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
        67384529,
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
