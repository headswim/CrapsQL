# CrapsQL Language Reference

## Overview

CrapsQL is a domain-specific language designed for craps table simulation. It provides SQL-like syntax for placing bets, managing game state, and defining complex betting strategies.

## Language Syntax

### Basic Structure

CrapsQL statements end with semicolons (`;`) and are case-insensitive. Keywords are typically written in UPPERCASE for clarity.

### Comments

```sql
-- Single line comment
/* Multi-line comment */
```

## Core Statements

### 1. Bet Placement

#### Basic Bet Syntax
```sql
PLACE $amount ON bet_type;
```

**Examples:**
```sql
PLACE $25 ON PASS_LINE;
PLACE $10 ON FIELD;
PLACE $12 ON PLACE_6;
PLACE $5 ON HARD_8;
```

#### Bet Modifiers
```sql
PLACE $amount ON bet_type WITH modifier;
```

**Available Modifiers:**
- `WORKING` - Bet works on come-out rolls
- `NOT WORKING` - Bet doesn't work on come-out rolls

**Examples:**
```sql
PLACE $10 ON PLACE_6 WITH WORKING;
PLACE $5 ON FIELD WITH NOT WORKING;
```

### 2. Query Statements

#### Show Commands
```sql
SHOW TABLE;           -- Show current table state
SHOW PLAYERS;         -- Show all players
SHOW BETS;            -- Show all active bets
SHOW PLAYER player_id; -- Show specific player
SHOW BET bet_type;    -- Show specific bet type
```

#### Examples:
```sql
SHOW TABLE;
-- Output: Come-out roll, Point: none, Players: 3, Active bets: 5

SHOW PLAYER john_doe;
-- Output: Player: John Doe, Bankroll: $1,250, Active bets: 2
```

### 3. Management Statements

#### Set Commands
```sql
SET variable = value;
```

**Available Variables:**
- `MIN_BET` - Minimum bet amount
- `MAX_BET` - Maximum bet amount
- `MAX_ODDS` - Maximum odds multiplier

**Examples:**
```sql
SET MIN_BET = $10;
SET MAX_BET = $1000;
SET MAX_ODDS = 5;
```

## Advanced Features
### 3. Conditional Logic

Execute statements based on game conditions:

```sql
IF condition THEN
    statement1;
    statement2;
END IF;
```

**Available Conditions:**
- `POINT = number` - Check current point
- `COME_OUT_ROLL` - Check if in come-out phase
- `POINT_ESTABLISHED` - Check if point is established
- `bet_type EXISTS` - Check if bet is active

**Examples:**
```sql
IF POINT = 6 THEN
    PLACE $30 ON PASS_ODDS;
END IF;

IF COME_OUT_ROLL THEN
    PLACE $10 ON FIELD;
END IF;

IF PLACE_6 EXISTS THEN
    PRESS PLACE_6 BY $6;
END IF;
```

### 4. Bet Management

#### Remove Bets
```sql
REMOVE BET ON bet_type;
REMOVE ALL;
```

**Examples:**
```sql
REMOVE BET ON PLACE_6;
REMOVE ALL;
```

#### Press Bets (Increase)
```sql
PRESS bet_type BY $amount;
```

**Examples:**
```sql
PRESS PLACE_6 BY $6;
PRESS FIELD BY $5;
```

#### Turn Bets On/Off
```sql
TURN ON bet_type;
TURN OFF bet_type;
```

**Examples:**
```sql
TURN ON PLACE_6;
TURN OFF FIELD;
```

## Bet Types Reference

### Line Bets
- `PASS_LINE` - Pass line bet
- `DONT_PASS` - Don't pass bet

### Come Bets
- `COME` - Come bet
- `DONT_COME` - Don't come bet

### Field Bets
- `FIELD` - Field bet (2, 3, 4, 9, 10, 11, 12)

### Place Bets
- `PLACE_4` - Place bet on 4
- `PLACE_5` - Place bet on 5
- `PLACE_6` - Place bet on 6
- `PLACE_8` - Place bet on 8
- `PLACE_9` - Place bet on 9
- `PLACE_10` - Place bet on 10

### Hard Ways
- `HARD_4` - Hard 4 (2-2)
- `HARD_6` - Hard 6 (3-3)
- `HARD_8` - Hard 8 (4-4)
- `HARD_10` - Hard 10 (5-5)

### Proposition Bets
- `ANY_7` - Any 7
- `ANY_CRAPS` - Any craps (2, 3, 12)
- `ELEVEN` - Eleven (6-5)
- `ACE_DEUCE` - Ace-deuce (1-2)
- `TWELVE` - Twelve (6-6)

### Odds Bets
- `PASS_ODDS` - Odds behind pass line
- `DONT_PASS_ODDS` - Odds behind don't pass
- `COME_ODDS` - Odds behind come bet
- `DONT_COME_ODDS` - Odds behind don't come

### Buy/Lay Bets
- `BUY_4`, `BUY_5`, `BUY_6`, `BUY_8`, `BUY_9`, `BUY_10`
- `LAY_4`, `LAY_5`, `LAY_6`, `LAY_8`, `LAY_9`, `LAY_10`

## Payout Reference

### Place Bet Payouts
- 4, 10: 9:5
- 5, 9: 7:5
- 6, 8: 7:6

### Field Bet Payouts
- 2: 2:1
- 3, 4, 9, 10, 11: 1:1
- 12: 3:1 (varies by casino)

### Hard Way Payouts
- Hard 4, 10: 7:1
- Hard 6, 8: 9:1

### Proposition Bet Payouts
- Any 7: 4:1
- Any Craps: 7:1
- Eleven: 15:1
- Ace-Deuce: 15:1
- Twelve: 30:1

## Error Handling

### Common Errors

1. **Syntax Errors**
   ```sql
   PLACE $25 ON INVALID_BET;  -- Error: unknown bet type
   ```

2. **Game State Errors**
   ```sql
   PLACE $25 ON PASS_ODDS;    -- Error: no point established
   ```

3. **Bet Limit Errors**
   ```sql
   PLACE $1000 ON PASS_LINE;  -- Error: bet exceeds maximum
   ```

4. **Insufficient Funds**
   ```sql
   PLACE $500 ON PASS_LINE;   -- Error: insufficient bankroll
   ```

### Error Recovery

CrapsQL provides detailed error messages with line numbers and suggestions for correction.

## Best Practices

1. **Use Strategies** for complex betting patterns
2. **Validate Conditions** before placing conditional bets
3. **Monitor Bankroll** with regular SHOW commands
4. **Use Progressive Betting** carefully with proper limits
5. **Test Strategies** before using in live play

## Examples

### Complete Game Session
```sql
-- Add player and place initial bets
PLACE $25 ON PASS_LINE;
PLACE $10 ON FIELD;

-- Roll dice (handled by game engine)
-- Check results
SHOW TABLE;

-- If point established, place odds
IF POINT = 6 THEN
    PLACE $30 ON PASS_ODDS;
END IF;

-- Manage bets
PRESS PLACE_6 BY $6;
TURN OFF FIELD;

-- Clean up
REMOVE ALL;
``` 