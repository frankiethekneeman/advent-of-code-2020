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

const (
    ON = '#'
    OFF = '.'
)

type Edge int
const (
    TOP Edge = iota
    RIGHT
    BOTTOM
    LEFT
)

var EDGES = []Edge{TOP, RIGHT, BOTTOM, LEFT}
func (e Edge) String() string {
    return [...]string{"TOP", "RIGHT", "BOTTOM", "LEFT"}[e]
}


type Tile struct {
    id int
    data [10][10] bool
}

//Each side is read numerically as if it had been rotated to the top.
func (t Tile) top(reversed bool) int {
    bits := t.data[0]
    if reversed {
        bits = reverse(bits)
    }
    return asBinaryNumber(bits)
}

func (t Tile) right(reversed bool) int {
    bits := getCol(t.data, 9)
    if reversed {
        bits = reverse(bits)
    }
    return asBinaryNumber(bits)
}

func (t Tile) bottom(reversed bool) int {
    bits := t.data[9]
    if !reversed { //Rotating the bottom to the top would reverse it.
        bits = reverse(bits)
    }
    return asBinaryNumber(bits)
}

func (t Tile) left(reversed bool) int {
    bits := getCol(t.data, 0)
    if !reversed { //Normally read bottom up, instead of top down
        bits = reverse(bits)
    }
    return asBinaryNumber(bits)
}

func (t Tile) edge(e Edge, reversed bool) int {
    switch e {
        case TOP: return t.top(reversed)
        case LEFT: return t.left(reversed)
        case RIGHT: return t.right(reversed)
        case BOTTOM: return t.bottom(reversed)
    }
    panic("Unrecognized edge: " + string(e))
}

type PlacedTile struct {
    tile Tile
    topEdge Edge //Edge in original rotated to the top
    flipped bool
}
func trueMod(n int, m int) int {
    if (n % m) < 0 {
        return (n % m) + m
    }
    return n % m
}
func (p PlacedTile) getBit(row int, column int) bool {
    switch p.topEdge {
        case TOP:
        case BOTTOM:
            row = 9 - row
            column = 9 - column
        case LEFT:
            if p.flipped {
                //RotR - to undo the RotL
                row, column = column, 9 - row
            } else {
                //RotL - to undo the RotR
                row, column = 9 - column, row
            }
        case RIGHT:
            if p.flipped {
                //RotL - to undo the RotR
                row, column = 9 - column, row
            } else {
                //RotR - to undo the RotL
                row, column = column, 9 - row
            }
        default:
            panic("Unrecognized topEdge")
    }
    if p.flipped {
        column = 9 - column
    }
    return p.tile.data[row][column]
}
func (p PlacedTile) edge(e Edge) int {
    return p.tile.edge(
        edgeFromRotatedTop(p.topEdge, e, p.flipped),
        p.flipped,
    )
}
func (p PlacedTile) top() int {
    return p.edge(TOP)
}
func (p PlacedTile) right() int {
    return p.edge(RIGHT)
}
func (p PlacedTile) left() int {
    return p.edge(LEFT)
}
func(p PlacedTile) bottom() int {
    return p.edge(BOTTOM)
}


type Picture = [][]PlacedTile

func solution(lines []string) RESULT_TYPE {
    validateEdgeLogic()
    validateTranspose()
    if len(lines) % 12 != 0 {
        panic("Unparseable Input")
    }

    nTiles := len(lines) / 12

    tiles := make(map[int]Tile)
    edgeMap := make(map[int] []int)

    for i := 0; i * 12 < len(lines); i ++ {
        tile := parseTile(lines[i * 12 : (i + 1) * 12])
        if (tile.id == 0) {
            panic(fmt.Sprintf("%dth tile gives %v", i, tile))
        }
        tiles[tile.id] = tile
        for _, edge := range edges(tile) {
            slice := edgeMap[edge]
            if slice == nil {
                slice = make([]int, 0)
            }
            edgeMap[edge] = append(slice, tile.id)
        }
    }
    
    count := 0
    for edge, tiles := range edgeMap {
        if (edge == 0) {
            panic("that seems unlikely")
        }
        if len(tiles) > 2 {
            panic(fmt.Sprintf("%v", tiles))
        }
        count += len(tiles)
    }
    if count != nTiles * 8 {
        panic("malformed Edgemap")
    }

    reconstituted := align(tiles, edgeMap)

    image := make([][]bool, len(reconstituted) * 8)
    for tileRow, tiles := range reconstituted {
        for tileCol, tile := range tiles {
            for row := 0; row < 8; row++ {
                targetRow := 8 * tileRow + row
                if image[targetRow] == nil {
                    image[targetRow] = make([]bool, len(tiles) * 8)
                }
                for col := 0; col < 8; col++ {
                    targetCol := 8 * tileCol + col
                    image[targetRow][targetCol] = tile.getBit(row + 1, col + 1)
                }
            }
        }
    }

    //Doing an extra rotation to keep `image` sacrosanct.
    for rotations := 1; rotations <=4; rotations++ {
        candidate := rotateN(image, rotations)
        if n := countSeaMonsters(candidate); n > 0 {
            return countWaves(candidate) - 15 * n
        }
        antiCandidate := rotateN(flipRows(image), rotations)
        if n := countSeaMonsters(antiCandidate); n > 0 {
            return countWaves(antiCandidate) - 15 * n
        }
    }

    return int(-1)
}

func countWaves(img [][]bool) int {
    count := 0
    for _, row := range img {
        for _, bit := range row {
            if bit {
                count++
            }
        }
    }
    return int(count)
}

var SEA_MONSTER = [][]rune {
	[]rune("                  # "),
	[]rune("#    ##    ##    ###"),
	[]rune(" #  #  #  #  #  #   "),
}

func countSeaMonsters(img [][]bool) int {
    found := 0
	for r := 0; r < len(img) - len(SEA_MONSTER); r++ {
        for c := 0; c < len(img) - len(SEA_MONSTER[0]); c++ {
            if detectSeaMonster(img, r, c) {
                found++
            }
        }
    }
    return found
}

func detectSeaMonster(img [][]bool, rowOffset int, columnOffset int) bool {
    for r, row := range SEA_MONSTER {
        for c, monsterPart := range row {
            if monsterPart == ' ' {
                continue
            }
            if !img[rowOffset + r][columnOffset + c] {
                return false
            }
        }
    }
    return true
}

func flipRows(in [][]bool) [][]bool {
    dim := len(in)
    out := make([][]bool, dim)
    for r, row := range in {
        out[r] = make([]bool, dim)
        for c, bit := range row {
            out[r][(dim - 1) - c] = bit
        }
    }
    return out
}

func rotateN(in [][]bool, times int) [][]bool {
    if times == 0 {
        return in
    }
    return rotateN(rotate(in), times - 1)
}

func rotate(in [][]bool) [][]bool {
    dim := len(in)
    out := make([][]bool, dim)
    for i := range in {
        out[i] = make([]bool, dim)
    }
    for r, row := range in {
        for c, bit := range row {
            out[c][(dim-1) - r] = bit
        }
    }
    return out
}

func align(tiles map[int]Tile, edgeMap map[int][]int) Picture {
    corners := findCorners(tiles, edgeMap)

    imageDim := sqrt(len(tiles))

    image := make(Picture, imageDim)
    for i := range image {
        image[i] = make([]PlacedTile, imageDim)
    }

    image[0][0] = rotateUntil(tiles[corners[0]], func(p PlacedTile) bool {
        return len(edgeMap[p.top()]) == 1 && len(edgeMap[p.left()]) == 1
    })
    
    for i := range image[1:] {
        leftNeighbor := image[0][i]
        left := transpose(leftNeighbor.right())
        matches := edgeMap[left]
        if len(matches) != 2 {
            panic(fmt.Sprintf("position 0, %d. left pattern %d matches: %v.", i + 1, left, matches))
        }
        toPlace := matches[0]
        if toPlace == leftNeighbor.tile.id {
            toPlace = matches[1]
        }
        image[0][i + 1] = rotateUntil(tiles[toPlace], func(p PlacedTile) bool {
            return p.left() == left
        })
    }
    
    for i := range image[1:] {
        topNeighbor := image[i][0]
        top := transpose(topNeighbor.bottom())
        matches := edgeMap[top]
        if len(matches) != 2 {
            panic(fmt.Sprintf("position %d, 0. top pattern %d matches: %v.", i + 1, top, matches))
        }
        toPlace := matches[0]
        if toPlace == topNeighbor.tile.id {
            toPlace = matches[1]
        }
        image[i + 1][0] = rotateUntil(tiles[toPlace], func(p PlacedTile) bool {
            return p.top() == top
        })
    }

    for i := range image[1:] {
        for j := range image[1:] {
            leftNeighbor := image[i+1][j]
            topNeighbor := image[i][j+1]
            left := transpose(leftNeighbor.right())
            top := transpose(topNeighbor.bottom())
            matches := edgeMap[left]
            if len(matches) != 2 {
                panic(fmt.Sprintf("position %d, %d. left pattern %d, %v.", i + 1, j + 1, left, matches))
            }
            toPlace := matches[0]
            if toPlace == leftNeighbor.tile.id {
                toPlace = matches[1]
            }
            image[i + 1][j + 1] = rotateUntil(tiles[toPlace], func(p PlacedTile) bool {
                return p.top() == top && p.left() == left
            })
        }
    }
    return image
}

func transpose(in int) int {
    if (in > 1023) {
        panic(string(in) + " is not a valid 10 bit integer")
    }
    out := 0;
    for i := 0; i < 10; i++ {
        if (in & (1 << i)) != 0 {
            out += (512 >> i)
        }
    }
    return out
}

func assertEqual(actual int, expected int) {
    if actual != expected {
        panic(fmt.Sprintf("Expected %d, but got %d", expected, actual))
    }
}

func validateTranspose() {
    assertEqual(transpose(1), 512)
    assertEqual(transpose(512), 1)
    assertEqual(transpose(513), 513)
    for i := 0; i < 10; i++ {
        assertEqual(transpose(1 << i), 512 >> i)
    }
}

func rotateUntil(tile Tile, predicate func(PlacedTile) bool) PlacedTile {
    for _, edge := range EDGES {
        placement := PlacedTile {tile, edge, false}
        if predicate(placement) {
            return placement
        }
        antiPlacement := PlacedTile{tile, edge, true}
        if predicate(antiPlacement) {
            return antiPlacement
        }
    }
    panic("Request not achievable")
}

func findCorners(tiles map[int]Tile, edgeMap map[int] []int) [4]int {
    var corners [4]int
    nCorner := 0
    for _, tile := range tiles {
        unmatchedEdges := 0
        for _, edge := range edges(tile) {
            if len(edgeMap[edge]) == 1 {
                unmatchedEdges++
            }
        }
        if unmatchedEdges == 4 {
            corners[nCorner] = tile.id
            nCorner ++
        }
    }

    if nCorner != 4 {
        fmt.Printf("%v\n", corners)
        panic("Corners not found")
    }
    return corners
}

func edgeFromRotatedTop(topEdge Edge, relativeEdge Edge, flipped bool) Edge {
    var mult = 1
    if flipped {
        mult = -1
    }
    return Edge(trueMod(int(topEdge) + mult * int(relativeEdge), 4))
}

func validateEdgeLogic() {
    cases := map[struct{
        top Edge
        flipped bool
    }] [4]Edge{
        {TOP, false}: [4]Edge{TOP, RIGHT, BOTTOM, LEFT},
        {RIGHT, false}: [4]Edge{RIGHT, BOTTOM, LEFT, TOP},
        {BOTTOM, false}: [4]Edge{BOTTOM, LEFT, TOP, RIGHT},
        {LEFT, false}: [4]Edge{LEFT, TOP, RIGHT, BOTTOM},
        {TOP, true}: [4]Edge{TOP, LEFT, BOTTOM, RIGHT},
        {RIGHT, true}: [4]Edge{RIGHT, TOP, LEFT, BOTTOM},
        {BOTTOM, true}: [4]Edge{BOTTOM, RIGHT, TOP, LEFT},
        {LEFT, true}: [4]Edge{LEFT, BOTTOM, RIGHT, TOP},
    }
    for transformation, sides := range cases {
        for i, edge := range EDGES {
            expected := sides[i]
            actual := edgeFromRotatedTop(transformation.top, edge, transformation.flipped)
            if (expected != actual) {
                panic(fmt.Sprintf(
                    "With edge %s rotated to the top, and flip?(%t), new %s edge should be %s - but calculated %s",
                    transformation.top,
                    transformation.flipped,
                    edge,
                    expected,
                    actual,
                ))
            }
        }
    }
    testTile := parseTile(strings.Split(`Tile 99:
........#.
..........
..........
..........
#.........
..........
..........
.........#
..........
...#......`, "\n"))

    idCases := map[struct{
        top Edge
        flipped bool
    }] [4]int {
        {TOP, false}: [4]int{2, 4, 8, 16},
        {RIGHT, false}: [4]int{4, 8, 16, 2},
        {BOTTOM, false}: [4]int{8, 16, 2, 4},
        {LEFT, false}: [4]int{16, 2, 4, 8},
        {TOP, true}: [4]int{256, 32, 64, 128},
        {RIGHT, true}: [4]int{128, 256, 32, 64},
        {BOTTOM, true}: [4]int{64, 128, 256, 32},
        {LEFT, true}: [4]int{32, 64, 128, 256},
    }

    for transformation, sides := range idCases {
        for i, edge := range EDGES {
            placed := PlacedTile { testTile, transformation.top, transformation.flipped}
            expected := sides[i]
            actual := placed.edge(edge)
            if (expected != actual) {
                panic(fmt.Sprintf(
                    "With edge %s rotated to the top, and flip?(%t), new %s edge should be %d - but calculated %d",
                    transformation.top,
                    transformation.flipped,
                    edge,
                    expected,
                    actual,
                ))
            }
        }
    }
    
    rotationCases := map[struct{
        top Edge
        flipped bool
    }] [4][2]int {
        {TOP, false}: [4][2]int{
            [2]int {0, 8},
            [2]int {7, 9},
            [2]int {9, 3},
            [2]int {4, 0},
        },
        {RIGHT, false}: [4][2]int{
            [2]int {1, 0},
            [2]int {0, 7},
            [2]int {6, 9},
            [2]int {9, 4},
        },
        {BOTTOM, false}: [4][2]int{
            [2]int {9, 1},
            [2]int {2, 0},
            [2]int {0, 6},
            [2]int {5, 9},
        },
        {LEFT, false}: [4][2]int{
            [2]int {8, 9},
            [2]int {9, 2},
            [2]int {3, 0},
            [2]int {0, 5},
        },
        {TOP, true}: [4][2]int{
            [2]int {0, 1},
            [2]int {7, 0},
            [2]int {9, 6},
            [2]int {4, 9},
        },
        {RIGHT, true}: [4][2]int{
            [2]int {1, 9},
            [2]int {0, 2},
            [2]int {6, 0},
            [2]int {9, 5},
        },
        {BOTTOM, true}: [4][2]int{
            [2]int {9, 8},
            [2]int {2, 9},
            [2]int {0, 3},
            [2]int {5, 0},
        },
        {LEFT, true}: [4][2]int{
            [2]int {8, 0},
            [2]int {9, 7},
            [2]int {3, 9},
            [2]int {0, 4},
        },
    }
    for transformation, ons := range rotationCases {
        placed := PlacedTile { testTile, transformation.top, transformation.flipped}
        for _, on := range ons {
            col, row := on[1], on[0] // Why did I type these backwards?
            actual := placed.getBit(row, col)
            if (!actual) {
                panic(fmt.Sprintf(
                    "With edge %s rotated to the top, and flip?(%t), row %d, column %d should be on.", 
                    transformation.top,
                    transformation.flipped,
                    row,
                    col,
                ))
            }
            if row == 0 || row == 9 {
                for offCol := 0; offCol <10; offCol++ {
                    if (offCol == col) {
                        continue
                    }
                    if placed.getBit(row, offCol) {
                        panic(fmt.Sprintf("(%s, %t) %d, %d should be off, but it's not...",
                            transformation.top,
                            transformation.flipped,
                            row,
                            offCol,
                        ))
                    }
                }
            }
            if col == 0 || col == 9 {
                for offRow := 0; offRow <10; offRow++ {
                    if (offRow == row) {
                        continue
                    }
                    if placed.getBit(offRow, col) {
                        panic(fmt.Sprintf("(%s, %t) %d, %d should be off, but it's not...",
                            transformation.top,
                            transformation.flipped,
                            offRow,
                            col,
                        ))
                    }
                }
            }
        }
    }
}


func parseTile(lines []string) Tile {
    var id int

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

func edges(tile Tile) [8]int {
    var toReturn [8]int
    for i, edge := range EDGES {
        toReturn[2 * i] = tile.edge(edge, false)
        toReturn[2 * i + 1] = tile.edge(edge, true)
    }
    return toReturn
}

//This is inefficient, but I didn't feel like detecting if casting to float was
//losing precision
func sqrt(n int) int {
    for i := 1; i <= n; i++ {
        if n % i == 0 && n / i == i {
            return i
        }
        if n / i < i {
            panic(strconv.Itoa(n) + " is not a perfect square")
        }
    }
    panic("Cannot find sqrt of n <= 0")
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
      273,
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
