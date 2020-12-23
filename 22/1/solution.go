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

func solution(lines []string) RESULT_TYPE {
    nCards := (len(lines) - 3) / 2
    player1 := parseDeck(lines[1:nCards + 1])
    player2 := parseDeck(lines[nCards + 3:])
    for true {
        if player1.empty() {
            return player2.score()
        }
        if player2.empty() {
            return player1.score()
        }
        card1 := player1.deal()
        card2 := player2.deal()
        if card1 > card2 {
            player1.accept(card1, card2)
        } else {
            player2.accept(card2, card1)
        }
    }
    return -1;
}

func parseDeck(lines []string) *Deck {
    deck := newDeck();
    for _, line := range lines {
        card, err := strconv.Atoi(line)
        checkErr(err)
        deck.accept(card)
    }
    return deck
}

type Deck struct {
    length int
    head *Card
    tail *Card
}

func newDeck() *Deck {
    return &Deck{0, nil, nil}
}

func (d *Deck) accept(values ...int) {
    for _, v := range values {
        card := newCard(v)
        if d.length == 0 {
            d.head = card
            d.tail = card
            d.length = 1
        } else {
            d.tail.next = card
            d.tail = card
            d.length ++
        }
    }
}

func (d *Deck) deal() int {
    if d.empty() {
        panic("Cannot deal from an empty deck")
    }
    toDeal := d.head
    d.head = d.head.next
    if d.head == nil {
        d.tail = nil
    }
    d.length --
    return toDeal.value
}

func (d* Deck) empty() bool {
    return d.length == 0
}

func (d* Deck) score() int {
    score := 0
    currCard := d.head
    currPoints := d.length
    for currCard != nil {
        score += currCard.value * currPoints
        currCard = currCard.next
        currPoints --
    }
    return score
}

type Card struct {
    value int
    next *Card
}

func newCard(value int) *Card {
    return &Card{value, nil}
}



/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        306,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 22
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
