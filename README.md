# CrapsQL - A Domain-Specific Language for Craps Game Simulation

## Overview

CrapsQL is a complete craps table simulation with a custom domain-specific language (DSL) for placing bets and managing game state. It consists of two main packages: `crapsql` (the DSL interpreter) and `crapsgame` (the core game engine).

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Testing](#testing)
- [Documentation](#documentation)

## Project Structure

### Package: `crapsql` - The CrapsQL Language Interpreter

The `crapsql` package implements a complete language stack for interacting with a craps table:

#### Core Components:

- **`crapsql.go`** - Main package interface and convenience functions
- **`lexer.go`** - Tokenizes CrapsQL input into lexical tokens
- **`parser.go`** - Parses tokens into an Abstract Syntax Tree (AST)
- **`interpreter.go`** - Executes parsed statements against the game state
- **`types.go`** - Defines all language constructs, tokens, and AST nodes
- **`validation.go`** - Validates bets and game state
- **`bet_registry.go`** - Maps bet type strings to enum values
- **`naming_conventions.go`** - Bet type naming validation
- **`crapsql_test.go`** - Comprehensive test suite (1800+ lines)

#### Language Features:

**Bet Placement Commands:**
```sql
PLACE $25 ON PASS_LINE;
PLACE $10 ON FIELD;
PLACE $15 ON PLACE_6;
PLACE $5 ON HARD_8;
```

**Query Commands:**
```sql
SHOW POINT;
SHOW BETS;
SHOW BANKROLL;
SHOW TABLE_MINIMUMS;
```

**Conditional Logic:**
```sql
IF POINT = 6 THEN
    PLACE $20 ON PLACE_6;
    PLACE $10 ON HARD_6;
END;
```

**Bet Management:**
```sql
REMOVE ALL PLACE_6;
PRESS PLACE_8 BY $10;
TURN OFF HARD_8;
```

#### Supported Bet Types:

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

### Package: `crapsgame` - The Core Game Engine

The `crapsgame` package provides the foundational craps table simulation:

#### Core Components:

- **`state.go`** - Game state management, player handling, and dice rolling (1581 lines)
- **`bets.go`** - Comprehensive bet resolution system with all payout calculations
- **`canonical_bets.go`** - Single source of truth for all bet definitions (1120 lines)

#### Key Features:

**Game State Management:**
- Tracks come-out vs. point phases
- Manages point numbers (4, 5, 6, 8, 9, 10)
- Handles shooter rotation
- Maintains player sessions

**Player Management:**
- Add/remove players with bankrolls
- Track individual betting limits
- Session management with win goals and loss limits
- Automatic shooter assignment

**Bet Resolution System:**
- Real-time bet resolution on each roll
- Accurate payout calculations for all bet types
- Proper odds handling (3x, 5x, etc.)
- Working/off bet modifiers
- Cryptographically secure random number generation

**Comprehensive Bet Types:**
- All standard casino craps bets with correct house edges
- Proper payout ratios and commission handling
- State-dependent bet validation
- Come bet lifecycle management
- Odds bet integration

## How It Works

### 1. Table Creation
```go
table := crapsql.NewTable(5.0, 100.0, 3) // $5 min, $100 max, 3x odds
```

### 2. Player Management
```go
crapsql.AddPlayer(table, "player1", "John", 1000.0) // player, Name, bankroll
```

### 3. Betting via CrapsQL
```go
interpreter := crapsql.NewInterpreter(table)
results, err := interpreter.ExecuteString("PLACE $25 ON PASS_LINE;")
```

### 4. Game Progression
```go
roll := crapsql.RollDice(table) // Automatically resolves all bets
```

### 5. State Queries
```go
if crapsql.IsComeOut(table) {
    fmt.Println("Come out roll")
} else {
    fmt.Printf("Point is: %d\n", crapsql.GetPoint(table))
}
```

## Testing

The package includes a comprehensive test suite with 30+ test functions covering:

- Basic bet placement and validation
- All bet types (50+ bet types tested)
- Bet resolution scenarios
- Multi-player scenarios
- Conditional logic and strategies
- Error handling and edge cases
- Performance and stress testing

To run the tests:
```bash
go test ./pkg/crapsql/ -v
```

## Documentation

- [API Reference](docs/API.md) - Detailed API documentation
- [CrapsQL Language Guide](docs/CrapsQL.md) - Language syntax and features
- [Quick Start Guide](docs/QuickStart.md) - Getting started examples
- [Craps Terminology](docs/jargon.md) - Craps game terminology

## Technical Architecture

### Language Design
- **Lexical Analysis:** Custom lexer with craps-specific tokens
- **Parsing:** Recursive descent parser with error recovery
- **AST:** Rich abstract syntax tree for complex expressions
- **Interpreter:** Direct execution with game state integration

### Game Engine
- **State Machine:** Clean state transitions (come-out ↔ point)
- **Bet Resolution:** Event-driven bet resolution system
- **Player Management:** Concurrent player support
- **Security:** Cryptographically secure random number generation

### Key Features
- **Complete Rule Set:** All standard craps rules implemented
- **Accurate Math:** Precise payout calculations with proper house edges
- **Multi-Player:** Support for multiple players with individual bankrolls
- **Session Management:** Win goals, loss limits, time tracking
- **Audit Trail:** Complete bet and roll history
- **Error Recovery:** Graceful error handling with detailed messages 