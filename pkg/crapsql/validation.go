package crapsql

import (
	"fmt"

	"github.com/headswim/CrapsQL/pkg/crapsgame"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s (value: %v)", e.Field, e.Message, e.Value)
}

// validateBetAmount validates that a bet amount is within acceptable limits
func validateBetAmount(amount float64, minBet float64, maxBet float64) error {
	if amount <= 0 {
		return ValidationError{
			Field:   "amount",
			Message: "bet amount must be positive",
			Value:   amount,
		}
	}

	if amount < minBet {
		return ValidationError{
			Field:   "amount",
			Message: fmt.Sprintf("bet amount $%.2f is below minimum $%.2f", amount, minBet),
			Value:   amount,
		}
	}

	if amount > maxBet {
		return ValidationError{
			Field:   "amount",
			Message: fmt.Sprintf("bet amount $%.2f exceeds maximum $%.2f", amount, maxBet),
			Value:   amount,
		}
	}

	return nil
}

// validateBankroll validates that a player has sufficient bankroll for a bet
func validateBankroll(player *crapsgame.Player, amount float64) error {
	if player == nil {
		return ValidationError{
			Field:   "player",
			Message: "player is nil",
			Value:   nil,
		}
	}

	if amount > player.Bankroll {
		return ValidationError{
			Field:   "bankroll",
			Message: fmt.Sprintf("insufficient bankroll: $%.2f (need $%.2f)", player.Bankroll, amount),
			Value:   player.Bankroll,
		}
	}

	if player.Bankroll < 0 {
		return ValidationError{
			Field:   "bankroll",
			Message: "player has negative bankroll",
			Value:   player.Bankroll,
		}
	}

	return nil
}

// validateBetType validates that a bet type is valid and exists in canonical definitions
func validateBetType(betType string) error {
	if betType == "" {
		return ValidationError{
			Field:   "betType",
			Message: "bet type cannot be empty",
			Value:   betType,
		}
	}

	// Check if bet type exists in canonical definitions
	if _, exists := crapsgame.CanonicalBetDefinitions[betType]; !exists {
		return ValidationError{
			Field:   "betType",
			Message: fmt.Sprintf("unknown bet type: %s", betType),
			Value:   betType,
		}
	}

	// Check if bet type follows naming conventions
	if !IsValidBetTypeName(betType) {
		return ValidationError{
			Field:   "betType",
			Message: fmt.Sprintf("bet type %s does not follow naming conventions", betType),
			Value:   betType,
		}
	}

	return nil
}

// validateGameState validates that a bet type can be placed in the current game state
func validateGameState(betType string, gameState crapsgame.GameState) error {
	// First validate the bet type exists
	if err := validateBetType(betType); err != nil {
		return err
	}

	// Get bet definition
	betDef, exists := crapsgame.CanonicalBetDefinitions[betType]
	if !exists {
		return ValidationError{
			Field:   "betType",
			Message: fmt.Sprintf("bet type %s not found in canonical definitions", betType),
			Value:   betType,
		}
	}

	// Validate game state requirements
	if betDef.RequiresPoint && gameState != crapsgame.StatePoint {
		return ValidationError{
			Field:   "gameState",
			Message: fmt.Sprintf("bet %s requires point to be established (current state: %d)", betType, gameState),
			Value:   gameState,
		}
	}

	if betDef.RequiresComeOut && gameState != crapsgame.StateComeOut {
		return ValidationError{
			Field:   "gameState",
			Message: fmt.Sprintf("bet %s requires come-out roll (current state: %d)", betType, gameState),
			Value:   gameState,
		}
	}

	// Validate game state is valid
	if gameState < crapsgame.StateComeOut || gameState > crapsgame.StateSevenOut {
		return ValidationError{
			Field:   "gameState",
			Message: fmt.Sprintf("invalid game state: %d", gameState),
			Value:   gameState,
		}
	}

	return nil
}

// validateBetPlacement performs comprehensive validation for bet placement
func validateBetPlacement(bet *crapsgame.Bet, player *crapsgame.Player, table *crapsgame.Table) error {
	// Validate bet is not nil
	if bet == nil {
		return ValidationError{
			Field:   "bet",
			Message: "bet cannot be nil",
			Value:   nil,
		}
	}

	// Validate player is not nil
	if player == nil {
		return ValidationError{
			Field:   "player",
			Message: "player cannot be nil",
			Value:   nil,
		}
	}

	// Validate table is not nil
	if table == nil {
		return ValidationError{
			Field:   "table",
			Message: "table cannot be nil",
			Value:   nil,
		}
	}

	// Validate bet amount
	if err := validateBetAmount(bet.Amount, table.MinBet, table.MaxBet); err != nil {
		return err
	}

	// Validate bankroll
	if err := validateBankroll(player, bet.Amount); err != nil {
		return err
	}

	// Validate bet type
	if err := validateBetType(bet.Type); err != nil {
		return err
	}

	// Validate game state
	if err := validateGameState(bet.Type, table.State); err != nil {
		return err
	}

	// Validate bet numbers if present
	if err := validateBetNumbers(bet); err != nil {
		return err
	}

	// Validate bet ID
	if bet.ID == "" {
		return ValidationError{
			Field:   "betID",
			Message: "bet must have a valid ID",
			Value:   bet.ID,
		}
	}

	// Validate bet player matches
	if bet.Player != player.ID {
		return ValidationError{
			Field:   "betPlayer",
			Message: fmt.Sprintf("bet player %s does not match player %s", bet.Player, player.ID),
			Value:   bet.Player,
		}
	}

	return nil
}

// validateBetNumbers validates that bet numbers are appropriate for the bet type
func validateBetNumbers(bet *crapsgame.Bet) error {
	if bet == nil {
		return ValidationError{
			Field:   "bet",
			Message: "bet cannot be nil",
			Value:   nil,
		}
	}

	// Get bet definition
	betDef, exists := crapsgame.CanonicalBetDefinitions[bet.Type]
	if !exists {
		return ValidationError{
			Field:   "betType",
			Message: fmt.Sprintf("unknown bet type: %s", bet.Type),
			Value:   bet.Type,
		}
	}

	// Check if bet requires specific numbers to be provided
	// Only certain bet types like PLACE_NUMBERS, HOP bets, etc. require numbers
	requiresNumbers := map[string]bool{
		"PLACE_NUMBERS": true,
		"HOP_1_2":       true,
		"HOP_1_3":       true,
		"HOP_1_4":       true,
		"HOP_1_5":       true,
		"HOP_1_6":       true,
		"HOP_2_3":       true,
		"HOP_2_4":       true,
		"HOP_2_5":       true,
		"HOP_2_6":       true,
		"HOP_3_4":       true,
		"HOP_3_5":       true,
		"HOP_3_6":       true,
		"HOP_4_5":       true,
		"HOP_4_6":       true,
		"HOP_5_6":       true,
	}

	if requiresNumbers[bet.Type] {
		// Bet requires specific numbers, validate they match
		if len(bet.Numbers) == 0 {
			return ValidationError{
				Field:   "numbers",
				Message: fmt.Sprintf("bet %s requires numbers but none provided", bet.Type),
				Value:   bet.Numbers,
			}
		}

		// Validate each number is in the valid range
		for _, num := range bet.Numbers {
			valid := false
			for _, validNum := range betDef.ValidNumbers {
				if num == validNum {
					valid = true
					break
				}
			}
			if !valid {
				return ValidationError{
					Field:   "numbers",
					Message: fmt.Sprintf("number %d is not valid for bet type %s", num, bet.Type),
					Value:   num,
				}
			}
		}
	} else {
		// Bet doesn't require specific numbers, but if numbers are provided, validate them
		for _, num := range bet.Numbers {
			if num < 1 || num > 12 {
				return ValidationError{
					Field:   "numbers",
					Message: fmt.Sprintf("number %d is outside valid range (1-12)", num),
					Value:   num,
				}
			}
		}
	}

	return nil
}

// validateBetModifiers validates that bet modifiers are valid
func validateBetModifiers(modifiers []*ModifierExpression) error {
	if modifiers == nil {
		return nil // No modifiers is valid
	}

	// Track used modifier types to detect conflicts
	usedModifiers := make(map[ModifierType]bool)

	for _, modifier := range modifiers {
		if modifier == nil {
			return ValidationError{
				Field:   "modifier",
				Message: "modifier cannot be nil",
				Value:   nil,
			}
		}

		// Check for conflicting modifiers
		switch modifier.Type {
		case ModWorking:
			if usedModifiers[ModOff] {
				return ValidationError{
					Field:   "modifiers",
					Message: "cannot have both WORKING and OFF modifiers",
					Value:   modifier.Type,
				}
			}
			usedModifiers[ModWorking] = true

		case ModOff:
			if usedModifiers[ModWorking] {
				return ValidationError{
					Field:   "modifiers",
					Message: "cannot have both WORKING and OFF modifiers",
					Value:   modifier.Type,
				}
			}
			usedModifiers[ModOff] = true

		case ModPress, ModOneRoll, ModMax, ModAmount:
			// These modifiers are valid and can be used together
			usedModifiers[modifier.Type] = true

		default:
			return ValidationError{
				Field:   "modifier",
				Message: fmt.Sprintf("unknown modifier type: %v", modifier.Type),
				Value:   modifier.Type,
			}
		}
	}

	return nil
}

// Error Recovery Functions for Phase 7.3

// recoverFromParseError attempts to recover from a parse error by skipping to the next statement
func recoverFromParseError(parser *Parser) Statement {
	// Log the error for debugging
	fmt.Printf("Parse error recovery: attempting to skip to next statement\n")

	// Skip tokens until we find a statement delimiter or new statement
	for parser.curToken.Type != SEMICOLON && parser.curToken.Type != EOF {
		parser.nextToken()

		// If we find a new statement keyword, try to parse it
		switch parser.curToken.Type {
		case PLACE, IF, SHOW, SET, DEFINE, EXECUTE, APPLY, REMOVE, PRESS, TURN, ROLL:
			return parser.parseStatement()
		}
	}

	// If we reach EOF, return nil
	if parser.curToken.Type == EOF {
		return nil
	}

	// Skip the semicolon and try to parse the next statement
	parser.nextToken()
	return parser.parseStatement()
}