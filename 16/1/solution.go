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
type Field struct {
    name string
    min1 int
    max1 int
    min2 int
    max2 int
}

func parseField(line string) Field{
    parts := strings.Split(line, ": ")
    var min1 int
    var max1 int
    var min2 int
    var max2 int
    _, err := fmt.Sscanf(parts[1], "%d-%d or %d-%d", &min1, &max1, &min2, &max2)
    checkErr(err)
    return Field {
        parts[0],
        min1,
        max1,
        min2,
        max2,
    }
}

func assertInput(actual string, expected string) {
    if actual != expected {
        panic("Unexpected input... wanted '" + expected + "', but got '" + actual + "'")
    }
}

func parseTicket(line string) []int {
    parts := strings.Split(line, ",")
    toReturn := make([]int, len(parts))
    for i, part := range parts {
        parsed, err := strconv.Atoi(part)
        checkErr(err)
        toReturn[i] = parsed
    }
    return toReturn
}

func isValid(value int, fields []Field) bool {
    for _, field := range fields {
        if (value >= field.min1) && (value <= field.max1) {
            return true
        }
        if (value >= field.min2) && (value <= field.max2) {
            return true
        }
    }
    return false
}

func solution(lines []string) RESULT_TYPE {
    i := 0

    fields := make([]Field, 0, len(lines))//Definite overprovisioning, but eh.
    for lines[i] != "" {
        fields = append(fields, parseField(lines[i]))
        i++
    }

    i++
    assertInput(lines[i], "your ticket:")
    i++

    //myTicket := parseTicket(lines[i])

    i++
    assertInput(lines[i], "")
    i++
    assertInput(lines[i], "nearby tickets:")
    i++

    nearbyTickets := make([][]int, 0, len(lines) - i)
    for i < len(lines) {
        nearbyTickets = append(nearbyTickets, parseTicket(lines[i]))
        i++
    }
    errorRate := 0
    for _, ticket := range nearbyTickets {
        for _, fieldValue := range ticket {
            if !isValid(fieldValue, fields) {
                errorRate += fieldValue
            }
        }
    }

    return errorRate;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        71,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 16
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
