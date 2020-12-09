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

type Instruction struct {
    ins string
    arg int
}

func parseInstruction(line string) Instruction {
    var ins string
    var arg int
    _, err := fmt.Sscanf(line, "%s %d", &ins, &arg)
    checkErr(err)
    return Instruction {
        ins,
        arg,
    }
}

func copyProgram(program []Instruction, replaceIndex int, replacement Instruction ) []Instruction {
    toReturn := make([]Instruction, len(program))
    for i, curr := range program {
        if i == replaceIndex {
            toReturn[i] = replacement
        } else {
            toReturn[i] = curr
        }
    }
    return toReturn
}

func repairProgram(program []Instruction) [][]Instruction {
    count := 0
    for _, instruction := range program {
        if instruction.ins == "nop" || instruction.ins == "jmp" {
            count++
        }
    }
    toReturn := make([][]Instruction, count)
    curr := 0
    for i := range toReturn {
        for program[curr].ins == "acc" {
            curr ++
        }
        toChange := program[curr]
        if toChange.ins == "nop" {
            toReturn[i] = copyProgram(program, curr, Instruction{
                "jmp",
                toChange.arg,
            })
        } else {
            toReturn[i] = copyProgram(program, curr, Instruction{
                "nop",
                toChange.arg,
            })
        }

        curr ++
    }
    return toReturn
}
func runProgram(program []Instruction) (int, bool) {
    executedInstructions := make([]bool, len(program))
    accumulator := 0
    insPointer := 0
    for insPointer >= 0 && insPointer < len(program) && !executedInstructions[insPointer] {
        executedInstructions[insPointer] = true
        action := program[insPointer]
        if action.ins == "nop" {
            insPointer++
        } else if action.ins == "acc" {
            insPointer++
            accumulator += action.arg
        } else if action.ins == "jmp" {
            insPointer += action.arg
        } else {
            panic("Unrecognized command: " + action.ins)
        }
    }

    return accumulator, insPointer == len(program);
}

func solution(lines []string) RESULT_TYPE {
    program := make([]Instruction, len(lines))
    for i, line := range lines {
        program[i] = parseInstruction(line)
    }
    for _, repaired := range repairProgram(program) {
        accumulated, terminated := runProgram(repaired)
        if (terminated) {
            return accumulated
        }
    }
    return -1;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        8,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 8
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
