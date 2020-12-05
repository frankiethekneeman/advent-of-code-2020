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

type validationFunction func(string) bool
var REQUIRED_FIELDS = map[string]validationFunction {
    "byr": func(val string) bool {
        if len(val) != 4 {
            return false
        }
        birthYear, err := strconv.Atoi(val)
        if err != nil || birthYear < 1920 || birthYear > 2002 {
            return false
        }

        return true
    },
    "iyr": func(val string) bool {
        if len(val) != 4 {
            return false
        }
        issueYear, err := strconv.Atoi(val)
        if err != nil || issueYear < 2010 || issueYear > 2020 {
            return false
        }

        return true
    },
    "eyr": func(val string) bool {
        if len(val) != 4 {
            return false
        }
        expirationYear, err := strconv.Atoi(val)
        if err != nil || expirationYear < 2020 || expirationYear > 2030 {
            return false
        }

        return true
    },
    "hgt": func(val string) bool {
        var height int
        var units string
        _, err := fmt.Sscanf(val, "%d%s", &height, &units)
        if err != nil {
            return false
        }
        if units == "cm" && height >=150 && height <=193 {
            return true
        }
        if units == "in" && height >=59 && height <=76 {
            return true
        }
        return false
    },
    "hcl": func(val string) bool {
        if len(val) != 7 {
            return false
        }
        var parsedColor int
        _, err := fmt.Sscanf(val, "#%x", &parsedColor)
        if err != nil {
            return false
        }
        return val == fmt.Sprintf("#%06x", parsedColor)
    },
    "ecl": func(color string) bool {
        return color == "amb" ||
            color == "blu" ||
            color == "brn" ||
            color == "gry" ||
            color == "grn" ||
            color == "hzl" ||
            color == "oth"
    },
    "pid": func(val string) bool {
        if len(val) != 9 {
            return false
        }
        id, err := strconv.Atoi(val)
        if err != nil {
          return false
        }
        return val == fmt.Sprintf("%09d", id)
    },
}

func isValid(passport map[string]string) bool {
    for field, validator := range REQUIRED_FIELDS {
        val, exists := passport[field]
        if !exists || !validator(val) {
            return false
        }
    }
    return true
}

func addFields(passport map[string] string, line string) {
    for _, entry := range strings.Split(line, " ") {
        pieces := strings.Split(entry, ":")
        passport[pieces[0]] = pieces[1]
    }
}

func solution(lines []string) RESULT_TYPE {
    passport := make(map[string] string)
    valid := 0

    for _, line := range lines {
        if line == "" {
            if isValid(passport) {
                valid ++
            }
            passport = make(map[string] string)
        } else {
            addFields(passport, line)
        }
    }
    if isValid(passport) {
        valid ++
    }

    return valid;
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {

    return []RESULT_TYPE {
        1,
        0,
        1,
        0,
        2,
        0,
        0,
        0,
        0,
        1,
        1,
        1,
        1,
        4,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 4
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
