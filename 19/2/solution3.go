package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "sync"
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

//Not actually a generator - a complicated bidirectional lock that allows for generator-like
//behaviour.
type Generator struct {
    channel chan int
    readSignal chan bool
    mutex *sync.Mutex
    closed bool
}

func (c *Generator) close() {
    c.mutex.Lock()
    if (c.closed) {
        c.mutex.Unlock()
        return
    }
    c.closed = true
    close(c.readSignal)
    close(c.channel)
    c.mutex.Unlock()
}

func (c *Generator) read() (int, bool) {
    c.mutex.Lock()
    if c.closed {
        c.mutex.Unlock()
        return 0, false
    }
    c.readSignal <- true
    c.mutex.Unlock()
    toReturn, more := <-c.channel
    return toReturn, more
}

func (c *Generator) write(msg int) bool {
    c.mutex.Lock()
    if c.closed {
        c.mutex.Unlock()
        return false
    }
    c.channel <- msg
    c.mutex.Unlock()
    return <- c.readSignal
}

func generator() *Generator {
    return &Generator {
        make(chan int, 1),
        make(chan bool, 1),
        &sync.Mutex{},
        false,
    }
}

type Grammar = map[int] Rule
type Rule interface {
    id() int
    consume(tokens []rune, start int, g Grammar, out *Generator)
}

func (t Terminal) id() int {
    return t.identity
}

func (t Terminal) consume(tokens []rune, start int, _ Grammar, out *Generator) {
    if start < len(tokens) && tokens[start] == t.literal {
        if ! out.write(start + 1) {
            return
        }
    }
    out.close()
}

func (n NonTerminal) id() int {
    return n.identity
}

func consume(replacement []int, tokens []rune, start int, g Grammar, out *Generator) {
    if len(replacement) == 0 {
        if ! out.write(start) {
            return
        }
    } else {
        nexts := generator()
        rule, ok := g[replacement[0]]
        if !ok {
            panic("No rule for id " + strconv.Itoa(replacement[0]))
        }
        go rule.consume(tokens, start, g, nexts)
        for next, ok := nexts.read(); ok; next, ok = nexts.read() {
            ends := generator()
            go consume(replacement[1:], tokens, next, g, ends)
            for end, ok := ends.read(); ok; end, ok = ends.read() {
                if ! out.write(end) {
                    nexts.close()
                    ends.close()
                    return
                }
            }
        }
    }
    out.close()
}

func (n NonTerminal) consume(tokens []rune, start int, g Grammar, out *Generator) {
    for _, replacement := range n.replacements {
        nexts := generator()
        go consume(replacement, tokens, start, g, nexts)
        for next, ok := nexts.read(); ok; next, ok = nexts.read() {
            if ! out.write(next) {
                nexts.close()
                return
            }
        }
    }
    out.close()
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
    endpoints := generator()
    tokens := []rune(input)
    go g[0].consume(tokens, 0, g, endpoints)
    for endpoint, ok := endpoints.read(); ok; endpoint, ok = endpoints.read() {
        if endpoint == len(tokens) {
            endpoints.close()
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
