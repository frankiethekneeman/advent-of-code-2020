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

const TARGET = 2020

func solution(lines []string) RESULT_TYPE {
    starters := strings.Split(strings.Join(lines, ""), ",")
    lastSeen := make(map[int]int)
    prev := 0
    for i := 0; i < TARGET; i++ {
        say := 0
        if i < len(starters) {
            parsed, err := strconv.Atoi(starters[i])
            checkErr(err)
            say = parsed
        } else {
            lastIndex, ok := lastSeen[prev]
            if ok {
                say = (i - 1) - lastIndex
            } else {
                say = 0
            }
        }
        if (i > 0) {
            lastSeen[prev] = i - 1
        }
        prev = say
    }
    return prev;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        436,
        1,
        10,
        27,
        78,
        438,
        1836,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 15
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
