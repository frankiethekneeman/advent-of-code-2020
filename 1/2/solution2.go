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
    RED = iota
    BLACK
)

type Ledger struct {
    color int
    value int
    count int
    largerTransactions *Ledger
    smallerTransactions *Ledger
    parent *Ledger
}

func (ledger *Ledger) Print() {
    if ledger == nil {
        return
    }
    ledger.smallerTransactions.Print()
    fmt.Printf("%d (%d) ", ledger.value, ledger.count)
    ledger.largerTransactions.Print()
}
func (ledger *Ledger) clone() *Ledger {
    if ledger == nil {
        return nil
    }
    toReturn := &Ledger{
        ledger.color,
        ledger.value,
        ledger.count,
        ledger.largerTransactions.clone(),
        ledger.smallerTransactions.clone(),
        nil,
    }
    if toReturn.largerTransactions != nil {
        toReturn.largerTransactions.parent = toReturn
    }
    if toReturn.smallerTransactions != nil {
        toReturn.smallerTransactions.parent = toReturn
    }
    return toReturn
}

func (ledger *Ledger) recursiveInsert(toAdd *Ledger ) *Ledger {
    if ledger.value == toAdd.value {
        ledger.count++
        return ledger;
    } else if ledger.value > toAdd.value {
        if ledger.smallerTransactions == nil {
            ledger.smallerTransactions = toAdd
            toAdd.parent = ledger
            return toAdd;
        }
        //insert
        return ledger.smallerTransactions.recursiveInsert(toAdd)
    } else {
        if ledger.largerTransactions == nil {
            ledger.largerTransactions = toAdd
            toAdd.parent = ledger
            return toAdd
        }
        return ledger.largerTransactions.recursiveInsert(toAdd)
    }
}

func (ledger *Ledger) uncle() *Ledger {
    if ledger == nil || ledger.parent == nil || ledger.parent.parent == nil {
        return nil
    }
    if ledger.parent == ledger.parent.parent.largerTransactions {
        return ledger.parent.parent.smallerTransactions
    }
    return ledger.parent.parent.largerTransactions
}

func (ledger *Ledger) rotateLeft() {
    if ledger == nil {
        return
    }
    toPromote := ledger.largerTransactions
    if toPromote == nil {
        panic("Left Rotate without right child")
    }
    fixed := ledger.parent
    ledger.largerTransactions = toPromote.smallerTransactions
    if ledger.largerTransactions != nil {
      ledger.largerTransactions.parent = ledger
    }
    toPromote.smallerTransactions = ledger
    ledger.parent = toPromote

    toPromote.parent = fixed;
    if fixed != nil {
        if fixed.smallerTransactions == ledger {
            fixed.smallerTransactions = toPromote
        } else {
            fixed.largerTransactions = toPromote
        }
    }
}

func (ledger *Ledger) rotateRight() {
    if ledger == nil {
        return
    }
    toPromote := ledger.smallerTransactions
    if toPromote == nil {
        panic("Right Rotate without Left child")
    }
    fixed := ledger.parent
    ledger.smallerTransactions = toPromote.largerTransactions
    if ledger.smallerTransactions != nil {
      ledger.smallerTransactions.parent = ledger
    }
    toPromote.largerTransactions = ledger
    ledger.parent = toPromote

    toPromote.parent = fixed;
    if fixed != nil {
        if fixed.largerTransactions == ledger {
            fixed.largerTransactions = toPromote
        } else {
            fixed.smallerTransactions = toPromote
        }
    }
}

func (ledger *Ledger) repair() {
    if ledger.parent == nil {
        //Sick, I'm the root node.  Nothing to do but paint this red door black.
        ledger.color = BLACK
        return
    }
    if ledger.parent.color == BLACK {
        //My parent is black, and both my children are black - I've extended no paths in
        //terms of black count, so nothing to do.
        return
    }
    //My parent is red.
    uncle := ledger.uncle()
    if uncle != nil && uncle.color == RED {
        // Paint my grand parent red, and its children black
        ledger.parent.parent.color = RED
        ledger.parent.color = BLACK
        uncle.color = BLACK
        //recurse
        ledger.parent.parent.repair()
        return
    }
    //... smaller === left, larger === right
    //Here comes the noise.  It's time for a rotation to keep that baby balanced.
    if ledger == ledger.parent.smallerTransactions && ledger.parent == ledger.parent.parent.largerTransactions {
        ledger.parent.rotateRight()
        //At this point, what was the new nodes grandparent is now its parent
        ledger.color = BLACK
        ledger.parent.color = RED
        ledger.parent.rotateLeft()
    } else if ledger == ledger.parent.largerTransactions && ledger.parent == ledger.parent.parent.smallerTransactions {
        ledger.parent.rotateLeft()
        //At this point, what was the new nodes grandparent is now its parent
        ledger.color = BLACK
        ledger.parent.color = RED
        ledger.parent.rotateRight()
    } else if (ledger == ledger.parent.smallerTransactions) {
        ledger.parent.color = BLACK
        ledger.parent.parent.color = RED
        ledger.parent.parent.rotateRight()
    } else {
        ledger.parent.color = BLACK
        ledger.parent.parent.color = RED
        ledger.parent.parent.rotateLeft()
    }

    
}

func (ledger *Ledger) insert(val int) *Ledger {
    if ledger == nil {
        return &Ledger{
            BLACK,
            val,
            1,
            nil,
            nil,
            nil,
        }
    }

    added := ledger.recursiveInsert(&Ledger {
        RED,
        val,
        1,
        nil,
        nil,
        nil,
    })
    if (added.count > 1) {
        //We already had this - no need to balance.
        return ledger
    }

    added.repair();
    root := added;
    for root.parent != nil {
        root = root.parent
    }
    return root
}

func (ledger *Ledger) rule1() {
    //Each node is red or black
    if ledger == nil {
        return //nil nodes are interpeted as black.
    }
    if ledger.color != BLACK && ledger.color != RED {
        panic("Node of unknown color?" + strconv.Itoa(ledger.color))
    }
    ledger.smallerTransactions.rule1()
    ledger.largerTransactions.rule1()
}

func (ledger *Ledger) rule4() {
    //Red nodes have ONLY black children
    if ledger == nil {
        return //nil nodes are interpeted as black.
    }
    if ledger.color == RED {
        if ledger.smallerTransactions != nil && ledger.smallerTransactions.color == RED {
            panic("Red node with red smaller child")
        }
        if ledger.largerTransactions != nil && ledger.largerTransactions.color == RED {
            panic("Red node with red smaller child")
        }
    }
    ledger.smallerTransactions.rule4()
    ledger.largerTransactions.rule4()
}
func (ledger *Ledger) blackDepth() int {
    if ledger == nil {
        return 1
    }
    subDepth := ledger.smallerTransactions.blackDepth() //doesn't matter which direction
    if ledger.color == BLACK {
        return subDepth + 1
    }
    return subDepth
}

func (ledger *Ledger) rule5() {
    //Any route to a leaf node should pas through the same number of black nodes.
    if ledger == nil {
        return // No children, no problem
    }
    if ledger.smallerTransactions.blackDepth() != ledger.largerTransactions.blackDepth() {
        panic("uneven black depth")
    }
    ledger.smallerTransactions.rule5()
    ledger.largerTransactions.rule5()
}

func (ledger *Ledger) validateRedBlackRules() {
    ledger.rule1()
    //Rule 2 - root is black
    if (ledger.color != BLACK) {
        panic("RED ROOT")
    }
    // Rule 3 is free - all leaves are nil
    ledger.rule4()
    ledger.rule5()
}

func (ledger *Ledger) getCount(val int) int {
    if ledger == nil {
        return 0
    }
    if ledger.value == val {
        return ledger.count
    }
    if ledger.value > val {
        return ledger.smallerTransactions.getCount(val)
    }
    return ledger.largerTransactions.getCount(val)
}

func (ledger *Ledger) produceLessThan(val int, exclusions *Ledger, c chan int) {
    if ledger == nil {
        return
    }
    if (val > ledger.value) {
        ledger.largerTransactions.produceLessThan(val, exclusions, c)
        if ledger.count - exclusions.getCount(ledger.value) > 0 {
            c <- ledger.value
        }
    }
    ledger.smallerTransactions.produceLessThan(val, exclusions, c)
}

func produceLessThan(ledger *Ledger, val int, exclusions *Ledger) chan int {
    c := make(chan int)
    go func() {
      ledger.produceLessThan(val, exclusions, c)
      close(c)
    }()
    return c
}

func findSumming(ledger *Ledger, sum int, limit int, alreadyUsed *Ledger) *Ledger {
    if limit == 1 {
        if ledger.getCount(sum) > alreadyUsed.getCount(sum) {
            return alreadyUsed.clone().insert(sum)
        }
        return nil
    }
    for candidate := range produceLessThan(ledger, sum, alreadyUsed) {
        solution := findSumming(ledger, sum - candidate, limit -1, alreadyUsed.clone().insert(candidate))
        if solution != nil {
          return solution
        }
    }
    return nil
}

func (ledger *Ledger) product() int {
    if ledger == nil {
        return 1
    }
    return ledger.smallerTransactions.product() * ledger.value * ledger.largerTransactions.product()
}
func solution(lines []string) RESULT_TYPE {
    TARGET := 2020
    var lineItems *Ledger
    for _, line := range lines {
        asInt, err := strconv.Atoi(line)
        checkErr(err)
        lineItems = lineItems.insert(asInt)
    }
    solution := findSumming(lineItems, TARGET, 3, nil)

    return solution.product();
}

/*
Test Cases
*/
func TEST_CASES() []RESULT_TYPE {
    return []RESULT_TYPE {
        241861950,
    }
}

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    DAY := 1
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
