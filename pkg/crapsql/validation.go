package crapsql

import (
	"fmt"

	"github.com/headswim/CrapsQL/pkg/crapsgame"
)

// ValidationError represents a validation error with field, message, and value
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
			Field:   "bet_amount",
			Message: "bet amount must be positive",
			Value:   amount,
		}
	}
	if amount < minBet {
		return ValidationError{
			Field:   "bet_amount",
			Message: fmt.Sprintf("bet amount $%.2f is below minimum $%.2f", amount, minBet),
			Value:   amount,
		}
	}
	if amount > maxBet {
		return ValidationError{
			Field:   "bet_amount",
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
			Message: fmt.Sprintf("insufficient bankroll: $%.2f available, $%.2f required", player.Bankroll, amount),
			Value:   player.Bankroll,
		}
	}
	return nil
}

// validateBetType validates that a bet type is valid and exists in canonical definitions
func validateBetType(betType string) error {
	if betType == "" {
		return ValidationError{
			Field:   "bet_type",
			Message: "bet type cannot be empty",
			Value:   betType,
		}
	}
	if !IsValidBetType(betType) {
		return ValidationError{
			Field:   "bet_type",
			Message: fmt.Sprintf("unknown bet type: %s", betType),
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

	// Get bet definition from canonical definitions
	betDef, exists := crapsgame.CanonicalBetDefinitions[betType]
	if !exists {
		// Try to get from bet registry as fallback
		if !IsValidBetType(betType) {
			return ValidationError{
				Field:   "bet_type",
				Message: fmt.Sprintf("unknown bet type: %s", betType),
				Value:   betType,
			}
		}
		// If it exists in registry but not canonical definitions, create a basic definition
		betDef = crapsgame.CanonicalBetDefinition{
			Name:              betType,
			Category:          "Unknown",
			Description:       "Bet type from registry",
			Payout:            "1:1",
			WorkingBehavior:   "ALWAYS",
			OneRoll:           false,
			PayoutNumerator:   1,
			PayoutDenominator: 1,
			ValidNumbers:      []int{},
			RequiresPoint:     false,
			RequiresComeOut:   false,
			HouseEdge:         0.0,
			Commission:        0.0,
		}
	}

	// Validate game state requirements
	if betDef.RequiresComeOut && gameState != crapsgame.StateComeOut {
		return ValidationError{
			Field:   "game_state",
			Message: fmt.Sprintf("bet type %s can only be placed during come-out phase", betType),
			Value:   gameState,
		}
	}

	if betDef.RequiresPoint && gameState != crapsgame.StatePoint {
		return ValidationError{
			Field:   "game_state",
			Message: fmt.Sprintf("bet type %s can only be placed during point phase", betType),
			Value:   gameState,
		}
	}

	// Validate game state is valid
	if gameState < crapsgame.StateComeOut || gameState > crapsgame.StateSevenOut {
		return ValidationError{
			Field:   "game_state",
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
			Message: "bet object is nil",
			Value:   nil,
		}
	}

	// Validate player is not nil
	if player == nil {
		return ValidationError{
			Field:   "player",
			Message: "player object is nil",
			Value:   nil,
		}
	}

	// Validate table is not nil
	if table == nil {
		return ValidationError{
			Field:   "table",
			Message: "table object is nil",
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
			Field:   "bet_id",
			Message: "bet must have an ID",
			Value:   bet.ID,
		}
	}

	// Validate bet player matches
	if bet.Player != player.ID {
		return ValidationError{
			Field:   "bet_player",
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
			Message: "bet is nil",
			Value:   nil,
		}
	}

	// Get bet definition from canonical definitions
	betDef, exists := crapsgame.CanonicalBetDefinitions[bet.Type]
	if !exists {
		// Try to get from bet registry as fallback
		if !IsValidBetType(bet.Type) {
			return ValidationError{
				Field:   "bet_type",
				Message: fmt.Sprintf("unknown bet type: %s", bet.Type),
				Value:   bet.Type,
			}
		}
		// If it exists in registry but not canonical definitions, create a basic definition
		betDef = crapsgame.CanonicalBetDefinition{
			Name:              bet.Type,
			Category:          "Unknown",
			Description:       "Bet type from registry",
			Payout:            "1:1",
			WorkingBehavior:   "ALWAYS",
			OneRoll:           false,
			PayoutNumerator:   1,
			PayoutDenominator: 1,
			ValidNumbers:      []int{},
			RequiresPoint:     false,
			RequiresComeOut:   false,
			HouseEdge:         0.0,
			Commission:        0.0,
		}
	}

	// Check if bet requires specific numbers
	if len(betDef.ValidNumbers) > 0 {
		// Bet requires specific numbers, validate they match
		if len(bet.Numbers) == 0 {
			return ValidationError{
				Field:   "bet_numbers",
				Message: fmt.Sprintf("bet type %s requires specific numbers", bet.Type),
				Value:   bet.Numbers,
			}
		}

		// Validate each number is in the valid range
		for _, num := range bet.Numbers {
			if num < 1 || num > 12 {
				return ValidationError{
					Field:   "bet_numbers",
					Message: fmt.Sprintf("invalid number %d for bet type %s (must be 1-12)", num, bet.Type),
					Value:   num,
				}
			}
		}
	} else {
		// Bet doesn't require specific numbers, but if numbers are provided, validate them
		if len(bet.Numbers) > 0 {
			for _, num := range bet.Numbers {
				if num < 1 || num > 12 {
					return ValidationError{
						Field:   "bet_numbers",
						Message: fmt.Sprintf("invalid number %d (must be 1-12)", num),
						Value:   num,
					}
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

	// Track modifier types to validate combinations
	modifierTypes := make(map[ModifierType]bool)

	for _, modifier := range modifiers {
		if modifier == nil {
			return ValidationError{
				Field:   "modifier",
				Message: "modifier is nil",
				Value:   nil,
			}
		}

		// Check for duplicate modifier types
		if modifierTypes[modifier.Type] {
			return ValidationError{
				Field:   "modifier_type",
				Message: fmt.Sprintf("duplicate modifier type: %v", modifier.Type),
				Value:   modifier.Type,
			}
		}
		modifierTypes[modifier.Type] = true

		// Validate modifier value if present
		if modifier.Value != nil {
			// For now, just check that the value is not nil
			// More specific validation could be added here based on modifier type
		}
	}

	return nil
}

// validateBetState validates the current state of a bet
func validateBetState(bet *crapsgame.Bet, table *crapsgame.Table) error {
	if bet == nil {
		return ValidationError{
			Field:   "bet",
			Message: "bet is nil",
			Value:   nil,
		}
	}

	if bet.Amount <= 0 {
		return ValidationError{
			Field:   "bet_amount",
			Message: fmt.Sprintf("bet amount must be positive: $%.2f", bet.Amount),
			Value:   bet.Amount,
		}
	}

	if bet.Player == "" {
		return ValidationError{
			Field:   "bet_player",
			Message: "bet must have a player ID",
			Value:   bet.Player,
		}
	}

	if bet.Type == "" {
		return ValidationError{
			Field:   "bet_type",
			Message: "bet must have a type",
			Value:   bet.Type,
		}
	}

	// Check if bet type exists in canonical definitions
	if _, exists := crapsgame.CanonicalBetDefinitions[bet.Type]; !exists {
		return ValidationError{
			Field:   "bet_type",
			Message: fmt.Sprintf("unknown bet type: %s", bet.Type),
			Value:   bet.Type,
		}
	}

	// Check if player exists
	if table != nil {
		if _, exists := table.Players[bet.Player]; !exists {
			return ValidationError{
				Field:   "bet_player",
				Message: fmt.Sprintf("player %s not found for bet", bet.Player),
				Value:   bet.Player,
			}
		}
	}

	return nil
}

// validateCommissionRate validates that a commission rate is valid
func validateCommissionRate(rate float64) error {
	if rate < 0 {
		return ValidationError{
			Field:   "commission_rate",
			Message: fmt.Sprintf("commission rate cannot be negative: %.2f", rate),
			Value:   rate,
		}
	}
	if rate >= 1 {
		return ValidationError{
			Field:   "commission_rate",
			Message: fmt.Sprintf("commission rate cannot be 100%% or greater: %.2f", rate),
			Value:   rate,
		}
	}
	return nil
}

// validateTableState validates the overall table state
func validateTableState(table *crapsgame.Table) error {
	if table == nil {
		return ValidationError{
			Field:   "table",
			Message: "table is nil",
			Value:   nil,
		}
	}

	// Validate game state
	if table.State < crapsgame.StateComeOut || table.State > crapsgame.StateSevenOut {
		return ValidationError{
			Field:   "game_state",
			Message: fmt.Sprintf("invalid game state: %d", table.State),
			Value:   table.State,
		}
	}

	// Validate point
	if err := validatePoint(table.Point); err != nil {
		return err
	}

	// Validate shooter if we have players
	if len(table.Players) > 0 {
		if err := validateShooter(table.Shooter, table); err != nil {
			return err
		}
	}

	// Validate that point is only set during point phase
	if table.State == crapsgame.StateComeOut && table.Point != crapsgame.PointOff {
		return ValidationError{
			Field:   "point",
			Message: "point should be off during come out phase",
			Value:   table.Point,
		}
	}

	if table.State == crapsgame.StatePoint && table.Point == crapsgame.PointOff {
		return ValidationError{
			Field:   "point",
			Message: "point should be set during point phase",
			Value:   table.Point,
		}
	}

	return nil
}

// validatePoint validates that a point number is valid
func validatePoint(point crapsgame.Point) error {
	switch point {
	case crapsgame.PointOff, crapsgame.Point4, crapsgame.Point5, crapsgame.Point6, crapsgame.Point8, crapsgame.Point9, crapsgame.Point10:
		return nil
	default:
		return ValidationError{
			Field:   "point",
			Message: fmt.Sprintf("invalid point: %d", point),
			Value:   point,
		}
	}
}

// validateShooter validates that the shooter exists and is valid
func validateShooter(shooterID string, table *crapsgame.Table) error {
	if shooterID == "" {
		return ValidationError{
			Field:   "shooter",
			Message: "no shooter assigned",
			Value:   shooterID,
		}
	}

	if table == nil {
		return ValidationError{
			Field:   "table",
			Message: "table is nil for shooter validation",
			Value:   nil,
		}
	}

	if _, exists := table.Players[shooterID]; !exists {
		return ValidationError{
			Field:   "shooter",
			Message: fmt.Sprintf("shooter %s not found", shooterID),
			Value:   shooterID,
		}
	}

	return nil
}

// validateStateTransition validates if a state transition is valid
func validateStateTransition(fromState crapsgame.GameState, toState crapsgame.GameState, roll *crapsgame.Roll, currentPoint crapsgame.Point) error {
	// Validate that we're not in an invalid state
	if fromState < crapsgame.StateComeOut || fromState > crapsgame.StateSevenOut {
		return ValidationError{
			Field:   "from_state",
			Message: fmt.Sprintf("invalid current state: %d", fromState),
			Value:   fromState,
		}
	}

	if toState < crapsgame.StateComeOut || toState > crapsgame.StateSevenOut {
		return ValidationError{
			Field:   "to_state",
			Message: fmt.Sprintf("invalid target state: %d", toState),
			Value:   toState,
		}
	}

	// Validate specific transitions
	switch fromState {
	case crapsgame.StateComeOut:
		switch toState {
		case crapsgame.StatePoint:
			// Valid: point establishment
			if roll == nil {
				return ValidationError{
					Field:   "roll",
					Message: "roll is nil for point establishment",
					Value:   nil,
				}
			}
			if roll.Total < 4 || roll.Total > 10 || roll.Total == 7 {
				return ValidationError{
					Field:   "roll_total",
					Message: fmt.Sprintf("invalid point number: %d", roll.Total),
					Value:   roll.Total,
				}
			}
		case crapsgame.StateComeOut:
			// Valid: natural or craps
			if roll == nil {
				return ValidationError{
					Field:   "roll",
					Message: "roll is nil for come out",
					Value:   nil,
				}
			}
			if roll.Total != 2 && roll.Total != 3 && roll.Total != 7 && roll.Total != 11 && roll.Total != 12 {
				return ValidationError{
					Field:   "roll_total",
					Message: fmt.Sprintf("invalid come out roll: %d", roll.Total),
					Value:   roll.Total,
				}
			}
		default:
			return ValidationError{
				Field:   "state_transition",
				Message: fmt.Sprintf("invalid transition from come out to state %d", toState),
				Value:   toState,
			}
		}
	case crapsgame.StatePoint:
		switch toState {
		case crapsgame.StateComeOut:
			// Valid: point resolution
			if roll == nil {
				return ValidationError{
					Field:   "roll",
					Message: "roll is nil for point resolution",
					Value:   nil,
				}
			}
			pointNumber, err := crapsgame.PointToNumber(currentPoint)
			if err != nil {
				return ValidationError{
					Field:   "point",
					Message: fmt.Sprintf("invalid point state: %v", err),
					Value:   currentPoint,
				}
			}
			if roll.Total != pointNumber && roll.Total != 7 {
				return ValidationError{
					Field:   "roll_total",
					Message: fmt.Sprintf("invalid point phase roll: %d", roll.Total),
					Value:   roll.Total,
				}
			}
		case crapsgame.StateSevenOut:
			// Valid: seven out
			if roll == nil {
				return ValidationError{
					Field:   "roll",
					Message: "roll is nil for seven out",
					Value:   nil,
				}
			}
			if roll.Total != 7 {
				return ValidationError{
					Field:   "roll_total",
					Message: fmt.Sprintf("invalid seven out roll: %d", roll.Total),
					Value:   roll.Total,
				}
			}
		default:
			return ValidationError{
				Field:   "state_transition",
				Message: fmt.Sprintf("invalid transition from point to state %d", toState),
				Value:   toState,
			}
		}
	case crapsgame.StateSevenOut:
		// Seven out should only transition to come out
		if toState != crapsgame.StateComeOut {
			return ValidationError{
				Field:   "state_transition",
				Message: fmt.Sprintf("invalid transition from seven out to state %d", toState),
				Value:   toState,
			}
		}
	}

	return nil
}

// validatePlayer validates that a player is valid
func validatePlayer(player *crapsgame.Player) error {
	if player == nil {
		return ValidationError{
			Field:   "player",
			Message: "player is nil",
			Value:   nil,
		}
	}

	if player.ID == "" {
		return ValidationError{
			Field:   "player_id",
			Message: "player ID cannot be empty",
			Value:   player.ID,
		}
	}

	if player.Name == "" {
		return ValidationError{
			Field:   "player_name",
			Message: "player name cannot be empty",
			Value:   player.Name,
		}
	}

	if player.Bankroll < 0 {
		return ValidationError{
			Field:   "player_bankroll",
			Message: fmt.Sprintf("player bankroll cannot be negative: $%.2f", player.Bankroll),
			Value:   player.Bankroll,
		}
	}

	return nil
}

// validateTable validates that a table is valid
func validateTable(table *crapsgame.Table) error {
	if table == nil {
		return ValidationError{
			Field:   "table",
			Message: "table is nil",
			Value:   nil,
		}
	}

	if table.MinBet <= 0 {
		return ValidationError{
			Field:   "min_bet",
			Message: fmt.Sprintf("minimum bet must be positive: $%.2f", table.MinBet),
			Value:   table.MinBet,
		}
	}

	if table.MaxBet <= 0 {
		return ValidationError{
			Field:   "max_bet",
			Message: fmt.Sprintf("maximum bet must be positive: $%.2f", table.MaxBet),
			Value:   table.MaxBet,
		}
	}

	if table.MinBet > table.MaxBet {
		return ValidationError{
			Field:   "bet_limits",
			Message: fmt.Sprintf("minimum bet $%.2f cannot be greater than maximum bet $%.2f", table.MinBet, table.MaxBet),
			Value:   fmt.Sprintf("min: $%.2f, max: $%.2f", table.MinBet, table.MaxBet),
		}
	}

	if table.MaxOdds < 0 {
		return ValidationError{
			Field:   "max_odds",
			Message: fmt.Sprintf("maximum odds cannot be negative: %d", table.MaxOdds),
			Value:   table.MaxOdds,
		}
	}

	return nil
}

// recoverFromParseError attempts to recover from parse errors by skipping to next statement
func recoverFromParseError(parser *Parser) Statement {
	// This is a simplified recovery mechanism
	// In a more robust implementation, you might want to:
	// 1. Skip tokens until you find a semicolon or statement boundary
	// 2. Log the error for debugging
	// 3. Return a special error statement or nil

	fmt.Println("Parse error recovery: attempting to skip to next statement")

	// For now, just return nil to indicate we couldn't recover
	return nil
}
