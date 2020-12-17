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


func isValid(value int, fields map[string]Field) bool {
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

func isValidTicket(ticket []int, fields map[string]Field) bool {
    for _, fieldValue := range ticket {
        if !isValid(fieldValue, fields) {
            return false
        }
    }
    return true
}

func filterTickets(tickets [][]int, fields map[string]Field) [][]int {
    toReturn := make([][]int, 0, len(tickets))
    for _, ticket := range tickets {
        if isValidTicket(ticket, fields) {
            toReturn = append(toReturn, ticket)
        }
    }
    return toReturn
}

func getPossibleFields(tickets [][]int, fields map[string]Field, position int) []string {
    remainingPossible := make(map[string] bool)
    for name := range fields {
        remainingPossible[name] = true
    }
    for _, ticket := range tickets {
        for name := range remainingPossible {
            value := ticket[position]
            field := fields[name]
            if (value >= field.min1) && (value <= field.max1) {
                continue
            }
            if (value >= field.min2) && (value <= field.max2) {
                continue
            }
            delete(remainingPossible, name)
        }
    }
    toReturn := make([]string, len(remainingPossible))
    i := 0
    for name := range remainingPossible {
        toReturn[i] = name
        i++
    }
    return toReturn
}
func deepLen(arr [][]string) int {
    toReturn := 0
    for _, subArr := range arr {
        toReturn += len(subArr)
    }
    return toReturn
}

func filterCandidates(candidates [][]string, known map[string]bool) [][]string {
    toReturn := make([][]string, len(candidates))
    for i, colMatches := range candidates {
        if len(colMatches) == 1 {
            toReturn[i] = colMatches
        } else {
            filtered := make([]string, 0, len(colMatches))
            for _, match := range colMatches {
                if !known[match] {
                    filtered = append(filtered, match)
                }
            }
            if len(filtered) == 0 {
                panic(fmt.Sprintf("Reduced column '%v' to zero matches (%v)", colMatches, candidates))
            }
            toReturn[i] = filtered
        }
    }
    return toReturn
}

func uniquify(candidates [][]string) []string {
    columns := len(candidates)
    matches := deepLen(candidates)
    for columns != matches {
        //Collect all the resolved ones
        known := make(map[string] bool)
        for _, possibilities := range candidates {
            if len(possibilities) == 1 {
                known[possibilities[0]] = true
            }
        }

        candidates = filterCandidates(candidates, known)
        nMatches := deepLen(candidates)
        if matches == nMatches {
            panic(fmt.Sprintf("Cannot reduce further: %v", candidates))
        }
        matches = nMatches
    }
    toReturn := make([]string, len(candidates))

    for i, result := range candidates {
        toReturn[i] = result[0]
    }
    return toReturn
}

func solution(lines []string) RESULT_TYPE {
    i := 0

    fields := make(map[string]Field)//Definite overprovisioning, but eh.
    for ; lines[i] != ""; i++ {
        field := parseField(lines[i])
        fields[field.name] = field
    }

    i++
    assertInput(lines[i], "your ticket:")
    i++

    myTicket := parseTicket(lines[i])

    i++
    assertInput(lines[i], "")
    i++
    assertInput(lines[i], "nearby tickets:")
    i++

    nearbyTickets := make([][]int, 0, len(lines) - i)
    for ;i < len(lines); i++ {
        nearbyTickets = append(nearbyTickets, parseTicket(lines[i]))
    }
    validTickets := append(filterTickets(nearbyTickets, fields), myTicket)

    candidates := make([][]string, len(myTicket))
    for i := range myTicket {
        candidates[i] = getPossibleFields(validTickets, fields, i)
    }

    assignments := uniquify(candidates)
    //fmt.Printf("%v\n", assignments)

    result := 1
    for i, name := range assignments {
        if strings.HasPrefix(name, "departure ") {
            result *= myTicket[i]
        }
    }


    return result;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        1,
        1,
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
