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

const (
    OUT_OF_SERVICE = "x"
)
func parseBusses(line string) []int {
    parts := strings.Split(line, ",")
    toReturn := make([]int, 0, len(parts))
    for _, part := range parts {
        if part == OUT_OF_SERVICE {
            continue
        }
        id, err := strconv.Atoi(part)
        checkErr(err)
        toReturn = append(toReturn, id)
    }
    return toReturn
}

func calcWait(bus int, arrival int) int {
    if arrival % bus == 0 {
        return 0
    }
    return bus - (arrival % bus)
}

func solution(lines []string) RESULT_TYPE {
    arrival, err := strconv.Atoi(lines[0])
    checkErr(err)
    busses := parseBusses(lines[1])
    firstBus := busses[0]
    for _, bus := range busses {
        if calcWait(bus, arrival) < calcWait(firstBus, arrival) {
            firstBus = bus
        }
    }
    return firstBus * calcWait(firstBus, arrival);
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        295,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 13
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
