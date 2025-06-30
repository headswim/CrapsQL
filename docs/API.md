# CrapsQL API Documentation

## Overview

CrapsQL is a domain-specific language for craps table simulation. This document provides complete API documentation for the Go package.

---

## Canonical Bet Logic

**All bet logic, payout calculations, and working status are defined in `pkg/crapsgame/canonical_bets.go`.**
- This file is the single source of truth for all bet types, categories, payout ratios, and working rules.
- All bet placement, validation, and resolution in the engine reference this canonical map.

---

## CanonicalBetDefinition Fields

Each bet type is defined by a `CanonicalBetDefinition` struct with the following fields:

```go
Name              string      // Human-readable name
Category          BetCategory // Bet category (see below)
Description       string      // Description of the bet
Payout            string      // Human-readable payout (e.g., "7:6")
WorkingBehavior   string      // "ALWAYS", "ONE_ROLL", "CONDITIONAL"
OneRoll           bool        // True if bet resolves in one roll
PayoutNumerator   int         // For payout calculations (e.g., 7 for 7:6)
PayoutDenominator int         // For payout calculations (e.g., 6 for 7:6)
ValidNumbers      []int       // For bets that work on specific numbers
RequiresPoint     bool        // True if bet only works when point is established
RequiresComeOut   bool        // True if bet only works during come out
HouseEdge         float64     // House edge percentage
Commission        float64     // Commission rate (0.05 for 5%)
```

---

## Bet Categories & Types

Below is a complete, canonical list of all supported bet types, grouped by category. Each bet is available for dynamic strategies, validation, and payout calculation.

### **Line Bets**
- PASS_LINE
- DONT_PASS

### **Come Bets**
- COME
- DONT_COME

### **Odds Bets**
- PASS_ODDS
- DONT_PASS_ODDS
- COME_ODDS
- DONT_COME_ODDS

### **Field Bets**
- FIELD

### **Place Bets**
- PLACE_4
- PLACE_5
- PLACE_6
- PLACE_8
- PLACE_9
- PLACE_10

### **Buy Bets**
- BUY_4
- BUY_5
- BUY_6
- BUY_8
- BUY_9
- BUY_10

### **Lay Bets**
- LAY_4
- LAY_5
- LAY_6
- LAY_8
- LAY_9
- LAY_10

### **Place-to-Lose Bets**
- PLACE_TO_LOSE_4
- PLACE_TO_LOSE_5
- PLACE_TO_LOSE_6
- PLACE_TO_LOSE_8
- PLACE_TO_LOSE_9
- PLACE_TO_LOSE_10

### **Hard Way Bets**
- HARD_4
- HARD_6
- HARD_8
- HARD_10

### **Proposition Bets**
- ANY_SEVEN
- ANY_CRAPS
- ELEVEN
- ACE_DEUCE
- ACES
- BOXCARS

### **Horn Bets**
- HORN
- HORN_HIGH_2
- HORN_HIGH_3
- HORN_HIGH_11
- HORN_HIGH_12

### **Hop Bets (Easy Hops)**
- HOP
- HOP_HARD_6
- HOP_EASY_8
- HOP_1_2
- HOP_1_3
- HOP_1_4
- HOP_1_5
- HOP_1_6
- HOP_2_3
- HOP_2_4
- HOP_2_5
- HOP_2_6
- HOP_3_4
- HOP_3_5
- HOP_3_6
- HOP_4_5
- HOP_4_6
- HOP_5_6

### **Big Bets**
- BIG_6
- BIG_8

### **Combination Bets**
- WORLD
- C_AND_E

---

## Core API Functions

### Table Management
```go
// Create a new craps table
func NewTable(minBet, maxBet float64, maxOdds int) *Table

// Add a player to the table
func AddPlayer(table *Table, id, name string, bankroll float64) error

// Remove a player from the table
func RemovePlayer(table *Table, id string) error

// Get a player by ID
func GetPlayer(table *Table, id string) (*Player, error)
```

### Game State
```go
// Roll the dice and resolve all bets
func RollDice(table *Table) *Roll

// Check if we're in come-out phase
func IsComeOut(table *Table) bool

// Check if we have a point established
func IsPoint(table *Table) bool

// Get the current point
func GetPoint(table *Table) Point

// Get the current shooter
func GetShooter(table *Table) string

// Get the current game state
func GetState(table *Table) GameState
```

### Betting
```go
// Place a bet directly (Go API)
func PlaceBet(table *Table, playerID, betType string, amount float64) (*Bet, error)

// Execute CrapsQL commands
func ExecuteString(input string, table *Table) ([]string, error)
```

### Interpreter
```go
// Create a new interpreter
func NewInterpreter(table *Table) *Interpreter

// Execute CrapsQL string
func (i *Interpreter) ExecuteString(input string) ([]string, error)

// Execute CrapsQL string for specific player
func (i *Interpreter) ExecuteStringForPlayer(input string, playerID string) ([]string, error)
```

---

## Dynamic Bet Info & SHOW BETS Command

You can query all available bet types and their details at runtime using the CrapsQL command:

```sql
SHOW BETS;
```

This will return a categorized list of all bet types, with their descriptions and payout ratios, as defined in `canonical_bets.go`.

---

## Example: Querying Bet Info

```go
results, err := interpreter.ExecuteString("SHOW BETS;")
for _, line := range results {
    fmt.Println(line)
}
```

**Sample Output:**
```
Available bet types:

Line Bets:
  PASS_LINE - Bet that shooter will win (7 or 11 on come out, then make point) (Payout: 1:1)
  DONT_PASS - Bet that shooter will lose (2, 3 on come out, then seven out) (Payout: 1:1 (12 is push))
... (other categories and bets)
```

---

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "CrapsQL/pkg/crapsql"
)

func main() {
    // Create table
    table := crapsql.NewTable(5.0, 100.0, 3)
    
    // Add player
    err := crapsql.AddPlayer(table, "player1", "John", 1000.0)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create interpreter
    interpreter := crapsql.NewInterpreter(table)
    
    // Place bets via CrapsQL
    results, err := interpreter.ExecuteString(`
        PLACE $25 ON PASS_LINE;
        PLACE $10 ON FIELD;
        SHOW BETS;
    `)
    if err != nil {
        log.Fatal(err)
    }
    
    // Print results
    for _, result := range results {
        fmt.Println(result)
    }
    
    // Roll dice
    roll := crapsql.RollDice(table)
    fmt.Printf("Rolled: %d-%d = %d\n", roll.Die1, roll.Die2, roll.Total)
    
    // Check game state
    if crapsql.IsComeOut(table) {
        fmt.Println("Come out roll")
    } else {
        fmt.Printf("Point is: %d\n", crapsql.GetPoint(table))
    }
}
```

---

## Error Handling

- Placing an invalid or unsupported bet will return an error: `unknown bet type`.
- Attempting to place a bet in the wrong game state (e.g., odds with no point) will return a descriptive error.
- All errors are surfaced via the API and CrapsQL interpreter.
- The system includes comprehensive validation and error recovery.

---

## Best Practices

- Use `SHOW BETS` to discover all available bets and their details.
- Reference bet types by their canonical string names as listed above.
- All bet logic, validation, and payout is driven by `canonical_bets.go`.
- House edge and commission are for reference and validation.
- Use the interpreter for complex betting strategies and conditional logic.

---

## Testing

The package includes comprehensive tests covering all bet types and scenarios:

```bash
# Run all tests
go test ./pkg/crapsql/ -v

# Run specific test categories
go test ./pkg/crapsql/ -run TestBasicBetPlacement -v
go test ./pkg/crapsql/ -run TestBetResolution -v
```

---

## For More Information

- See `docs/QuickStart.md` for usage examples
- See `docs/CrapsQL.md` for language syntax
- See `docs/jargon.md` for terminology
- See `pkg/crapsgame/canonical_bets.go` for the canonical bet definitions 