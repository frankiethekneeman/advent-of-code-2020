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

func parse(lines []string) []int {
    toReturn := make([]int, len(lines))
    for i, line := range lines {
        parsed, err := strconv.Atoi(line)
        checkErr(err)
        toReturn[i] = parsed
    }

    return toReturn
}

func add(set map[int]int, value int) map[int]int {
    curr, exists := set[value]
    if exists {
        set[value] = curr + 1
    } else {
        set[value] = 1
    }
    return set
}

func remove(set map[int]int, value int) map[int]int {
    curr, exists := set[value]
    if !exists {
        panic(strconv.Itoa(value) + " is not a valid key in the provied map.")
    } else if curr == 1 {
        delete(set, value)
    } else {
        set[value] = curr - 1
    }
    return set
}

func canSum(numbers map[int] int, target int) bool {
    for first, count := range numbers {
        if 2 * first == target {
            if count > 1 {
                return true
            } else {
                continue
            }
        }
        second := target - first
        _, ok := numbers[second]
        if ok {
            return true
        }
    }
    return false
}

func getInvalidNumber(numbers []int, preambleLength int) int {
    sumFrom := make(map[int]int)
    for i := 0; i < preambleLength; i++ {
        sumFrom = add(sumFrom, numbers[i])
    }
    for i := preambleLength; i < len(numbers); i++ {
        if !canSum(sumFrom, numbers[i]) {
            return numbers[i]
        }
        sumFrom = add(remove(sumFrom, numbers[i-preambleLength]), numbers[i])
    }
    panic("failed to find invalid number")
}

func findContiguousSum(numbers []int, target int64) (int, int) {
    sum := int64(0)
    partials := make(map[int64] int)
    for i, number := range numbers {
        sum += int64(number)
        if sum < 0 {
            panic ("int64 overflow")
        }
        partials[sum] = i
    }
    for endSum, end := range partials {
        if endSum <= target {
            continue
        }
        preSum := endSum - target
        beginning, ok := partials[preSum]
        if ok && end - beginning > 1 {
            return beginning + 1, end
        }
    }
    panic("No contiguous sequence found.")
}


func solution(lines []string, preambleLength int) RESULT_TYPE {
    numbers := parse(lines)
    target := getInvalidNumber(numbers, preambleLength)
    begin, end := findContiguousSum(numbers, int64(target))
    min := numbers[begin]
    max := numbers[begin]
    for i := begin; i <= end; i++ {
        if numbers[i] > max {
            max = numbers[i]
        }
        if numbers[i] < min {
            min = numbers[i]
        }
    }
    return min + max
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        62,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 9
    passedTests := true
    for i, expected := range TEST_CASES(){
        result := solution(getTest(DAY, i + 1), 5)
        if result != expected {
            fmt.Printf("Test %d: expected %v, but got %v instead!\n", i+1, expected, result)
            passedTests = false
        }
    }
    if passedTests {
        fmt.Printf("%v\n", solution(getInput(DAY), 25))
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
