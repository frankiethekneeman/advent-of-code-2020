package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    input := getInput(1)
    TARGET := 2020
    lineItems := make(map[int]bool) // Poor Man's Set.
    for _, line := range input {
        asInt, err := strconv.Atoi(line)
        checkErr(err)
        lineItems[asInt] = true
    }

    for item := range lineItems {
        if lineItems[TARGET - item] {
            fmt.Printf("%d\n", item * (TARGET - item))
            return
        }
    }
    panic("No result found")
}

func getInput(day int) []string {
    inputLocation := strconv.Itoa(day) + "/input"
    file, err := os.Open(inputLocation)
    checkErr(err)
    scanner := bufio.NewScanner(file)
    var lines []string
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    checkErr(scanner.Err())
    return lines
}
