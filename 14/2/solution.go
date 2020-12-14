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
    forceOnes int64
    elimXs int64
    floatingBits []int64
}

type Computer struct {
    mask Mask
    memory map[int64]int64
}

func getXLocations(maskSpec string) []uint {
    toReturn := make([]uint, 0, len(maskSpec))
    for i, char := range []rune(maskSpec) {
        if char == 'X' {
            // length minus index is always >= 1, so uint cast is safe.
            location := uint(len(maskSpec) - i - 1)
            toReturn = append(toReturn, location)
        }
    }
    return toReturn
}

func generateFloatingBits(maskSpec string) []int64 {
    xLocations := getXLocations(maskSpec)
    if len(xLocations) == 0 {
        return []int64{0}
    }
    l := uint(len(xLocations))
    toReturn := make([]int64, 1 << l)
    for i, _ := range toReturn {
        floatingBit := int64(0)
        for seedPos, finalPos := range xLocations {
            seedShift := uint(seedPos) //Indices are definitionally unsigned
            seed := (int64(i) & (int64(1) << seedShift)) >> seedShift
            final := seed << finalPos
            floatingBit += final
        }
        toReturn[i] = floatingBit
    }

    return toReturn
}

func parseMask(maskSpec string) Mask {
    forceOnes, err := strconv.ParseInt(
        strings.Replace(maskSpec, "X", "0", -1),
        2,
        64,
    )
    checkErr(err)
    elimXs, err := strconv.ParseInt(
        strings.Replace(strings.Replace(maskSpec, "0", "1", -1), "X", "0", -1),
        2,
        64,
    )
    checkErr(err)
    return Mask {
        forceOnes,
        elimXs,
        generateFloatingBits(maskSpec),
    }

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
        mask := parseMask(parts[1])
        return func(c *Computer) *Computer {
            c.mask = mask
            return c
        }
    } else {
        value, err := strconv.ParseInt(parts[1], 10, 64)
        checkErr(err)
        address := getMemoryAddress(parts[0])
        return func(c *Computer) *Computer {
            for _, floatingBits := range c.mask.floatingBits {
                resolvedAddress := ((address | c.mask.forceOnes) & c.mask.elimXs) | floatingBits
                c.memory[resolvedAddress] = value
            }
            return c
        }
    }
}

func solution(lines []string) RESULT_TYPE {
    computer := &Computer {
        Mask {
            0,
            (1 << 36) - 1,
            []int64{0},
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
        208,
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
