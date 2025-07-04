# CrapsQL Language Reference

*A Domain-Specific Language for Craps Simulation*

---

## Table of Contents

- [Overview](#-overview)
- [Getting Started](#-getting-started)
- [Language Syntax](#-language-syntax)
- [Core Statements](#-core-statements)
- [Bet Types Reference](#-bet-types-reference)
- [Game Management](#-game-management)
- [Advanced Features](#-advanced-features)
- [Error Handling](#-error-handling)
- [Examples](#-examples)
- [Best Practices](#-best-practices)

---

## Overview

**CrapsQL** is a domain-specific language designed for craps table simulation and betting strategy development. It combines SQL-like syntax with craps-specific operations, allowing you to:

- **Place and manage bets** with intuitive commands
- **Simulate complete craps games** with random dice rolling
- **Implement complex betting strategies** using conditional logic
- **Track game state and player statistics** in real-time
- **Test betting systems** before using them in live play

### Key Features

- **Complete Bet Coverage**: Support for all standard craps bets plus exotic combinations
- **Real Game Simulation**: Accurate dice rolling and payout calculations
- **Bankroll Management**: Built-in player and table management
- **Error Prevention**: Comprehensive validation and helpful error messages

---

## Getting Started

### Basic Structure

CrapsQL statements are executed sequentially and end with semicolons (`;`). The language is case-sensitive - all commands are UPPERCASE

```sql
-- Basic bet placement
PLACE $25 ON PASS_LINE;

-- Check current game state
SHOW POINT;

-- Roll the dice
ROLL DICE;
```

### Your First Game

```sql
-- Set up your bankroll
SET BANKROLL = $1000;

-- Place a basic pass line bet
PLACE $25 ON PASS_LINE;

-- Roll the dice to see what happens
ROLL DICE;

-- Check your results
SHOW BANKROLL;
```

---

## 📝 Language Syntax

### Basic Operators

| Operator | Description           | Example            |
|----------|-----------------------|--------------------|
|   `=`    | Equal to              | `POINT = 6`        |
|   `!=`   | Not equal to          | `POINT != 4`       |
|   `>`    | Greater than          | `BANKROLL > 500`   |
|   `<`    | Less than             | `BANKROLL < 100`   |
|   `>=`   | Greater than or equal | `BANKROLL >= 1000` |
|   `<=`   | Less than or equal    | `BANKROLL <= 2000` |

### Amount Notation

```sql
$25         -- Explicit dollar amount
$25.50      -- Decimal amounts supported
$5.00       -- Both formats acceptable
```

---

## 💰 Core Statements

### 1. Bet Placement

#### Basic Syntax
```sql
PLACE $amount ON bet_type;
```

#### Examples
```sql
PLACE $25 ON PASS_LINE;
PLACE $10 ON FIELD;
PLACE $12 ON PLACE_6;
PLACE $5 ON HARD_8;
PLACE $15 ON ANY_SEVEN;
```

#### Combination Bets
```sql
-- Place multiple numbers at once
PLACE $32 ON PLACE_INSIDE;     -- Covers 5, 6, 8, 9
PLACE $44 ON PLACE_OUTSIDE;    -- Covers 4, 5, 9, 10
PLACE $5 ON ALL_HARDWAYS;      -- Covers all hard ways

-- Horn bets with high numbers
PLACE $5 ON HORN_HIGH_11;      -- Horn bet with extra on 11
PLACE $4 ON HORN;              -- Standard horn bet
```

### 2. Dice Rolling

```sql
ROLL DICE;                     -- Roll the dice and resolve all bets
```

When you roll the dice, the system will:
- Generate random dice results
- Resolve all active bets
- Update game state (point establishment/resolution)
- Display win/loss results
- Update player bankrolls

### 3. Bet Management

#### Remove Bets
```sql
REMOVE ALL;                    -- Remove all bets and return money
REMOVE PLACE_6;               -- Remove specific bet type
```

#### Press Bets (Increase Amount)
```sql
PRESS PLACE_6 BY $6;          -- Increase Place 6 bet by $6
PRESS HARD_8 BY $5;           -- Press hard 8 by $5
```

#### Turn Bets On/Off
```sql
TURN ON PLACE_6;              -- Make bet active for next roll
TURN OFF PLACE_6;             -- Make bet inactive for next roll
```

### 4. Query Statements

#### Game State Queries
```sql
SHOW POINT;                   -- Display current point
SHOW BANKROLL;                -- Show your current bankroll
SHOW BETS;                    -- List all available bet types
SHOW TABLE_MINIMUMS;          -- Display table limits
```

---

## 🎲 Bet Types Reference

### Line Bets
*The foundation of craps betting*

| Bet Type | Description | Payout | House Edge |
|----------|-------------|--------|------------|
| `PASS_LINE` | Win on 7/11, lose on 2/3/12, then make point | 1:1 | 1.41% |
| `DONT_PASS` | Opposite of pass line (12 pushes) | 1:1 | 1.36% |

### Come Bets
*Similar to line bets but placed after come-out*

| Bet Type | Description | Payout | House Edge |
|----------|-------------|--------|------------|
| `COME` | Personal pass line bet | 1:1 | 1.41% |
| `DONT_COME` | Personal don't pass bet | 1:1 | 1.36% |

### Odds Bets
*The best bets on the table - no house edge!*

| Bet Type | Description | Payout | House Edge |
|----------|-------------|--------|------------|
| `PASS_ODDS` | Odds behind pass line | True odds | 0.00% |
| `DONT_PASS_ODDS` | Odds behind don't pass | True odds | 0.00% |
| `COME_ODDS` | Odds behind come bet | True odds | 0.00% |
| `DONT_COME_ODDS` | Odds behind don't come | True odds | 0.00% |

### Place Bets
*Bet that a number will roll before 7*

| Bet Type | Number | Payout | House Edge |
|----------|--------|--------|------------|
| `PLACE_4` | 4 | 9:5 | 6.67% |
| `PLACE_5` | 5 | 7:5 | 4.00% |
| `PLACE_6` | 6 | 7:6 | 1.52% |
| `PLACE_8` | 8 | 7:6 | 1.52% |
| `PLACE_9` | 9 | 7:5 | 4.00% |
| `PLACE_10` | 10 | 9:5 | 6.67% |

#### Place Bet Combinations
| Bet Type | Numbers Covered | Description |
|----------|-----------------|-------------|
| `PLACE_INSIDE` | 5, 6, 8, 9 | Inside numbers |
| `PLACE_OUTSIDE` | 4, 5, 9, 10 | Outside numbers |
| `PLACE_NUMBERS` | Custom | Specify which numbers |

### Buy Bets
*Pay commission for true odds*

| Bet Type | Number | Payout | Commission |
|----------|--------|--------|------------|
| `BUY_4` | 4 | 2:1 | 5% |
| `BUY_5` | 5 | 3:2 | 5% |
| `BUY_6` | 6 | 6:5 | 5% |
| `BUY_8` | 8 | 6:5 | 5% |
| `BUY_9` | 9 | 3:2 | 5% |
| `BUY_10` | 10 | 2:1 | 5% |

### Lay Bets
*Bet against numbers (7 before number)*

| Bet Type | Number | Payout | Commission |
|----------|--------|--------|------------|
| `LAY_4` | 4 | 1:2 | 5% |
| `LAY_5` | 5 | 2:3 | 5% |
| `LAY_6` | 6 | 5:6 | 5% |
| `LAY_8` | 8 | 5:6 | 5% |
| `LAY_9` | 9 | 2:3 | 5% |
| `LAY_10` | 10 | 1:2 | 5% |

### Place-to-Lose Bets
*Like lay bets but without commission*

| Bet Type | Number | Payout | House Edge |
|----------|--------|--------|------------|
| `PLACE_TO_LOSE_4` | 4 | 1:2 | 2.44% |
| `PLACE_TO_LOSE_5` | 5 | 2:3 | 3.23% |
| `PLACE_TO_LOSE_6` | 6 | 5:6 | 4.00% |
| `PLACE_TO_LOSE_8` | 8 | 5:6 | 4.00% |
| `PLACE_TO_LOSE_9` | 9 | 2:3 | 3.23% |
| `PLACE_TO_LOSE_10` | 10 | 1:2 | 2.44% |

### Field Bets
*One-roll bet on specific numbers*

| Bet Type | Description | Payout | House Edge |
|----------|-------------|--------|------------|
| `FIELD` | Win on 2,3,4,9,10,11,12 | 1:1 (2 pays 2:1, 12 pays 3:1) | 2.78% |

### Hard Way Bets
*Both dice must show the same number*

| Bet Type | Number | Required Roll | Payout | House Edge |
|----------|--------|---------------|--------|------------|
| `HARD_4` | 4 | 2-2 | 7:1 | 11.11% |
| `HARD_6` | 6 | 3-3 | 9:1 | 9.09% |
| `HARD_8` | 8 | 4-4 | 9:1 | 9.09% |
| `HARD_10` | 10 | 5-5 | 7:1 | 11.11% |
| `ALL_HARDWAYS` | All | All hard ways | Various | Various |

### Proposition Bets
*One-roll bets with high payouts*

| Bet Type | Description | Payout | House Edge |
|----------|-------------|--------|------------|
| `ANY_SEVEN` | Next roll is 7 | 4:1 | 16.67% |
| `ANY_CRAPS` | Next roll is 2, 3, or 12 | 7:1 | 11.11% |
| `ELEVEN` | Next roll is 11 | 15:1 | 11.11% |
| `ACE_DEUCE` | Next roll is 3 (1-2) | 15:1 | 11.11% |
| `ACES` | Next roll is 2 (1-1) | 30:1 | 13.89% |
| `BOXCARS` | Next roll is 12 (6-6) | 30:1 | 13.89% |

### Horn Bets
*Combination bets on 2, 3, 11, 12*

| Bet Type | Description | Special Payout |
|----------|-------------|----------------|
| `HORN` | Equal on 2, 3, 11, 12 | Standard payouts |
| `HORN_HIGH_2` | Extra on 2 | 2 pays 27:4 |
| `HORN_HIGH_3` | Extra on 3 | 3 pays 15:1 |
| `HORN_HIGH_11` | Extra on 11 | 11 pays 15:1 |
| `HORN_HIGH_12` | Extra on 12 | 12 pays 27:4 |

### Hop Bets
*Bet on exact dice combinations*

| Bet Type | Description | Payout |
|----------|-------------|--------|
| `HOP_1_2` | Exactly 1-2 | 15:1 |
| `HOP_1_3` | Exactly 1-3 | 15:1 |
| `HOP_1_4` | Exactly 1-4 | 15:1 |
| `HOP_1_5` | Exactly 1-5 | 15:1 |
| `HOP_1_6` | Exactly 1-6 | 15:1 |
| `HOP_2_3` | Exactly 2-3 | 15:1 |
| `HOP_2_4` | Exactly 2-4 | 15:1 |
| `HOP_2_5` | Exactly 2-5 | 15:1 |
| `HOP_2_6` | Exactly 2-6 | 15:1 |
| `HOP_3_4` | Exactly 3-4 | 15:1 |
| `HOP_3_5` | Exactly 3-5 | 15:1 |
| `HOP_3_6` | Exactly 3-6 | 15:1 |
| `HOP_4_5` | Exactly 4-5 | 15:1 |
| `HOP_4_6` | Exactly 4-6 | 15:1 |
| `HOP_5_6` | Exactly 5-6 | 15:1 |
| `HOP_HARD_6` | Hard 6 (3-3) | 30:1 |
| `HOP_EASY_8` | Easy 8 (not 4-4) | 15:1 |

### Big Bets
*Alternative to place bets (worse odds)*

| Bet Type | Description | Payout | House Edge |
|----------|-------------|--------|------------|
| `BIG_6` | 6 before 7 | 1:1 | 9.09% |
| `BIG_8` | 8 before 7 | 1:1 | 9.09% |

### Combination Bets

| Bet Type | Description | Coverage |
|----------|-------------|----------|
| `WORLD` | Any 7 + Any Craps | 2, 3, 7, 11, 12 |
| `C_AND_E` | Craps + Eleven | 2, 3, 11, 12 |

---

## 🎮 Game Management

### Bankroll Management

```sql
-- Set your starting bankroll
SET BANKROLL = $1000;

-- Set betting limits
SET MAX_BET = $100;
SET MIN_BET = $5;

-- Set win/loss goals
SET WIN_GOAL = $1500;
SET LOSS_LIMIT = $500;
```

### Table Settings

```sql
-- View current table limits
SHOW TABLE_MINIMUMS;

-- Table limits are set when creating the game engine
-- but can be queried for strategy decisions
```

---

## 🧠 Advanced Features

### Conditional Logic

Execute statements based on game conditions:

```sql
IF condition THEN
    statement1;
    statement2;
ELSE
    alternative_statement;
END;
```

#### Available Conditions

| Condition | Description | Example |
|-----------|-------------|---------|
| `POINT = number` | Check if point equals number | `IF POINT = 6 THEN` |
| `POINT != number` | Check if point doesn't equal | `IF POINT != 4 THEN` |
| `BANKROLL > amount` | Check bankroll level | `IF BANKROLL > 500 THEN` |
| `BANKROLL < amount` | Check if running low | `IF BANKROLL < 100 THEN` |

#### Advanced Conditional Examples

```sql
-- Aggressive odds strategy
IF POINT = 6 OR POINT = 8 THEN
    PLACE $30 ON PASS_ODDS;
ELSE
    PLACE $20 ON PASS_ODDS;
END;

-- Conservative bankroll management
IF BANKROLL < 200 THEN
    REMOVE ALL;
    PLACE $5 ON PASS_LINE;
ELSE
    PLACE $25 ON PASS_LINE;
    PLACE $10 ON FIELD;
END;

-- Point-specific place betting
IF POINT = 4 THEN
    PLACE $12 ON PLACE_6;
    PLACE $12 ON PLACE_8;
END;
```

### Strategy Examples

#### The Iron Cross
*Cover most numbers with field + place bets*

```sql
PLACE $10 ON FIELD;
PLACE $12 ON PLACE_5;
PLACE $12 ON PLACE_6;
PLACE $12 ON PLACE_8;
-- Covers everything except 7
```

#### Conservative Pass Line Strategy
*Basic strategy with odds*

```sql
-- Come out roll
PLACE $25 ON PASS_LINE;
ROLL DICE;

-- If point established, take odds
IF POINT != 0 THEN
    IF POINT = 6 OR POINT = 8 THEN
        PLACE $30 ON PASS_ODDS;  -- 6x odds on 6/8
    ELSE
        PLACE $50 ON PASS_ODDS;  -- Higher odds on other points
    END;
END;
```

#### Aggressive Place Betting
*Cover all place numbers*

```sql
-- Wait for point to be established
IF POINT != 0 THEN
    PLACE $10 ON PLACE_4;
    PLACE $10 ON PLACE_5;
    PLACE $12 ON PLACE_6;
    PLACE $12 ON PLACE_8;
    PLACE $10 ON PLACE_9;
    PLACE $10 ON PLACE_10;
END;
```

---

## ⚠️ Error Handling

### Common Error Types

#### 1. Syntax Errors
```sql
PLACE $25 ON INVALID_BET;
-- Error: Unknown bet type: INVALID_BET
```

#### 2. Game State Errors
```sql
PLACE $25 ON PASS_ODDS;
-- Error: Cannot place pass odds - no point established
```

#### 3. Bankroll Errors
```sql
PLACE $1000 ON PASS_LINE;
-- Error: Insufficient bankroll: $500 available, $1000 required
```

#### 4. Table Limit Errors
```sql
PLACE $2000 ON PASS_LINE;
-- Error: Bet amount $2000 exceeds maximum $1000
```

### Error Recovery

The system provides:
- **Detailed error messages** with specific reasons
- **Suggestions for correction** when possible
- **Line numbers** for syntax errors
- **Graceful handling** - other bets continue to work

---

## 📚 Examples

### Complete Game Session

```sql
-- Start a new session
SET BANKROLL = $1000;

-- Place initial come-out bet
PLACE $25 ON PASS_LINE;
PLACE $5 ON FIELD;

-- Roll to establish point
ROLL DICE;

-- If point is 6 or 8, play aggressively
IF POINT = 6 OR POINT = 8 THEN
    PLACE $30 ON PASS_ODDS;
    PLACE $18 ON PLACE_6;
    PLACE $18 ON PLACE_8;
END;

-- If point is 4 or 10, play conservatively
IF POINT = 4 OR POINT = 10 THEN
    PLACE $50 ON PASS_ODDS;  -- Better odds on outside numbers
END;

-- Continue rolling until point resolution
ROLL DICE;
ROLL DICE;

-- Check results
SHOW BANKROLL;
SHOW POINT;
```

### Martingale System Example

```sql
-- Start with base bet
SET START_BET = $10;
SET CURRENT_BET = $10;

-- Place initial bet
PLACE $CURRENT_BET ON PASS_LINE;
ROLL DICE;

-- If lost, double the bet (simplified - real implementation would need more logic)
IF BANKROLL < PREVIOUS_BANKROLL THEN
    SET CURRENT_BET = CURRENT_BET * 2;
ELSE
    SET CURRENT_BET = START_BET;
END;
```

### Hedge Betting Example

```sql
-- Place pass line bet
PLACE $30 ON PASS_LINE;

-- Hedge with any craps on come out
PLACE $3 ON ANY_CRAPS;

ROLL DICE;

-- If point established, remove hedge and add odds
IF POINT != 0 THEN
    PLACE $40 ON PASS_ODDS;
END;
```

---

## 🎲 Quick Reference

### Essential Commands
```sql
PLACE $amount ON bet_type;     -- Place a bet
ROLL DICE;                     -- Roll the dice
SHOW POINT;                    -- Check current point
SHOW BANKROLL;                 -- Check your money
REMOVE ALL;                    -- Clear all bets
```

### Best Bets (Lowest House Edge)
1. **Pass/Don't Pass** (1.36-1.41%)
2. **Pass/Don't Pass Odds** (0.00%)
3. **Come/Don't Come** (1.36-1.41%)
4. **Place 6/8** (1.52%)

### Avoid These Bets (High House Edge)
- **Any Seven** (16.67%)
- **Hard Ways** (9.09-11.11%)
- **Proposition Bets** (11.11-16.67%)
- **Field** (2.78%)

---

*Happy rolling! 🎲* 