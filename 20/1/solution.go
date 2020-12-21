package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

type RESULT_TYPE = int64

/*
Begin Solution
*/

const (
    ON = '#'
    OFF = '.'
)

type Tile struct {
    id int64
    data [10][10] bool
}

func solution(lines []string) RESULT_TYPE {
    if len(lines) % 12 != 0 {
        panic("Unparseable Input")
    }

    nTiles := len(lines) / 12

    tiles := make([]Tile, nTiles)
    edgeMap := make(map[int]int)

    for i := 0; i * 12 < len(lines); i ++ {
        tiles[i] = parseTile(lines[i * 12 : (i + 1) * 12])
        //fmt.Printf("%v\n", tiles[i])
        for _, edge := range edges(tiles[i].data) {
            edgeMap[edge] = edgeMap[edge] + 1
        }
    }

    found := make([]int64,0,4)
    for _, tile := range tiles {
        unmatchedEdges := 0
        for _, edge := range edges(tile.data) {
            if edgeMap[edge] == 1 {
                unmatchedEdges++
            }
        }
        if unmatchedEdges == 4 {
            found = append(found, tile.id)
        }
    }
    if len(found) != 4 {
        fmt.Printf("%v\n", found)
        panic("Couldn't cheat my way to the corners.")
    }

    return found[0] * found[1] * found[2] * found[3]
}

func parseTile(lines []string) Tile {
    var id int64

    _, err := fmt.Sscanf(lines[0], "Tile %d:", &id)
    checkErr(err)

    var bits [10][10] bool

    for y, row := range lines[1:11] {
        for x, bit := range []rune(row) {
            if bit == ON {
                bits[y][x] = true
            } else if bit == OFF {
                bits[y][x] = false
            } else {
                panic("Unrecognized bit: " + string(bit))
            }
        }
    }

    return Tile{
        id,
        bits,
    }

}

func edges(tile [10][10] bool) [8]int {
    top := tile[0]
    bottom := tile[9]
    left := getCol(tile, 0)
    right := getCol(tile, 9)
    return [8]int {
        asBinaryNumber(top),
        asBinaryNumber(reverse(top)),
        asBinaryNumber(left),
        asBinaryNumber(reverse(left)),
        asBinaryNumber(bottom),
        asBinaryNumber(reverse(bottom)),
        asBinaryNumber(right),
        asBinaryNumber(reverse(right)),
    }
}

func getCol(tile[10][10] bool, col int) [10]bool {
    var toReturn [10]bool
    for i := range toReturn {
        toReturn[i] = tile[i][col]
    }
    return toReturn
}

func asBinaryNumber(bits [10]bool) int {
    toReturn := 0
    for i, bit := range bits {
        if bit {
            toReturn += 1 << (9 - i)
        }
    }
    return toReturn
}

func reverse(bits [10]bool) [10]bool {
    var toReturn [10]bool
    for i, bit := range(bits) {
        toReturn[9-i] = bit
    }
    return toReturn
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
      20899048083289,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 20
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
