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

const (
    TREE = '#'
    OPEN = '.'
)

func solution(lines []string) RESULT_TYPE {
    treeLocations := make([][]rune, len(lines))
    for i, line := range lines {
        treeLocations[i] = []rune(line)
    }
    row := 0
    col := 0
    DROP := 1
    TRAVEL := 3
    trees := 0
    for ; row < len(lines); row += DROP {
        if treeLocations[row][col] == TREE {
            trees++
        }
        col = (col + TRAVEL) % len(treeLocations[row])
    }
    return trees;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        7,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 3
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
