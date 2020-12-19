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
    consume(tokens []rune, start int, g Grammar, out chan<- int)
}

func (t Terminal) id() int {
    return t.identity
}

func (t Terminal) consume(tokens []rune, start int, _ Grammar, out chan<- int) {
    if start < len(tokens) && tokens[start] == t.literal {
        out <- start + 1
    }
    close(out)
}

func (n NonTerminal) id() int {
    return n.identity
}

func consume(replacement []int, tokens []rune, start int, g Grammar, out chan<- int) {
    if len(replacement) == 0 {
        out <- start
    } else {
        nexts := make(chan int, 1)
        rule, ok := g[replacement[0]]
        if !ok {
            panic("No rule for id " + strconv.Itoa(replacement[0]))
        }
        go rule.consume(tokens, start, g, nexts)
        for next := range nexts {
            ends := make(chan int, 1)
            go consume(replacement[1:], tokens, next, g, ends)
            for end := range ends {
                out <- end
            }
        }
    }
    close(out)
}

func (n NonTerminal) consume(tokens []rune, start int, g Grammar, out chan<- int) {
    for _, replacement := range n.replacements {
        nexts := make(chan int, 1)
        go consume(replacement, tokens, start, g, nexts)
        for next := range nexts {
            out <- next
        }
    }
    close(out)
}

func parseRule(line string) Rule {
    parts := strings.Split(line, ": ")
    id, err := strconv.Atoi(parts[0])
    checkErr(err)
    if id == 8 {
        return NonTerminal {
            8,
            [][]int {
                []int { 42 },
                []int { 42, 8 },
            },
        }
    }
    if id == 11 {
        return NonTerminal {
            11,
            [][]int {
                []int {42, 31},
                []int {42, 11, 31},
            },
        }
    }
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

func matches(g Grammar, input string) bool {
    endpoints := make(chan int, 1)
    tokens := []rune(input)
    go g[0].consume(tokens, 0, g, endpoints)
    for endpoint := range endpoints {
        if endpoint == len(tokens) {
            return true
        }
    }
    return false
}

func solution(lines []string) RESULT_TYPE {
    grammar := Grammar{}
    i := 0
    for ; lines[i] != ""; i++ {
        rule := parseRule(lines[i])
        grammar[rule.id()] = rule
    }
    count := 0
    for i++; i < len(lines); i++ {
        if matches(grammar, lines[i]) {
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
        12,
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
