package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
)

type RESULT_TYPE = int

/*
Begin Solution
*/
var REQUIRED_FIELDS = [...]string {
    "byr",
    "iyr",
    "eyr",
    "hgt",
    "hcl",
    "ecl",
    "pid",
}

func isValid(passport map[string]string) bool {
    for _, field := range REQUIRED_FIELDS {
        _, exists := passport[field]
        if !exists {
            return false
        }
    }
    return true
}

func addFields(passport map[string] string, line string) {
    for _, entry := range strings.Split(line, " ") {
        pieces := strings.Split(entry, ":")
        passport[pieces[0]] = pieces[1]
    }
}

func solution(lines []string) RESULT_TYPE {
    passport := make(map[string] string)
    valid := 0

    for _, line := range lines {
        if line == "" {
            if isValid(passport) {
                valid ++
            }
            passport = make(map[string] string)
        } else {
            addFields(passport, line)
        }
    }
    if isValid(passport) {
        valid ++
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
        0,
        2,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 4
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
