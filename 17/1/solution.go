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

type ConwayCube struct {
    x, y, z int
}

type PocketDimension = map[ConwayCube] bool

const (
    STEPS = 6
    ON = '#'
)

func parse(lines []string) PocketDimension {
    toReturn := make(PocketDimension)
    z := 0
    for y, line := range lines {
        for x, state := range []rune(line) {
            if state == ON {
                toReturn[ConwayCube{ x,y,z }] = true
            }
        }
    }
    return toReturn
}

func neighbors(cube ConwayCube) []ConwayCube {
    toReturn := make([]ConwayCube, 26)
    i := 0;
    for dx := -1; dx <= 1; dx ++ {
        for dy := -1; dy <= 1; dy ++ {
            for dz := -1; dz <= 1; dz ++ {
                if dx == 0 && dy == 0 && dz == 0 {
                    continue
                }
                toReturn[i] = ConwayCube {
                    cube.x + dx,
                    cube.y + dy,
                    cube.z + dz,
                }
                i++
            }
        }
    }
    return toReturn
}

func countActiveNeighbors(cube ConwayCube, prev PocketDimension) int {
    activeNeighbors := 0
    for _, neighbor := range neighbors(cube) {
        if prev[neighbor] {
            activeNeighbors++
        }
    }
    return activeNeighbors
}

func getCubesRemainingActive(prev PocketDimension) PocketDimension {
    next := make(PocketDimension)
    for cube := range prev {
        activeNeighbors := countActiveNeighbors(cube, prev)
        if (activeNeighbors == 2) || (activeNeighbors == 3) {
            next[cube] = true
        }
    }
    return next
}

func getActiveNeighborsOfInactiveCubes(prev PocketDimension) map[ConwayCube] int {
    toReturn := make(map[ConwayCube] int)
    for cube := range prev {
        for _, neighbor := range neighbors(cube) {
            if !prev[neighbor] {
                toReturn[neighbor] = toReturn[neighbor] + 1
            }
        }
    }
    return toReturn
}


func increment(prev PocketDimension) PocketDimension {
    toReturn := getCubesRemainingActive(prev)
    for cube, neighbors := range getActiveNeighborsOfInactiveCubes(prev) {
        if neighbors == 3 {
            toReturn[cube] = true
        }
    }

    return toReturn
}

func solution(lines []string) RESULT_TYPE {
    state := parse(lines)
    for i := 0; i < 6; i++ {
        state = increment(state)
    }
    return len(state);
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        112,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 17
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
