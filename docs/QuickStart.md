# CrapsQL Quick Start Guide

## Overview

CrapsQL is a complete craps table simulation with a custom domain-specific language (DSL) for placing bets and managing game state. This guide will help you get started quickly.

## Installation

### Prerequisites
- Go 1.23 or later

### Setup
```bash
go get github.com/headswim/crapsql
```

## Basic Usage

### 1. Create a Table
```go
package main

import "github.com//crapsql"

func main() {
    // Create a table with $5 minimum, $500 maximum, 3x odds
    table := crapsql.NewTable(5.0, 500.0, 3)
    
    // Add a player
    err := crapsql.AddPlayer(table, "player1", "John", 1000.0)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. Create an Interpreter
```go
// Create interpreter for executing CrapsQL commands
interpreter := crapsql.NewInterpreter(table)
```

### 3. Place Your First Bet
```go
// Place a pass line bet
results, err := interpreter.ExecuteString("PLACE $25 ON PASS_LINE;")
if err != nil {
    log.Fatal(err)
}

for _, result := range results {
    fmt.Println(result)
}
```

### 4. Roll the Dice
```go
// Roll the dice and resolve all bets
roll := table.RollDice()
fmt.Printf("Rolled: %d-%d (%d)\n", roll.Die1, roll.Die2, roll.Total)
```

## Common Commands

### Basic Betting
```sql
-- Pass line bet
PLACE $25 ON PASS_LINE;

-- Field bet
PLACE $10 ON FIELD;

-- Place bet on 6
PLACE $12 ON PLACE_6;

-- Hard 8
PLACE $5 ON HARD_8;

-- Place-to-lose bet on 10
PLACE $10 ON PLACE_TO_LOSE_10;

-- Buy bet on 4
PLACE $20 ON BUY_4;

-- Lay bet against 8
PLACE $15 ON LAY_8;

-- Hop bet (easy hop 1-2)
PLACE $2 ON HOP_1_2;

-- Horn bet (horn high 12)
PLACE $4 ON HORN_HIGH_12;
```

### Check Game State
```sql
-- Show current table state
SHOW TABLE;

-- Show your bets
SHOW BETS;

-- Show your bankroll
SHOW BANKROLL;
```

### Advanced Features
```sql
-- Conditional betting
IF POINT = 6 THEN
    PLACE $30 ON PASS_ODDS;
END;
```

### Bet Management
```sql
-- Remove all place 6 bets
REMOVE ALL PLACE_6;

-- Press (increase) place 8 bet by $10
PRESS PLACE_8 BY $10;

-- Turn off hard 8 bet (make it not working)
TURN OFF HARD_8;
```

### Discovering Bets and Payouts

You can list all available bet types, their descriptions, and payout ratios:

```sql
SHOW BETS;
```

**Sample Output:**
```
Available bet types:

Line Bets:
  PASS_LINE - Bet that shooter will win (7 or 11 on come out, then make point) (Payout: 1:1)
  DONT_PASS - Bet that shooter will lose (2, 3 on come out, then seven out) (Payout: 1:1 (12 is push))
... (other categories and bets)
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/headswim/crapsql"
)

func main() {
    // Setup
    table := crapsql.NewTable(5.0, 500.0, 3)
    table.AddPlayer("player1", "John", 1000.0)
    interpreter := crapsql.NewInterpreter(table)
    
    // Place initial bets
    commands := []string{
        "PLACE $25 ON PASS_LINE;",
        "PLACE $10 ON FIELD;",
        "PLACE $2 ON HOP_1_2;",
        "SHOW BETS;",
    }
    
    for _, cmd := range commands {
        results, err := interpreter.ExecuteString(cmd)
        if err != nil {
            log.Printf("Error: %v", err)
            continue
        }
        
        for _, result := range results {
            fmt.Println(result)
        }
    }
    
    // Roll dice
    roll := table.RollDice()
    fmt.Printf("🎲 Rolled: %d-%d (%d)\n", roll.Die1, roll.Die2, roll.Total)
    
    // Check results
    results, _ := interpreter.ExecuteString("SHOW TABLE;")
    for _, result := range results {
        fmt.Println(result)
    }
}
```

## Testing

The package includes a comprehensive test suite. To run tests:

```bash
# Run all tests
go test ./pkg/crapsql/ -v

# Run specific test
go test ./pkg/crapsql/ -run TestBasicBetPlacement -v

# Run tests with coverage
go test ./pkg/crapsql/ -cover
```

## Supported Bet Types

The system supports 50+ bet types including:

- **Line Bets:** PASS_LINE, DONT_PASS, COME, DONT_COME
- **Field Bets:** FIELD (one-roll)
- **One-Roll Bets:** ANY_SEVEN, ANY_CRAPS, ELEVEN, ACE_DEUCE, ACES, BOXCARS
- **Place Bets:** PLACE_4, PLACE_5, PLACE_6, PLACE_8, PLACE_9, PLACE_10
- **Hard Ways:** HARD_4, HARD_6, HARD_8, HARD_10
- **Odds Bets:** PASS_ODDS, DONT_PASS_ODDS, COME_ODDS, DONT_COME_ODDS
- **Buy/Lay Bets:** BUY_4, BUY_5, BUY_6, BUY_8, BUY_9, BUY_10, LAY_4, LAY_5, LAY_6, LAY_8, LAY_9, LAY_10
- **Place-to-Lose Bets:** PLACE_TO_LOSE_4, PLACE_TO_LOSE_5, etc.
- **Hop Bets:** HOP_1_2, HOP_1_3, HOP_2_3, etc. (all combinations)
- **Horn Bets:** HORN, HORN_HIGH_2, HORN_HIGH_3, HORN_HIGH_11, HORN_HIGH_12
- **Proposition Bets:** WORLD, C_AND_E

## Next Steps

- Check out the [API Reference](API.md) for detailed function documentation
- Read the [CrapsQL Language Guide](CrapsQL.md) for advanced language features
- Review [Craps Terminology](jargon.md) to understand the game concepts

## Getting Help

- **Language Reference**: `docs/CrapsQL.md`
- **API Documentation**: `docs/API.md`
- **Examples**: See the `examples/` directory
- **Tests**: Run `go test ./...` to see usage examples

## Common Issues

### "Unknown bet type" Error
Make sure you're using the correct bet type names. See the Bet Types Reference in the language guide or use `SHOW BETS;`.

### "No point established" Error
Some bets (like odds) can only be placed after a point is established. Check the game state first.

### "Bet exceeds maximum" Error
Your bet amount is above the table's maximum. Check the table limits with `SHOW TABLE;`. 