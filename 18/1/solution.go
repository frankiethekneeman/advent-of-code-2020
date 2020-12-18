package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
)

type RESULT_TYPE = int

/*
Begin Solution
*/

type Expr interface {
    evaluate() int
}

type BinaryExpr struct {
    lhs, rhs Expr
    operator rune
}

type Term struct {
    value int
}

func (e BinaryExpr) evaluate() int {
    l := e.lhs.evaluate()
    r := e.rhs.evaluate()
    if e.operator == '*' {
        return l * r
    }else if e.operator == '+' {
        return l + r
    }
    panic("Unrecognized operator: " + string(e.operator))
}

func (t Term) evaluate() int {
    return t.value
}

var TOKENIZER *regexp.Regexp = regexp.MustCompile(`\d+|[*+()]`)

func consumeTerm(tokens []string) (Term, []string) {
    value, err := strconv.Atoi(tokens[0])
    checkErr(err)
    return Term { value }, tokens[1:]
}

func consume(expected string, tokens []string) []string {
    if tokens [0] == expected {
        return tokens[1:]
    }
    panic("Expected " + expected + " but found "+ tokens[0])
}

func consumeLHS(tokens []string) (Expr, []string) {
    if tokens[0] == ")" {
        tokens = consume(")", tokens)
        // Deal with parentheses
        lhs, tokens := consumeExpr(tokens)
        tokens = consume("(", tokens)
        return lhs, tokens
    } else {
        lhs, rest := consumeTerm(tokens)
        return lhs, rest
    }
}

func consumeOperator(tokens []string) (rune, []string) {
    if tokens[0] == "*" || tokens[0] == "+" {
        return []rune(tokens[0])[0], tokens[1:]
    }
    panic("Expected operator but got " + string(tokens[0]))
}

func consumeExpr(tokens []string) (Expr, []string) {
    lhs, tokens := consumeLHS(tokens)
    if len(tokens) == 0 || tokens[0] == "(" {
        return lhs, tokens
    }

    operator, tokens := consumeOperator(tokens)
    
    rhs, tokens := consumeExpr(tokens)

    //fmt.Printf("%v %c %v\n", lhs, operator, rhs)

    return BinaryExpr{lhs, rhs, operator}, tokens
}
func parse(line string) Expr {
    tokens := TOKENIZER.FindAllString(line, -1)
    //fmt.Printf("%v\n", tokens)
    for i, j := 0, len(tokens) - 1; i < j; i, j = i + 1, j - 1 {
        tokens[i], tokens[j] = tokens[j], tokens[i]
    }
    //fmt.Printf("%v\n", tokens)
    toReturn, remaining := consumeExpr(tokens)
    if len(remaining) != 0 {
        panic("Failed Parse")
    }
    return toReturn
}


func solution(lines []string) RESULT_TYPE {
    sum := 0
    for _, line := range(lines) {
        sum += parse(line).evaluate()
    }
    return sum;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        71,
        51,
        26,
        437,
        12240,
        13632,
        26457,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 18
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
