package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type RESULT_TYPE = int64

/*
Begin Solution
*/

type Mask struct {
    ones int64
    zeros int64
}

type Computer struct {
    mask Mask
    memory map[int64]int64
}

func getOneMask(maskSpec string) int64 {
    result, err := strconv.ParseInt(
        strings.Replace(maskSpec, "X", "0", -1),
        2,
        64,
    )
    checkErr(err)
    return result
}

func getZeroMask(maskSpec string) int64 {
    result, err := strconv.ParseInt(
        strings.Replace(maskSpec, "X", "1", -1),
        2,
        64,
    )
    checkErr(err)
    return result
}

func getMemoryAddress(lhs string) int64 {
    var address int64
    _, err := fmt.Sscanf(lhs, "mem[%d]", &address)
    checkErr(err)
    return address
}

func parse(line string) func(*Computer) *Computer {
    parts := strings.Split(line, " = ")
    if parts[0] == "mask" {
        oneMask := getOneMask(parts[1])
        zeroMask := getZeroMask(parts[1])
        return func(c *Computer) *Computer {
            c.mask = Mask{
                oneMask,
                zeroMask,
            }
            return c
        }
    } else {
        value, err := strconv.ParseInt(parts[1], 10, 64)
        checkErr(err)
        address := getMemoryAddress(parts[0])
        return func(c *Computer) *Computer {
            c.memory[address] = (value | c.mask.ones) & c.mask.zeros
            return c
        }
    }
}

func solution(lines []string) RESULT_TYPE {
    computer := &Computer {
        Mask {
            0,
            (1 << 36) - 1,
        },
        make(map[int64] int64),
    }
    for _, line := range lines {
        computer = parse(line)(computer)
    }
    sum := int64(0)
    for _, value := range computer.memory {
        sum += value
    }
    return sum;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        51,
        165,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 14
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
