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

func solution(lines []string) RESULT_TYPE {
    recipes := make([]Recipe, len(lines))
    for i, line := range lines {
        recipes[i] = parseRecipe(line)
    }
    candidateSets := make(map[string]Set)
    allergenFreeIngredients := make(Set)
    for _, recipe := range recipes {
        allergenFreeIngredients.addAll(recipe.ingredients)
        for allergen := range recipe.allergens {
            candidateSet, ok := candidateSets[allergen]
            if (!ok) {
                candidateSet = duplicate(recipe.ingredients)
            }
            candidateSet.keepOnly(recipe.ingredients)
            candidateSets[allergen] = candidateSet
        }
    }
    for _, candidateSet := range candidateSets {
        allergenFreeIngredients.removeAll(candidateSet)
    }

    appearances := 0
    for _, recipe := range recipes {
        for ingredient := range recipe.ingredients {
            if allergenFreeIngredients.contains(ingredient) {
                appearances ++
            }
        }
    }
    return appearances;
}

type Recipe struct {
    ingredients Set
    allergens Set
}

func parseRecipe(line string) Recipe {
    parts := strings.Split(line, " (contains ")
    ingredients := make(Set)
    for _, ingredient := range strings.Split(parts[0], " ") {
        ingredients.add(ingredient)
    }
    allergens := make(Set)
    //Trim off that last closing paren with an ugly slice
    for _, allergen := range strings.Split(parts[1][:len(parts[1]) - 1], ", ") {
        allergens.add(allergen)
    }
    return Recipe{
        ingredients,
        allergens,
    }
}

type Set map[string]bool

func (s Set) addAll (other Set) {
    for item := range other {
        s[item] = true
    }
}

func (s Set) keepOnly (other Set) {
    for item := range s {
        if _, ok := other[item]; !ok {
            delete(s, item)
        }
    }
}

func (s Set) removeAll(other Set) {
    for item := range other {
        delete(s, item)
    }
}

func (s Set) contains(key string) bool {
    _, ok := s[key]
    return ok
}

func (s Set) add(key string) {
    s[key] = true
}

func duplicate(s Set) Set {
    toReturn := make(Set)
    for item := range s {
        toReturn[item] = true
    }
    return toReturn
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        5,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 21
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
