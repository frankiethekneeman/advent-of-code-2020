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
func xor(A bool, B bool) bool {
    return (A || B) && !(A && B)
}

func validate(pw string, target rune, pos1 int, pos2 int) bool {
    runes := []rune(pw);
    return xor(runes[pos1 - 1] == target, runes[pos2 - 1] == target)
}

func solution(lines []string) RESULT_TYPE {
    valid := 0
    for _, line := range lines {
        var pos1 int
        var pos2 int
        var character rune
        var password string
        fmt.Sscanf(line, "%d-%d %c: %s", &pos1, &pos2, &character, &password)
        if validate(password, character, pos1, pos2) {
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
        0,
        1,
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
