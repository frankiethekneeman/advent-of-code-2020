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


func solution(lines []string) RESULT_TYPE {
    black := make(map[Tile] bool)
    for _, line := range lines {
        current := Tile{}
        for _, step := range parse(line) {
            current = current.update(BASIS_VECTORS[step])
        }
        if black[current] {
            delete(black, current)
        } else {
            black[current] = true
        }
    }

    return len(black);
}

var DIRECTION_REGEX = regexp.MustCompile("(e|w|ne|nw|se|sw)")
func parse(line string) []string {
    return DIRECTION_REGEX.FindAllString(line, -1)
}
type Tile struct {x,y,z int}

var BASIS_VECTORS = map[string]Tile {
    "w": Tile{-1, 1, 0},
    "e": Tile{1, -1, 0},
    "se": Tile{0, -1, +1},
    "sw": Tile{-1, 0, +1},
    "ne": Tile{1, 0, -1},
    "nw": Tile{0, 1, -1},
}

func (t Tile) update (o Tile) Tile {
    return Tile{
        t.x + o.x,
        t.y + o.y,
        t.z + o.z,
    }
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        10,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 24
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
