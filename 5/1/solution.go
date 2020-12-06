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

func interpretAsBinary(encodedNumber string, on rune, off rune) int {
    runes := []rune(encodedNumber)
    total := 0
    for curr, val := len(runes) - 1, 1; curr >= 0; curr, val = curr - 1, val * 2 {
        if runes[curr] == on {
            total += val
        } else if runes[curr] != off {
            panic("Unexpected Rune: " + string(runes[curr]))
        }
    }
    return total
}

func solution(lines []string) RESULT_TYPE {
    maxId := 0
    for _, line := range lines {
        row := interpretAsBinary(line[:7], 'B', 'F')
        column := interpretAsBinary(line[7:], 'R', 'L')
        id := row * 8 + column
        if id > maxId {
            maxId = id
        }
    }
    return maxId;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        357,
        567,
        119,
        820,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 5
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
