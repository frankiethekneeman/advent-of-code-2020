package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type RESULT_TYPE = int

/*
Begin Solution
*/


type NonTerminal struct {
    identity int
    replacements [][]int
}

type Terminal struct {
    identity int
    literal rune
}

type Grammar = map[int] Rule
type Rule interface {
    id() int
    asRegularExpression(g Grammar) string
}

func (t Terminal) id() int {
    return t.identity
}
func (t Terminal) asRegularExpression(_ Grammar) string {
    return string(t.literal)
}

func (n NonTerminal) id() int {
    return n.identity
}

func (n NonTerminal) asRegularExpression(g Grammar) string {
    subExpressions := make([]string, len(n.replacements))
    for i, replacement := range n.replacements {
        var expressionBuilder strings.Builder
        for _, id := range replacement {
            r, ok := g[id]
            if !ok {
                panic("No rule for id " + strconv.Itoa(id))
            }
            expressionBuilder.WriteString(r.asRegularExpression(g))
        }
        subExpressions[i] = expressionBuilder.String()
    }
    if len(subExpressions) == 1 {
        return subExpressions[0]
    }
    return "(" + strings.Join(subExpressions, "|") + ")"
}

func parseRule(line string) Rule {
    parts := strings.Split(line, ": ")
    id, err := strconv.Atoi(parts[0])
    checkErr(err)
    if parts[1] == "\"a\"" || parts[1] == "\"b\"" {
        return Terminal {
            id,
            []rune(parts[1])[1],
        }
    }
    replacementStrings := strings.Split(parts[1], " | ")
    replacements := make([][]int, len(replacementStrings))
    for i, replacementString := range replacementStrings {
        ids := strings.Split(replacementString, " ")
        replacement := make([]int, len(ids))
        for j, id := range ids {
            parsedId, err := strconv.Atoi(id)
            checkErr(err)
            replacement[j] = parsedId
        }
        replacements[i] = replacement
    }

    return NonTerminal { id, replacements }

}

func solution(lines []string) RESULT_TYPE {
    grammar := Grammar{}
    i := 0
    for ; lines[i] != ""; i++ {
        rule := parseRule(lines[i])
        grammar[rule.id()] = rule
    }
    expr := regexp.MustCompile("^" + grammar[0].asRegularExpression(grammar) + "$")
    
    count := 0
    for i++; i < len(lines); i++ {
        if expr.MatchString(lines[i]) {
            count++
        }
    }


    return count;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        2,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 19
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
