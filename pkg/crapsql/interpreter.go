package crapsql

import (
	"fmt"
	"strings"
	"time"

	"github.com/headswim/CrapsQL/pkg/crapsgame"
)

// Type aliases for backward compatibility
type Table = crapsgame.Table
type Player = crapsgame.Player
type Bet = crapsgame.Bet
type Roll = crapsgame.Roll

// Convenience functions for backward compatibility
func NewTable(minBet, maxBet float64, maxOdds int) *Table {
	return crapsgame.NewTable(minBet, maxBet, maxOdds)
}

func AddPlayer(table *Table, id, name string, bankroll float64) error {
	return table.AddPlayer(id, name, bankroll)
}

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

	return i.executeBetStatementForPlayer(stmt, playerID)
}

func (i *Interpreter) executeBetStatementForPlayer(stmt *BetStatement, playerID string) (string, error) {
	betType := i.betTypeToString(stmt.BetType.Type)
	numbers := extractNumbersForBetType(stmt.BetType)

	// Place the bet using the game engine
	placedBet, err := i.table.PlaceBet(playerID, betType, stmt.Amount.Value, numbers)
	if err != nil {
		return "", fmt.Errorf("failed to place bet: %v", err)
	}

	return fmt.Sprintf("✅ Placed $%.2f on %s", placedBet.Amount, betType), nil
}

func (i *Interpreter) executeConditionalStatement(stmt *ConditionalStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	return i.executeConditionalStatementForPlayer(stmt, playerID)
}

func (i *Interpreter) executeConditionalStatementForPlayer(stmt *ConditionalStatement, playerID string) (string, error) {
	// Evaluate condition using game state
	condition, err := i.evaluateConditionForPlayer(stmt.Condition, playerID)
	if err != nil {
		return "", fmt.Errorf("condition evaluation failed: %v", err)
	}

	var results []string
	if condition {
		// Execute consequence block
		for _, consequenceStmt := range stmt.Consequence.Statements {
			result, err := i.executeStatementForPlayer(consequenceStmt, playerID)
			if err != nil {
				return "", err
			}
			if result != "" {
				results = append(results, result)
			}
		}
	} else if stmt.Alternative != nil {
		// Execute alternative block
		for _, alternativeStmt := range stmt.Alternative.Statements {
			result, err := i.executeStatementForPlayer(alternativeStmt, playerID)
			if err != nil {
				return "", err
			}
			if result != "" {
				results = append(results, result)
			}
		}
	}

	return strings.Join(results, "\n"), nil
}

func (i *Interpreter) executeQueryStatement(stmt *QueryStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	return i.executeQueryStatementForPlayer(stmt, playerID)
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
	default:
		return "", fmt.Errorf("unknown query type: %v", stmt.Type)
	}
}

func (i *Interpreter) executeManagementStatement(stmt *ManagementStatement) (string, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return "", fmt.Errorf("no players at table - add a player first")
	}

	return i.executeManagementStatementForPlayer(stmt, playerID)
}

func (i *Interpreter) executeManagementStatementForPlayer(stmt *ManagementStatement, playerID string) (string, error) {
	amount, err := i.extractAmountFromExpression(stmt.Value)
	if err != nil {
		return "", fmt.Errorf("invalid amount: %v", err)
	}

	switch stmt.Type {
	case ManageBankroll:
		return i.executeSetBankroll(playerID, amount)
	case ManageMaxBet:
		return i.executeSetMaxBet(playerID, amount)
	case ManageMinBet:
		return i.executeSetMinBet(playerID, amount)
	case ManageWinGoal:
		return i.executeSetWinGoal(playerID, amount)
	case ManageLossLimit:
		return i.executeSetLossLimit(playerID, amount)
	default:
		return "", fmt.Errorf("unknown management type: %v", stmt.Type)
	}
}

func (i *Interpreter) executeSetBankroll(playerID string, amount float64) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	player.Bankroll = amount
	return fmt.Sprintf("✅ Set bankroll to $%.2f", amount), nil
}

func (i *Interpreter) executeSetMaxBet(playerID string, amount float64) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	player.MaxBet = amount
	return fmt.Sprintf("✅ Set max bet to $%.2f", amount), nil
}

func (i *Interpreter) executeSetMinBet(playerID string, amount float64) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	player.MinBet = amount
	return fmt.Sprintf("✅ Set min bet to $%.2f", amount), nil
}

func (i *Interpreter) executeSetWinGoal(playerID string, amount float64) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	player.WinGoal = amount
	return fmt.Sprintf("✅ Set win goal to $%.2f", amount), nil
}

func (i *Interpreter) executeSetLossLimit(playerID string, amount float64) (string, error) {
	player, err := i.table.GetPlayer(playerID)
	if err != nil {
		return "", fmt.Errorf("player %s not found", playerID)
	}

	player.LossLimit = amount
	return fmt.Sprintf("✅ Set loss limit to $%.2f", amount), nil
}

func (i *Interpreter) extractAmountFromExpression(expr Expression) (float64, error) {
	switch e := expr.(type) {
	case *NumberExpression:
		return e.Value, nil
	case *AmountExpression:
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

	return i.executeRemoveStatementForPlayer(stmt, playerID)
}

func (i *Interpreter) executeRemoveStatementForPlayer(stmt *RemoveStatement, playerID string) (string, error) {
	betType := i.betTypeToString(stmt.BetType.Type)

	// Remove the bet using the game engine
	err := i.table.RemoveBet(playerID, betType)
	if err != nil {
		return "", fmt.Errorf("failed to remove bet: %v", err)
	}

	return fmt.Sprintf("✅ Removed %s bet", betType), nil
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

	return i.executePressStatementForPlayer(stmt, playerID)
}

func (i *Interpreter) executePressStatementForPlayer(stmt *PressStatement, playerID string) (string, error) {
	betType := i.betTypeToString(stmt.BetType.Type)

	// Press the bet using the game engine
	err := i.table.PressBet(playerID, betType, stmt.Amount.Value)
	if err != nil {
		return "", fmt.Errorf("failed to press bet: %v", err)
	}

	return fmt.Sprintf("✅ Pressed %s bet by $%.2f", betType, stmt.Amount.Value), nil
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

	return i.executeTurnStatementForPlayer(stmt, playerID)
}

func (i *Interpreter) executeTurnStatementForPlayer(stmt *TurnStatement, playerID string) (string, error) {
	betType := i.betTypeToString(stmt.BetType.Type)

	// Turn the bet on/off using the game engine
	err := i.table.TurnBet(playerID, betType, stmt.Action == "ON")
	if err != nil {
		return "", fmt.Errorf("failed to turn bet %s: %v", stmt.Action, err)
	}

	return fmt.Sprintf("✅ Turned %s bet %s", betType, strings.ToLower(stmt.Action)), nil
}

func (i *Interpreter) evaluateConditionForPlayer(expr Expression, playerID string) (bool, error) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		return i.evaluateIdentifierConditionForPlayer(e, playerID)
	case *InfixExpression:
		return i.evaluateInfixConditionForPlayer(e, playerID)
	default:
		return false, fmt.Errorf("unsupported condition expression type: %T", expr)
	}
}

func (i *Interpreter) evaluateIdentifierConditionForPlayer(expr *IdentifierExpression, playerID string) (bool, error) {
	// Simple identifier conditions like "POINT"
	switch expr.Value {
	case "POINT":
		return i.table.IsPoint(), nil
	default:
		return false, fmt.Errorf("unknown condition identifier: %s", expr.Value)
	}
}

func (i *Interpreter) evaluateInfixConditionForPlayer(expr *InfixExpression, playerID string) (bool, error) {
	// Handle infix conditions like "POINT = 6"
	left, err := i.evaluateExpressionForPlayer(expr.Left, playerID)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpressionForPlayer(expr.Right, playerID)
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
		return false, fmt.Errorf("unsupported operator: %s", expr.Operator)
	}
}

func (i *Interpreter) evaluateExpressionForPlayer(expr Expression, playerID string) (float64, error) {
	switch e := expr.(type) {
	case *IdentifierExpression:
		return i.evaluateIdentifierExpressionForPlayer(e, playerID)
	case *NumberExpression:
		return e.Value, nil
	default:
		return 0, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (i *Interpreter) executeRollStatement(stmt *RollStatement) (string, error) {
	// Use the new clean game flow
	roll, results := i.table.ExecuteGameTurn()

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
	roll, allResults := i.table.RollDiceAndResolve()

	// Filter results for this player
	var playerResults []string
	for _, result := range allResults {
		// Check if this result is for our player
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
	case BetHorn:
		return "HORN"
	case BetHornHigh2:
		return "HORN_HIGH_2"
	case BetHornHigh3:
		return "HORN_HIGH_3"
	case BetHornHigh11:
		return "HORN_HIGH_11"
	case BetHornHigh12:
		return "HORN_HIGH_12"
	case BetWorld:
		return "WORLD"
	case BetCAndE:
		return "C_AND_E"
	default:
		return fmt.Sprintf("UNKNOWN_BET_TYPE_%d", betType)
	}
}

// Helper functions that should remain in the interpreter (language-specific)

func generateBetID() string {
	return fmt.Sprintf("bet_%d", time.Now().UnixNano())
}

func extractNumbersForBetType(expr *BetTypeExpression) []int {
	var numbers []int

	switch expr.Type {
	case BetPlaceNumbers:
		// Extract numbers from arguments
		for _, arg := range expr.Args {
			if numExpr, ok := arg.(*NumberExpression); ok {
				numbers = append(numbers, int(numExpr.Value))
			}
		}
	case BetHop:
		// Extract hop combination from arguments
		for _, arg := range expr.Args {
			if numExpr, ok := arg.(*NumberExpression); ok {
				numbers = append(numbers, int(numExpr.Value))
			}
		}
	}

	return numbers
}

// These functions are for language evaluation only, not game logic
func (i *Interpreter) evaluateCondition(expr Expression) (bool, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return false, fmt.Errorf("no players at table")
	}

	return i.evaluateConditionForPlayer(expr, playerID)
}

func (i *Interpreter) evaluateIdentifierCondition(expr *IdentifierExpression) (bool, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return false, fmt.Errorf("no players at table")
	}

	return i.evaluateIdentifierConditionForPlayer(expr, playerID)
}

func (i *Interpreter) evaluateInfixCondition(expr *InfixExpression) (bool, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return false, fmt.Errorf("no players at table")
	}

	return i.evaluateInfixConditionForPlayer(expr, playerID)
}

func (i *Interpreter) evaluateExpression(expr Expression) (float64, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return 0, fmt.Errorf("no players at table")
	}

	return i.evaluateExpressionForPlayer(expr, playerID)
}

func (i *Interpreter) evaluateIdentifierExpression(expr *IdentifierExpression) (float64, error) {
	var playerID string
	for id := range i.table.Players {
		playerID = id
		break
	}

	if playerID == "" {
		return 0, fmt.Errorf("no players at table")
	}

	return i.evaluateIdentifierExpressionForPlayer(expr, playerID)
}

func (i *Interpreter) evaluateIdentifierExpressionForPlayer(expr *IdentifierExpression, playerID string) (float64, error) {
	switch expr.Value {
	case "POINT":
		return float64(i.table.GetPointNumber()), nil
	case "BANKROLL":
		player, err := i.table.GetPlayer(playerID)
		if err != nil {
			return 0, err
		}
		return player.Bankroll, nil
	default:
		return 0, fmt.Errorf("unknown identifier: %s", expr.Value)
	}
}
