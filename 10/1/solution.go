package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strconv"
)

type RESULT_TYPE = int

/*
Begin Solution
*/

func solution(lines []string) RESULT_TYPE {
    adapters := make([]int, len(lines))
    for i, line := range lines {
        asInt, err := strconv.Atoi(line)
        checkErr(err)
        adapters[i] = asInt
    }
    sort.Ints(adapters)
    diff1 := 0
    diff3 := 1
    for i, output := range adapters {
        input := 0
        if (i > 0) {
            input = adapters[i - 1]
        }
        switch output - input {
            case 1: diff1 ++
            case 3: diff3 ++
        }
    }
    return diff1 * diff3;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        35,
        220,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 10
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
