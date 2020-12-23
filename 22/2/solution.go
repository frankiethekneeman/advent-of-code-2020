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
    nCards := (len(lines) - 3) / 2
    player1 := parseDeck(lines[1:nCards + 1])
    player2 := parseDeck(lines[nCards + 3:])
    if recursiveCombat(player1, player2) == PLAYER1 {
        return player1.score()
    } else {
        return player2.score()
    }
    return -1;
}

type Winner int
const (
    PLAYER1 = iota
    PLAYER2
)

func recursiveCombat(p1 *Deck, p2 *Deck) Winner {
    seen := make(map[string] bool)
    for true {
        state := getState(p1, p2)
        if seen[state] {
            return PLAYER1
        }
        seen[state] = true
        c1 := p1.deal()
        c2 := p2.deal()
        var roundWinner Winner
        if (c1 <= p1.length && c2 <= p2.length) {
            roundWinner = recursiveCombat(p1.duplicate(c1), p2.duplicate(c2))
        } else if c1 > c2 {
            roundWinner = PLAYER1
        } else {
            roundWinner = PLAYER2
        }
        if roundWinner == PLAYER1 {
            p1.accept(c1, c2)
        } else {
            p2.accept(c2, c1)
        }

        if p1.empty() || p2.empty() {
            return roundWinner
        }
    }
    return PLAYER1
}

func getState(p1 *Deck, p2 *Deck) string {
    var builder strings.Builder
    builder.Grow(p1.length * 3 + p2.length * 3)
    p1.addState(&builder)
    builder.WriteString("!")
    p2.addState(&builder)
    return builder.String()
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

func (d* Deck) duplicate(length int) *Deck {
    if length > d.length {
        panic("Duplicate called with too big a size")
    }
    n := newDeck()
    curr := d.head
    for n.length < length {
        n.accept(curr.value)
        curr = curr.next
    }
    return n
}

func (d* Deck) addState(builder *strings.Builder) {
    currCard := d.head
    for currCard != nil {
        builder.WriteString(strconv.Itoa(currCard.value))
        builder.WriteString(":")
        currCard = currCard.next
    }
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
      291,
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
