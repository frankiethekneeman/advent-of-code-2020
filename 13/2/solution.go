package main

import (
    "bufio"
    "fmt"
    "math/big"
    "os"
    "strconv"
    "strings"
)

type RESULT_TYPE = int64

/*
Begin Solution
*/
const (
    NO_CONSTRAINT = "x"
)

type TargetDef struct {
    value *big.Int
    modulo *big.Int
}

func parseBusses(line string) []TargetDef {
    parts := strings.Split(line, ",")
    toReturn := make([]TargetDef, 0, len(parts))
    for i, part := range parts {
        if part == NO_CONSTRAINT {
            continue
        }
        id, err := strconv.ParseInt(part, 10, 64)
        checkErr(err)
        bigI := big.NewInt(int64(i))
        bigId := big.NewInt(int64(id))
        var val big.Int
        val.Sub(bigId, bigI)
        val.Mod(&val, bigId)
        toReturn = append(toReturn, TargetDef{
            &val,
            bigId,
        })
    }
    return toReturn
}

func multInv(value int64, modulo int64) int64 {
    for i := int64(0); i < modulo; i++ {
        if value * i % modulo == 1 {
            return i
        }
    }
    panic("No multiplicative Inverse")
}
func combine(l TargetDef, r TargetDef) TargetDef {
    if l.modulo.Cmp(r.modulo) == -1 {
        return combine(r, l)
    }
    coFactor := big.NewInt(multInv(l.modulo.Int64(), r.modulo.Int64()))
    var newMod big.Int
    newMod.Mul(l.modulo, r.modulo)
    var newVal big.Int
    newVal.Sub(r.value, l.value)
    newVal.Mul(&newVal, l.modulo)
    newVal.Mul(&newVal, coFactor)
    newVal.Add(&newVal, l.value)
    newVal.Mod(&newVal, &newMod)
    return TargetDef{
        &newVal,
        &newMod,
    }
}

func solution(lines []string) RESULT_TYPE {
    busses := parseBusses(lines[1])
    target := busses[0]
    for _, bus := range busses[1:] {
        target = combine(target, bus)
    }

    return target.value.Int64()
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        1068781,
        3417,
        754018,
        779210,
        1261476,
        1202161486,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 13
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
