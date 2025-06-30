package crapsql

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/headswim/CrapsQL/pkg/crapsgame"
)

// Interpreter executes CrapsQL statements
type Interpreter struct {
	table   *crapsgame.Table
	results []string
}

// NewInterpreter creates a new interpreter
func NewInterpreter(table *crapsgame.Table) *Interpreter {
	return &Interpreter{
		table: table,
	}
}

// Execute executes a CrapsQL program
func (i *Interpreter) Execute(program *Program) ([]string, error) {
	var results []string

	for _, stmt := range program.Statements {
		result, err := i.executeStatement(stmt)
		if err != nil {
			return results, err
		}
		if result != "" {
			results = append(results, result)
		}
	}

	return results, nil
}

// ExecuteString parses and executes a CrapsQL string
func (i *Interpreter) ExecuteString(input string) ([]string, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(parser.Errors(), "; "))
	}

	return i.Execute(program)
}

// ExecuteStringForPlayer parses and executes a CrapsQL string for a specific player
func (i *Interpreter) ExecuteStringForPlayer(input string, playerID string) ([]string, error) {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(parser.Errors(), "; "))
	}

	return i.ExecuteForPlayer(program, playerID)
}

// ExecuteForPlayer executes a CrapsQL program for a specific player
func (i *Interpreter) ExecuteForPlayer(program *Program, playerID string) ([]string, error) {
	var results []string

	for _, stmt := range program.Statements {
		result, err := i.executeStatementForPlayer(stmt, playerID)
		if err != nil {
			return results, err
		}
		if result != "" {
			results = append(results, result)
		}
	}

	return results, nil
}

func (i *Interpreter) executeStatement(stmt Statement) (string, error) {
	switch s := stmt.(type) {
	case *BetStatement:
		return i.executeBetStatement(s)
	case *ConditionalStatement:
		return i.executeConditionalStatement(s)
	case *QueryStatement:
		return i.executeQueryStatement(s)
	case *ManagementStatement:
		return i.executeManagementStatement(s)
	case *RemoveStatement:
		return i.executeRemoveStatement(s)
	case *PressStatement:
		return i.executePressStatement(s)
	case *TurnStatement:
		return i.executeTurnStatement(s)
	case *RollStatement:
		return i.executeRollStatement(s)
	default:
		return "", fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (i *Interpreter) executeStatementForPlayer(stmt Statement, playerID string) (string, error) {
	switch s := stmt.(type) {
	case *BetStatement:
		return i.executeBetStatementForPlayer(s, playerID)
	case *ConditionalStatement:
		return i.executeConditionalStatementForPlayer(s, playerID)
	case *QueryStatement:
		return i.executeQueryStatementForPlayer(s, playerID)
	case *ManagementStatement:
		return i.executeManagementStatementForPlayer(s, playerID)
	case *RemoveStatement:
		return i.executeRemoveStatementForPlayer(s, playerID)
	case *PressStatement:
		return i.executePressStatementForPlayer(s, playerID)
	case *TurnStatement:
		return i.executeTurnStatementForPlayer(s, playerID)
	case *RollStatement:
		return i.executeRollStatementForPlayer(s, playerID)
	default:
		return "", fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (i *Interpreter) executeBetStatement(stmt *BetStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	numbers := extractNumbersForBetType(stmt.BetType)

	// Create bet object for validation
	bet := &crapsgame.Bet{
		ID:      generateBetID(),
		Type:    betType,
		Amount:  stmt.Amount.Value,
		Player:  playerID,
		Working: true,
		Numbers: numbers,
	}

	// Perform comprehensive validation
	if err := validateBetPlacement(bet, player, i.table); err != nil {
		return "", fmt.Errorf("bet validation failed: %v", err)
	}

	// Validate bet modifiers if present
	if len(stmt.Modifiers) > 0 {
		if err := validateBetModifiers(stmt.Modifiers); err != nil {
			return "", fmt.Errorf("modifier validation failed: %v", err)
		}
	}

	// Place the bet
	placedBet, err := i.table.PlaceBet(playerID, betType, stmt.Amount.Value, numbers)
	if err != nil {
		return "", fmt.Errorf("failed to place bet: %v", err)
	}

	return fmt.Sprintf("✅ Placed $%.2f on %s", placedBet.Amount, betType), nil
}

func (i *Interpreter) executeBetStatementForPlayer(stmt *BetStatement, playerID string) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	numbers := extractNumbersForBetType(stmt.BetType)

	// Create bet object for validation
	bet := &crapsgame.Bet{
		ID:      generateBetID(),
		Type:    betType,
		Amount:  stmt.Amount.Value,
		Player:  playerID,
		Working: true,
		Numbers: numbers,
	}

	// Perform comprehensive validation
	if err := validateBetPlacement(bet, player, i.table); err != nil {
		return "", fmt.Errorf("bet validation failed: %v", err)
	}

	// Validate bet modifiers if present
	if len(stmt.Modifiers) > 0 {
		if err := validateBetModifiers(stmt.Modifiers); err != nil {
			return "", fmt.Errorf("modifier validation failed: %v", err)
		}
	}

	// Place the bet
	placedBet, err := i.table.PlaceBet(playerID, betType, stmt.Amount.Value, numbers)
	if err != nil {
		return "", fmt.Errorf("failed to place bet: %v", err)
	}

	return fmt.Sprintf("✅ Placed $%.2f on %s", placedBet.Amount, betType), nil
}

// generateBetID generates a unique bet ID
func generateBetID() string {
	// Generate a random 8-character alphanumeric ID
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to timestamp if crypto/rand fails
			return fmt.Sprintf("bet_%d", time.Now().UnixNano())
		}
		result[i] = charset[n.Int64()]
	}

	return fmt.Sprintf("bet_%s", string(result))
}

// extractNumbersForBetType extracts numbers for place numbers bets
func extractNumbersForBetType(expr *BetTypeExpression) []int {
	var numbers []int

	switch expr.Type {
	case BetPlaceNumbers:
		// Extract place numbers from args
		for _, arg := range expr.Args {
			if numExpr, ok := arg.(*NumberExpression); ok {
				numbers = append(numbers, int(numExpr.Value))
			} else if amtExpr, ok := arg.(*AmountExpression); ok {
				numbers = append(numbers, int(amtExpr.Value))
			}
		}

	case BetHop, BetHop12, BetHop13, BetHop14, BetHop15, BetHop16,
		BetHop23, BetHop24, BetHop25, BetHop26, BetHop34, BetHop35, BetHop36,
		BetHop45, BetHop46, BetHop56:
		// Extract die values from hop bet args
		for _, arg := range expr.Args {
			if numExpr, ok := arg.(*NumberExpression); ok {
				numbers = append(numbers, int(numExpr.Value))
			} else if amtExpr, ok := arg.(*AmountExpression); ok {
				numbers = append(numbers, int(amtExpr.Value))
			}
		}

	case BetPassOdds, BetDontPassOdds, BetComeOdds, BetDontComeOdds:
		// Odds bets don't need specific numbers - they use the point
		// The point will be extracted from game state when placing the bet
		return nil

	case BetPlace4:
		numbers = append(numbers, 4)
	case BetPlace5:
		numbers = append(numbers, 5)
	case BetPlace6:
		numbers = append(numbers, 6)
	case BetPlace8:
		numbers = append(numbers, 8)
	case BetPlace9:
		numbers = append(numbers, 9)
	case BetPlace10:
		numbers = append(numbers, 10)

	case BetBuy4, BetLay4, BetPlaceToLose4:
		numbers = append(numbers, 4)
	case BetBuy5, BetLay5, BetPlaceToLose5:
		numbers = append(numbers, 5)
	case BetBuy6, BetLay6, BetPlaceToLose6:
		numbers = append(numbers, 6)
	case BetBuy8, BetLay8, BetPlaceToLose8:
		numbers = append(numbers, 8)
	case BetBuy9, BetLay9, BetPlaceToLose9:
		numbers = append(numbers, 9)
	case BetBuy10, BetLay10, BetPlaceToLose10:
		numbers = append(numbers, 10)

	case BetHard4, BetHard6, BetHard8, BetHard10:
		// Hard way bets use the number in their name
		switch expr.Type {
		case BetHard4:
			numbers = append(numbers, 4)
		case BetHard6:
			numbers = append(numbers, 6)
		case BetHard8:
			numbers = append(numbers, 8)
		case BetHard10:
			numbers = append(numbers, 10)
		}

	case BetHornHigh2, BetHornHigh3, BetHornHigh11, BetHornHigh12:
		// Horn high bets use the number in their name
		switch expr.Type {
		case BetHornHigh2:
			numbers = append(numbers, 2)
		case BetHornHigh3:
			numbers = append(numbers, 3)
		case BetHornHigh11:
			numbers = append(numbers, 11)
		case BetHornHigh12:
			numbers = append(numbers, 12)
		}

	case BetBig6, BetBig8:
		// Big bets use the number in their name
		switch expr.Type {
		case BetBig6:
			numbers = append(numbers, 6)
		case BetBig8:
			numbers = append(numbers, 8)
		}

	default:
		// For other bet types that don't require specific numbers
		return nil
	}

	// Validate extracted numbers
	for _, num := range numbers {
		if num < 1 || num > 12 {
			// Invalid number range for craps
			return nil
		}
	}

	return numbers
}

func (i *Interpreter) executeConditionalStatement(stmt *ConditionalStatement) (string, error) {
	// Evaluate condition
	condition, err := i.evaluateCondition(stmt.Condition)
	if err != nil {
		return "", err
	}

	var block *BlockStatement
	if condition {
		block = stmt.Consequence
	} else if stmt.Alternative != nil {
		block = stmt.Alternative
	} else {
		return "", nil // No action needed
	}

	// Execute the block
	var results []string
	for _, s := range block.Statements {
		result, err := i.executeStatement(s)
		if err != nil {
			return "", err
		}
		if result != "" {
			results = append(results, result)
		}
	}

	return strings.Join(results, "\n"), nil
}

func (i *Interpreter) executeConditionalStatementForPlayer(stmt *ConditionalStatement, playerID string) (string, error) {
	condition, err := i.evaluateConditionForPlayer(stmt.Condition, playerID)
	if err != nil {
		return "", err
	}

	if condition {
		// Execute consequence
		var results []string
		for _, s := range stmt.Consequence.Statements {
			result, err := i.executeStatementForPlayer(s, playerID)
			if err != nil {
				return "", err
			}
			if result != "" {
				results = append(results, result)
			}
		}
		return strings.Join(results, "\n"), nil
	} else if stmt.Alternative != nil {
		// Execute alternative
		var results []string
		for _, s := range stmt.Alternative.Statements {
			result, err := i.executeStatementForPlayer(s, playerID)
			if err != nil {
				return "", err
			}
			if result != "" {
				results = append(results, result)
			}
		}
		return strings.Join(results, "\n"), nil
	}

	return "", nil
}

func (i *Interpreter) executeQueryStatement(stmt *QueryStatement) (string, error) {
	switch stmt.Type {
	case QueryPoint:
		return i.executeShowPoint(), nil
	case QueryBets:
		return i.executeShowBets(), nil
	case QueryBankroll:
		var playerID string
		for id := range i.table.Players {
			playerID = id
			break
		}
		if playerID == "" {
			return "", fmt.Errorf("no players at table - add a player first")
		}
		return i.executeShowBankroll(playerID), nil
	case QueryTableMinimums:
		return i.executeShowTableMinimums(), nil
	case QueryOddsAllowed:
		return fmt.Sprintf("Max odds: %dx", i.table.MaxOdds), nil
	default:
		return "", fmt.Errorf("unknown query type: %d", stmt.Type)
	}
}

func (i *Interpreter) executeQueryStatementForPlayer(stmt *QueryStatement, playerID string) (string, error) {
	switch stmt.Type {
	case QueryPoint:
		return i.executeShowPoint(), nil
	case QueryBets:
		return i.executeShowBets(), nil
	case QueryBankroll:
		return i.executeShowBankroll(playerID), nil
	case QueryTableMinimums:
		return i.executeShowTableMinimums(), nil
	case QueryOddsAllowed:
		return fmt.Sprintf("Max odds: %dx", i.table.MaxOdds), nil
	default:
		return "", fmt.Errorf("unknown query type: %d", stmt.Type)
	}
}

func (i *Interpreter) executeManagementStatement(stmt *ManagementStatement) (string, error) {
	value, err := i.extractAmountFromExpression(stmt.Value)
	if err != nil {
		return "", fmt.Errorf("invalid amount: %v", err)
	}

	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	switch stmt.Type {
	case ManageBankroll:
		return i.executeSetBankroll(playerID, value)
	case ManageMaxBet:
		return i.executeSetMaxBet(playerID, value)
	case ManageMinBet:
		return i.executeSetMinBet(playerID, value)
	case ManageWinGoal:
		return i.executeSetWinGoal(playerID, value)
	case ManageLossLimit:
		return i.executeSetLossLimit(playerID, value)
	case ManageSessionTime:
		return fmt.Sprintf("Session time limit set to %.0f hours", value), nil
	default:
		return "", fmt.Errorf("unknown management type: %d", stmt.Type)
	}
}

func (i *Interpreter) executeManagementStatementForPlayer(stmt *ManagementStatement, playerID string) (string, error) {
	value, err := i.extractAmountFromExpression(stmt.Value)
	if err != nil {
		return "", fmt.Errorf("invalid amount: %v", err)
	}

	switch stmt.Type {
	case ManageBankroll:
		return i.executeSetBankroll(playerID, value)
	case ManageMaxBet:
		return i.executeSetMaxBet(playerID, value)
	case ManageMinBet:
		return i.executeSetMinBet(playerID, value)
	case ManageWinGoal:
		return i.executeSetWinGoal(playerID, value)
	case ManageLossLimit:
		return i.executeSetLossLimit(playerID, value)
	case ManageSessionTime:
		return fmt.Sprintf("Session time limit set to %.0f hours", value), nil
	default:
		return "", fmt.Errorf("unknown management type: %d", stmt.Type)
	}
}

func (i *Interpreter) executeSetBankroll(playerID string, amount float64) (string, error) {
	if amount < 0 {
		return "", fmt.Errorf("bankroll cannot be negative")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found: %v", playerID, err)
	}

	player.Bankroll = amount
	return fmt.Sprintf("Player %s bankroll set to $%.2f", playerID, amount), nil
}

func (i *Interpreter) executeSetMaxBet(playerID string, amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("max bet must be positive")
	}

	if amount < i.table.MinBet {
		return "", fmt.Errorf("max bet cannot be less than table minimum ($%.2f)", i.table.MinBet)
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found: %v", playerID, err)
	}

	player.MaxBet = amount
	return fmt.Sprintf("Player %s max bet set to $%.2f", playerID, amount), nil
}

func (i *Interpreter) executeSetMinBet(playerID string, amount float64) (string, error) {
	if amount < 0 {
		return "", fmt.Errorf("min bet cannot be negative")
	}

	if amount > i.table.MaxBet {
		return "", fmt.Errorf("min bet cannot be greater than table maximum ($%.2f)", i.table.MaxBet)
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found: %v", playerID, err)
	}

	player.MinBet = amount
	return fmt.Sprintf("Player %s min bet set to $%.2f", playerID, amount), nil
}

func (i *Interpreter) executeSetWinGoal(playerID string, amount float64) (string, error) {
	if amount < 0 {
		return "", fmt.Errorf("win goal cannot be negative")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found: %v", playerID, err)
	}

	player.WinGoal = amount
	return fmt.Sprintf("Player %s win goal set to $%.2f", playerID, amount), nil
}

func (i *Interpreter) executeSetLossLimit(playerID string, amount float64) (string, error) {
	if amount < 0 {
		return "", fmt.Errorf("loss limit cannot be negative")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found: %v", playerID, err)
	}

	player.LossLimit = amount
	return fmt.Sprintf("Player %s loss limit set to $%.2f", playerID, amount), nil
}

func (i *Interpreter) extractAmountFromExpression(expr Expression) (float64, error) {
	switch e := expr.(type) {
	case *AmountExpression:
		return e.Value, nil
	case *NumberExpression:
		return e.Value, nil
	default:
		return 0, fmt.Errorf("unsupported expression type for amount: %T", expr)
	}
}

func (i *Interpreter) executeRemoveStatement(stmt *RemoveStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", err
	}

	if stmt.BetType == nil {
		// REMOVE ALL
		removedCount := 0
		var remainingBets []*crapsgame.Bet
		for _, bet := range player.Bets {
			if bet.Working {
				player.Bankroll += bet.Amount
				removedCount++
			} else {
				remainingBets = append(remainingBets, bet)
			}
		}
		player.Bets = remainingBets
		return fmt.Sprintf("✅ Removed %d bets", removedCount), nil
	}

	// Remove specific bet type
	betType := i.betTypeToString(stmt.BetType.Type)
	removedCount := 0
	var remainingBets []*crapsgame.Bet
	for _, bet := range player.Bets {
		if bet.Type == betType && bet.Working {
			player.Bankroll += bet.Amount
			removedCount++
		} else {
			remainingBets = append(remainingBets, bet)
		}
	}
	player.Bets = remainingBets
	return fmt.Sprintf("✅ Removed %d %s bets", removedCount, betType), nil
}

func (i *Interpreter) executeRemoveStatementForPlayer(stmt *RemoveStatement, playerID string) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	if stmt.BetType == nil {
		// Remove all bets for the player
		removedCount := 0
		var remainingBets []*crapsgame.Bet
		for _, bet := range player.Bets {
			if bet.Working {
				player.Bankroll += bet.Amount
				removedCount++
			} else {
				remainingBets = append(remainingBets, bet)
			}
		}
		player.Bets = remainingBets
		return fmt.Sprintf("Removed %d bets for player %s", removedCount, playerID), nil
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	removedCount := 0
	var remainingBets []*crapsgame.Bet
	for _, bet := range player.Bets {
		if bet.Type == betType && bet.Working {
			player.Bankroll += bet.Amount
			removedCount++
		} else {
			remainingBets = append(remainingBets, bet)
		}
	}
	player.Bets = remainingBets
	return fmt.Sprintf("Removed %d %s bets for player %s", removedCount, betType, playerID), nil
}

func (i *Interpreter) executePressStatement(stmt *PressStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", err
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	pressedCount := 0

	for _, bet := range player.Bets {
		if bet.Type == betType && bet.Working {
			// Add the press amount to the bet
			bet.Amount += stmt.Amount.Value
			player.Bankroll -= stmt.Amount.Value // Deduct from bankroll
			pressedCount++
		}
	}

	if pressedCount == 0 {
		return "", fmt.Errorf("no active %s bets to press", betType)
	}

	return fmt.Sprintf("✅ Pressed %d %s bets by $%.2f", pressedCount, betType, stmt.Amount.Value), nil
}

func (i *Interpreter) executePressStatementForPlayer(stmt *PressStatement, playerID string) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	amount := stmt.Amount.Value
	pressedCount := 0

	for _, bet := range player.Bets {
		if bet.Type == betType && bet.Working {
			// Add the press amount to the bet
			bet.Amount += amount
			player.Bankroll -= amount // Deduct from bankroll
			pressedCount++
		}
	}

	if pressedCount == 0 {
		return "", fmt.Errorf("no active %s bets to press for player %s", betType, playerID)
	}

	return fmt.Sprintf("Pressed %d %s bets by $%.2f for player %s", pressedCount, betType, amount, playerID), nil
}

func (i *Interpreter) executeTurnStatement(stmt *TurnStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", err
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	turnedCount := 0

	for _, bet := range player.Bets {
		if bet.Type == betType {
			bet.Working = (stmt.Action == "ON")
			turnedCount++
		}
	}

	if turnedCount == 0 {
		return "", fmt.Errorf("no %s bets to turn %s", betType, stmt.Action)
	}

	return fmt.Sprintf("✅ Turned %s %d %s bets", stmt.Action, turnedCount, betType), nil
}

func (i *Interpreter) executeTurnStatementForPlayer(stmt *TurnStatement, playerID string) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	betType := i.betTypeToString(stmt.BetType.Type)
	action := stmt.Action
	turnedCount := 0

	for _, bet := range player.Bets {
		if bet.Type == betType {
			bet.Working = (action == "ON")
			turnedCount++
		}
	}

	if turnedCount == 0 {
		return "", fmt.Errorf("no %s bets to turn %s for player %s", betType, action, playerID)
	}

	return fmt.Sprintf("Turned %s %d %s bets for player %s", action, turnedCount, betType, playerID), nil
}

func (i *Interpreter) evaluateConditionForPlayer(expr Expression, playerID string) (bool, error) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		return i.evaluateIdentifierConditionForPlayer(e, playerID)
	case *InfixExpression:
		return i.evaluateInfixConditionForPlayer(e, playerID)
	default:
		return false, fmt.Errorf("unsupported condition type: %T", expr)
	}
}

func (i *Interpreter) evaluateIdentifierConditionForPlayer(expr *IdentifierExpression, playerID string) (bool, error) {
	value, err := i.evaluateIdentifierExpressionForPlayer(expr, playerID)
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

func (i *Interpreter) evaluateInfixConditionForPlayer(expr *InfixExpression, playerID string) (bool, error) {
	left, err := i.evaluateExpressionForPlayer(expr.Left, playerID)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpressionForPlayer(expr.Right, playerID)
	if err != nil {
		return false, err
	}

	switch expr.Operator {
	case ">":
		return left > right, nil
	case "<":
		return left < right, nil
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", expr.Operator)
	}
}

func (i *Interpreter) evaluateExpressionForPlayer(expr Expression, playerID string) (float64, error) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		return i.evaluateIdentifierExpressionForPlayer(e, playerID)
	case *NumberExpression:
		return e.Value, nil
	case *AmountExpression:
		return e.Value, nil
	default:
		return 0, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// PlaceBet places a bet on the table
func (i *Interpreter) PlaceBet(playerID, betType string, amount float64) error {
	// Validate bet type against our definitions
	if _, exists := crapsgame.CanonicalBetDefinitions[betType]; !exists {
		// Check if it's a legacy bet type
		legacyTypes := []string{"PASS_LINE", "DONT_PASS", "COME", "DONT_COME", "FIELD", "ANY_SEVEN",
			"ANY_CRAPS", "ELEVEN", "ACE_DEUCE", "ACES", "BOXCARS", "PLACE_4",
			"PLACE_5", "PLACE_6", "PLACE_8", "PLACE_9", "PLACE_10", "HARD_4",
			"HARD_6", "HARD_8", "HARD_10", "PASS_ODDS", "DONT_PASS_ODDS", "BUY_4",
			"BUY_10", "LAY_4", "LAY_10", "BIG_6", "BIG_8", "HOP", "WORLD",
			"C_AND_E", "HORN"}

		found := false
		for _, legacy := range legacyTypes {
			if betType == legacy {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unsupported bet type: %s. Use SHOW BETS to see available bet types", betType)
		}
	}

	_, err := i.table.PlaceBet(playerID, betType, amount, []int{})
	if err != nil {
		return fmt.Errorf("failed to place bet: %v", err)
	}

	i.results = append(i.results, fmt.Sprintf("Placed $%.2f on %s", amount, betType))
	return nil
}

// ShowBets lists all available bet types with their descriptions
func (i *Interpreter) ShowBets() {
	i.results = append(i.results, "Available bet types:")

	// Group bets by category
	categories := make(map[string][]string)
	for betType, betDef := range crapsgame.CanonicalBetDefinitions {
		categoryStr := string(betDef.Category)
		categories[categoryStr] = append(categories[categoryStr], betType)
	}

	for category, betTypes := range categories {
		i.results = append(i.results, fmt.Sprintf("\n%s:", category))
		for _, betType := range betTypes {
			betDef := crapsgame.CanonicalBetDefinitions[betType]
			i.results = append(i.results, fmt.Sprintf("  %s - %s (Payout: %s)",
				betType, betDef.Description, betDef.Payout))
		}
	}
}

func (i *Interpreter) executeRollStatement(stmt *RollStatement) (string, error) {
	// Roll the dice and resolve all bets
	roll := i.table.RollDice()

	// Get the results from bet resolution
	results := i.table.BetResolver.ResolveBets(roll)

	// Format the output
	var output strings.Builder
	output.WriteString(fmt.Sprintf("🎲 Rolled %d (%d + %d)", roll.Total, roll.Die1, roll.Die2))

	if len(results) > 0 {
		output.WriteString("\n")
		output.WriteString(strings.Join(results, "\n"))
	}

	return output.String(), nil
}

func (i *Interpreter) executeRollStatementForPlayer(stmt *RollStatement, playerID string) (string, error) {
	// For player-specific rolls, we still roll for the whole table
	// but we can filter results for the specific player
	roll := i.table.RollDice()

	// Get all results
	allResults := i.table.BetResolver.ResolveBets(roll)

	// Filter results for this player
	var playerResults []string
	for _, result := range allResults {
		// Check if this result is for our player
		// This is a simple check - in a real implementation you'd want more sophisticated filtering
		if strings.Contains(result, playerID) || !strings.Contains(result, "player") {
			playerResults = append(playerResults, result)
		}
	}

	// Format the output
	var output strings.Builder
	output.WriteString(fmt.Sprintf("🎲 Rolled %d (%d + %d)", roll.Total, roll.Die1, roll.Die2))

	if len(playerResults) > 0 {
		output.WriteString("\n")
		output.WriteString(strings.Join(playerResults, "\n"))
	}

	return output.String(), nil
}

func (i *Interpreter) executeShowPoint() string {
	pointNumber := i.table.GetPointNumber()
	if pointNumber == 0 {
		return "Point: OFF"
	} else {
		return fmt.Sprintf("Point: %d", pointNumber)
	}
}

func (i *Interpreter) executeShowBets() string {
	var output strings.Builder
	output.WriteString("=== AVAILABLE BET TYPES ===\n\n")

	// Group bets by category
	categories := crapsgame.GetBetsByCategory()

	// Define category display order
	categoryOrder := []crapsgame.BetCategory{
		crapsgame.LineBets,
		crapsgame.ComeBets,
		crapsgame.OddsBets,
		crapsgame.PlaceBets,
		crapsgame.BuyBets,
		crapsgame.LayBets,
		crapsgame.PlaceToLoseBets,
		crapsgame.FieldBets,
		crapsgame.HardWayBets,
		crapsgame.HopBets,
		crapsgame.HornBets,
		crapsgame.PropositionBets,
	}

	for _, category := range categoryOrder {
		betTypes := categories[category]
		if len(betTypes) == 0 {
			continue
		}

		// Category header
		output.WriteString(fmt.Sprintf("=== %s (%d bets) ===\n", string(category), len(betTypes)))

		for _, betType := range betTypes {
			betDef := crapsgame.CanonicalBetDefinitions[betType]
			output.WriteString(fmt.Sprintf("  %s\n", betType))
			output.WriteString(fmt.Sprintf("    Description: %s\n", betDef.Description))
			output.WriteString(fmt.Sprintf("    Payout: %s\n", betDef.Payout))
			output.WriteString(fmt.Sprintf("    Working: %s\n", betDef.WorkingBehavior))
			if betDef.OneRoll {
				output.WriteString("    Type: One-roll bet\n")
			} else {
				output.WriteString("    Type: Multi-roll bet\n")
			}
			output.WriteString(fmt.Sprintf("    House Edge: %.2f%%\n", betDef.HouseEdge))
			if betDef.Commission > 0 {
				output.WriteString(fmt.Sprintf("    Commission: %.1f%%\n", betDef.Commission*100))
			}
			output.WriteString("\n")
		}
	}

	return output.String()
}

func (i *Interpreter) executeShowBankroll(playerID string) string {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return fmt.Sprintf("Error: Player %s not found", playerID)
	}
	return fmt.Sprintf("Player %s Bankroll: $%.2f", playerID, player.Bankroll)
}

func (i *Interpreter) executeShowTableMinimums() string {
	return fmt.Sprintf("Table Limits:\n  Minimum Bet: $%.2f\n  Maximum Bet: $%.2f\n  Maximum Odds: %dx",
		i.table.MinBet, i.table.MaxBet, i.table.MaxOdds)
}

func (i *Interpreter) betTypeToString(betType BetType) string {
	switch betType {
	case BetPassLine:
		return "PASS_LINE"
	case BetDontPass:
		return "DONT_PASS"
	case BetCome:
		return "COME"
	case BetDontCome:
		return "DONT_COME"
	case BetField:
		return "FIELD"
	case BetAnySeven:
		return "ANY_SEVEN"
	case BetAnyCraps:
		return "ANY_CRAPS"
	case BetEleven:
		return "ELEVEN"
	case BetAceDeuce:
		return "ACE_DEUCE"
	case BetAces:
		return "ACES"
	case BetBoxcars:
		return "BOXCARS"
	case BetPlace4:
		return "PLACE_4"
	case BetPlace5:
		return "PLACE_5"
	case BetPlace6:
		return "PLACE_6"
	case BetPlace8:
		return "PLACE_8"
	case BetPlace9:
		return "PLACE_9"
	case BetPlace10:
		return "PLACE_10"
	case BetPlaceNumbers:
		return "PLACE_NUMBERS"
	case BetPlaceInside:
		return "PLACE_INSIDE"
	case BetPlaceOutside:
		return "PLACE_OUTSIDE"
	case BetHard4:
		return "HARD_4"
	case BetHard6:
		return "HARD_6"
	case BetHard8:
		return "HARD_8"
	case BetHard10:
		return "HARD_10"
	case BetAllHardways:
		return "ALL_HARDWAYS"
	case BetPassOdds:
		return "PASS_ODDS"
	case BetDontPassOdds:
		return "DONT_PASS_ODDS"
	case BetBuy4:
		return "BUY_4"
	case BetBuy10:
		return "BUY_10"
	case BetLay4:
		return "LAY_4"
	case BetLay10:
		return "LAY_10"
	case BetBig6:
		return "BIG_6"
	case BetBig8:
		return "BIG_8"
	case BetHop:
		return "HOP"
	case BetHopHard6:
		return "HOP_HARD_6"
	case BetHopEasy8:
		return "HOP_EASY_8"
	case BetWorld:
		return "WORLD"
	case BetCAndE:
		return "C_AND_E"
	case BetHorn:
		return "HORN"
	case BetHornHigh11:
		return "HORN_HIGH_11"
	case BetHornHighAceDeuce:
		return "HORN_HIGH_ACE_DEUCE"
	// Missing bet type cases from canonical definitions
	// Buy bets
	case BetBuy5:
		return "BUY_5"
	case BetBuy6:
		return "BUY_6"
	case BetBuy8:
		return "BUY_8"
	case BetBuy9:
		return "BUY_9"
	// Lay bets
	case BetLay5:
		return "LAY_5"
	case BetLay6:
		return "LAY_6"
	case BetLay8:
		return "LAY_8"
	case BetLay9:
		return "LAY_9"
	// Place-to-lose bets
	case BetPlaceToLose4:
		return "PLACE_TO_LOSE_4"
	case BetPlaceToLose5:
		return "PLACE_TO_LOSE_5"
	case BetPlaceToLose6:
		return "PLACE_TO_LOSE_6"
	case BetPlaceToLose8:
		return "PLACE_TO_LOSE_8"
	case BetPlaceToLose9:
		return "PLACE_TO_LOSE_9"
	case BetPlaceToLose10:
		return "PLACE_TO_LOSE_10"
	// Horn high bets
	case BetHornHigh2:
		return "HORN_HIGH_2"
	case BetHornHigh3:
		return "HORN_HIGH_3"
	case BetHornHigh12:
		return "HORN_HIGH_12"
	// Hop bets (all combinations)
	case BetHop12:
		return "HOP_1_2"
	case BetHop13:
		return "HOP_1_3"
	case BetHop14:
		return "HOP_1_4"
	case BetHop15:
		return "HOP_1_5"
	case BetHop16:
		return "HOP_1_6"
	case BetHop23:
		return "HOP_2_3"
	case BetHop24:
		return "HOP_2_4"
	case BetHop25:
		return "HOP_2_5"
	case BetHop26:
		return "HOP_2_6"
	case BetHop34:
		return "HOP_3_4"
	case BetHop35:
		return "HOP_3_5"
	case BetHop36:
		return "HOP_3_6"
	case BetHop45:
		return "HOP_4_5"
	case BetHop46:
		return "HOP_4_6"
	case BetHop56:
		return "HOP_5_6"
	// Odds bets (specific types)
	case BetComeOdds:
		return "COME_ODDS"
	case BetDontComeOdds:
		return "DONT_COME_ODDS"
	default:
		return "UNKNOWN"
	}
}

func (i *Interpreter) evaluateCondition(expr Expression) (bool, error) {
	// Handle different types of conditions
	switch e := expr.(type) {
	case *IdentifierExpression:
		return i.evaluateIdentifierCondition(e)
	case *InfixExpression:
		return i.evaluateInfixCondition(e)
	default:
		return false, fmt.Errorf("unsupported condition type: %T", expr)
	}
}

func (i *Interpreter) evaluateIdentifierCondition(expr *IdentifierExpression) (bool, error) {
	switch expr.Value {
	case "COME_OUT_ROLL":
		return i.table.IsComeOut(), nil
	case "POINT_ESTABLISHED":
		return i.table.IsPoint(), nil
	default:
		return false, fmt.Errorf("unknown condition: %s", expr.Value)
	}
}

func (i *Interpreter) evaluateInfixCondition(expr *InfixExpression) (bool, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return false, err
	}

	switch expr.Operator {
	case "=":
		return left == right, nil
	case "!=":
		return left != right, nil
	case ">":
		return left > right, nil
	case "<":
		return left < right, nil
	case ">=":
		return left >= right, nil
	case "<=":
		return left <= right, nil
	default:
		return false, fmt.Errorf("unknown operator: %s", expr.Operator)
	}
}

func (i *Interpreter) evaluateExpression(expr Expression) (float64, error) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		return i.evaluateIdentifierExpression(e)
	case *NumberExpression:
		return e.Value, nil
	case *AmountExpression:
		return e.Value, nil
	default:
		return 0, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (i *Interpreter) evaluateIdentifierExpression(expr *IdentifierExpression) (float64, error) {
	switch expr.Value {
	case "POINT":
		point := i.table.GetPoint()
		if point == crapsgame.PointOff {
			return 0, nil
		}
		return float64(point), nil
	case "BANKROLL":
		// Get the first available player's bankroll
		for _, player := range i.table.Players {
			return player.Bankroll, nil
		}
		return 0, fmt.Errorf("no players at table")
	default:
		return 0, fmt.Errorf("unknown identifier: %s", expr.Value)
	}
}

func (i *Interpreter) evaluateIdentifierExpressionForPlayer(expr *IdentifierExpression, playerID string) (float64, error) {
	switch expr.Value {
	case "POINT":
		point := i.table.GetPoint()
		if point == crapsgame.PointOff {
			return 0, nil
		}
		return float64(point), nil
	case "BANKROLL":
		player, err := i.table.GetPlayer(playerID)
		if err != nil {
			return 0, fmt.Errorf("player %s not found", playerID)
		}
		return player.Bankroll, nil
	default:
		return 0, fmt.Errorf("unknown identifier: %s", expr.Value)
	}
}
