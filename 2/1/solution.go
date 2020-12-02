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

func count(pw string, target rune) int {
    matching := 0
    for _, char := range pw {
        if char == target {
            matching ++
        }
    }
    return matching
}
func solution(lines []string) RESULT_TYPE {
    valid := 0
    for _, line := range lines {
        var min int
        var max int
        var character rune
        var password string
        fmt.Sscanf(line, "%d-%d %c: %s", &min, &max, &character, &password)
        targetCount := count(password, character)
        if targetCount >= min && targetCount <= max {
            valid++
        }
    }
    return valid;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        1,
        0,
        1,
        2,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 2
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
