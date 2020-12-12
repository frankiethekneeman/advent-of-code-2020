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

type Ship struct {
    x int
    y int
    waypoint Heading
}

func rot(h Heading, dir rune, times int) Heading {
    if times == 0 {
        return h
    }
    h = rot(h, dir, times - 1)
    switch dir {
        case RIGHT: return Heading{
            h.y,
            -h.x,
        }
        case LEFT: return Heading{
            -h.y,
            h.x,
        }
    }
    panic("Unrecognized rotation direction: " + string(dir))
}

func parse (line string) func(*Ship) {
    instruction := []rune(line)[0]
    magnitude, err := strconv.Atoi(string([]rune(line)[1:]))
    checkErr(err)
    if (instruction == LEFT || instruction == RIGHT) && magnitude % 90 != 0 {
        panic("Turn on non-right angle!")
    }
    switch instruction {
        case NORTH: return func (s *Ship) {
            s.waypoint.y += magnitude
        }
        case SOUTH: return func (s *Ship) {
            s.waypoint.y -= magnitude
        }
        case EAST: return func (s *Ship) {
            s.waypoint.x += magnitude
        }
        case WEST: return func (s *Ship) {
            s.waypoint.x -= magnitude
        }
        case LEFT: return func (s *Ship) {
            s.waypoint = rot(s.waypoint, instruction, magnitude / 90)
        }
        case RIGHT: return func (s *Ship) {
            s.waypoint = rot(s.waypoint, instruction, magnitude / 90)
        }
        case FORWARD: return func (s *Ship) {
            s.x += magnitude * s.waypoint.x
            s.y += magnitude * s.waypoint.y
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
        Heading{10,1},
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
        286,
        10,
        190,
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
