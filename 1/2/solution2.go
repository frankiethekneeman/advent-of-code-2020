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

func solution(lines []string) RESULT_TYPE {
    TARGET := 2020
    lineItems := make(map[int]bool) // Poor Man's Set.
    for _, line := range lines {
        asInt, err := strconv.Atoi(line)
        checkErr(err)
        lineItems[asInt] = true
    }

    for firstItem := range lineItems {
        for secondItem := range lineItems {
            if secondItem == firstItem {
                continue
            }
            remainder := TARGET - firstItem - secondItem
            if lineItems[remainder] {
                return firstItem * secondItem * remainder
            }
        }
    }
    return -1;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        241861950,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 1
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
