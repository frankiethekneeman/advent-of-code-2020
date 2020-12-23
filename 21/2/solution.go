package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strconv"
    "strings"
)

type RESULT_TYPE = string

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

    for !allDistinct(candidateSets) {
        knownDangerous := getKnownDangerous(candidateSets)
        for _, candidates := range candidateSets {
            if len(candidates) == 1 {
                continue
            }
            candidates.removeAll(knownDangerous)
            //Maps are a pointer type, so I shouldn't need to re-assign
        }
    }

    allergens := getKeys(candidateSets)
    sort.Strings(allergens)
    dangerousIngredients := make([]string, len(allergens))

    for i, allergen := range allergens {
        dangerousIngredients[i] = candidateSets[allergen].getOnly()
    }

    return strings.Join(dangerousIngredients, ",")
}

func allDistinct(candidateSets map[string]Set) bool {
    for _, set := range candidateSets {
        if len(set) == 0 {
            panic("Overly Greedy elimination.")
        }
        if len(set) > 1 {
            return false
        }
    }
    return true
}

func getKnownDangerous(candidateSets map[string]Set) Set {
    toReturn := make(Set)
    for _, set := range candidateSets {
        if len(set) == 1 {
            toReturn.addAll(set)
        }
    }
    return toReturn
}

func getKeys(m map[string]Set) []string {
    toReturn := make([]string, len(m))
    i := 0
    for key := range m {
        toReturn[i] = key
        i++
    }
    return toReturn
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

func (s Set) getOnly() string {
    if len(s) != 1 {
        panic("Called getOnly on an unqualified Set.")
    }
    for k := range s { //WOW this is a weird construction.
        return k
    }
    panic("Static analysis isn't perfect - this code is unreachable")
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
        "mxmxvkd,sqjhc,fvjkl",
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
