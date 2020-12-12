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

const (
    NORTH   = 'N'
    SOUTH   = 'S'
    EAST    = 'E'
    WEST    = 'W'
    LEFT    = 'L'
    RIGHT   = 'R'
    FORWARD = 'F'
)

type Heading struct {
    x int
    y int
}

var HEADING_EAST = Heading{1, 0}
var HEADING_WEST = Heading{-1, 0}
var HEADING_NORTH = Heading{0, 1,}
var HEADING_SOUTH = Heading{0, -1}

type Ship struct {
    x int
    y int
    facing [4]Heading
}

func actualMod(n int, m int) int {
    toReturn := n % m
    if toReturn < 0 {
        return toReturn + m
    }
    return toReturn
}
func rot(arr [4]Heading, dist int) [4]Heading {
    toReturn := [4]Heading{}
    for i, h := range arr {
        toReturn[actualMod(i - dist, 4)] = h
    }
    return toReturn
}
func parse (line string) func(*Ship) {
    instruction := line[0]
    magnitude, err := strconv.Atoi(line[1:])
    checkErr(err)
    if (instruction == LEFT || instruction == RIGHT) && magnitude % 90 != 0 {
        panic("Turn on non-right angle!")
    }
    switch instruction {
        case NORTH: return func (s *Ship) {
            s.y += magnitude
        }
        case SOUTH: return func (s *Ship) {
            s.y -= magnitude
        }
        case EAST: return func (s *Ship) {
            s.x += magnitude
        }
        case WEST: return func (s *Ship) {
            s.x -= magnitude
        }
        case LEFT: return func (s *Ship) {
            s.facing = rot(s.facing, -1 * (magnitude / 90))
        }
        case RIGHT: return func (s *Ship) {
            s.facing = rot(s.facing, magnitude / 90)
        }
        case FORWARD: return func (s *Ship) {
            s.x += magnitude * s.facing[0].x
            s.y += magnitude * s.facing[0].y
        }
    }
    panic("unknown instruction: " + string(instruction))
}

//Seriously, Go!?
func abs(in int) int {
    if in < 0 {
        return -in
    }
    return in
}

func solution(lines []string) RESULT_TYPE {
    ship := &Ship{
        0,
        0,
        [4]Heading{HEADING_EAST, HEADING_SOUTH, HEADING_WEST, HEADING_NORTH},
    }

    for _, line := range lines {
        parse(line)(ship)
    }
    return abs(ship.x) + abs(ship.y)
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        25,
        0,
        0,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 12
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
