package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

type RESULT_TYPE = int64

/*
Begin Solution
*/

func countPaths(input int, target int, adapters map[int] bool, memos map[int]int64) int64 {
    count, memoized := memos[input];
    if memoized {
        return count
    }
    count = 0
    for next := input + 1; next <= input + 3; next++ {
        if next == target {
            count++
        } else if adapters[next] {
            count += countPaths(next, target, adapters, memos)
        }
    }
    memos[input] = count
    return count

}

func solution(lines []string) RESULT_TYPE {
    adapters := make(map[int]bool)
    max := 0
    for _, line := range lines {
        output, err := strconv.Atoi(line)
        checkErr(err)
        adapters[output] = true
        if output > max {
            max = output
        }
    }
    
    return countPaths(0, max, adapters, make(map[int]int64))
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        8,
        19208,
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
