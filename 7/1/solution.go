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

type BagCount struct {
    num int
    color string
}

func parseContents(contents string) []*BagCount {
    if (contents == "no other bags.") {
        return make([]*BagCount, 0)
    }
    parts := strings.Split(contents, ", ")
    toReturn := make([]*BagCount, len(parts))
    for i, part := range parts {
        var count int
        var adj, color string
        _, err := fmt.Sscanf(part, "%d %s %s bag", &count, &adj, &color)
        //fmt.Printf("%v, %d, %s, %s, 
        checkErr(err);
        toReturn[i] = &BagCount {
            count,
            adj + " " + color,
        }
    }
    return toReturn
}

func canContain(startingColor string, targetColor string, rules map[string][]*BagCount) bool {
    ruleSet, ok := rules[startingColor]
    if !ok {
        panic("Cannot find bag contents: " + startingColor)
    }

    for _, rule := range ruleSet {
        if rule.color == targetColor {
            return true
        }
    }
    for _, rule := range ruleSet {
        if canContain(rule.color, targetColor, rules) {
            return true
        }
    }

    return false
}

func parseRule(line string) (string, []*BagCount) {
    var adj, color string
    _, err := fmt.Sscanf(line, "%s %s bags contain", &adj, &color)
    checkErr(err);
    stemLength := len(adj) + len(color) + 15
    return adj + " " + color, parseContents(line[stemLength:])
}

func solution(lines []string) RESULT_TYPE {
    grammar := make(map[string][]*BagCount)
    for _, line := range lines {
        color, contents := parseRule(line)
        grammar[color] = contents
    }

    count := 0
    for startingColor := range grammar {
        if canContain(startingColor, "shiny gold", grammar) {
            count ++
        }
    }
    return count;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        4,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 7
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
