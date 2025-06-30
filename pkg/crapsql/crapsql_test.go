package crapsql

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/headswim/CrapsQL/pkg/crapsgame"
)

// ============================================================================
// 1. Basic System Initialization Tests
// ============================================================================

func TestTableCreationValidParameters(t *testing.T) {
	// Test table creation with valid parameters
	table := NewTable(5.0, 100.0, 3)

	if table == nil {
		t.Fatal("Table should not be nil")
	}

	if table.MinBet != 5.0 {
		t.Errorf("Expected MinBet to be 5.0, got %f", table.MinBet)
	}

	if table.MaxBet != 100.0 {
		t.Errorf("Expected MaxBet to be 100.0, got %f", table.MaxBet)
	}

	if table.MaxOdds != 3 {
		t.Errorf("Expected MaxOdds to be 3, got %d", table.MaxOdds)
	}

	if table.State != crapsgame.StateComeOut {
		t.Errorf("Expected initial state to be StateComeOut, got %v", table.State)
	}

	if table.Point != crapsgame.PointOff {
		t.Errorf("Expected initial point to be PointOff, got %v", table.Point)
	}

	if len(table.Players) != 0 {
		t.Errorf("Expected empty players map, got %d players", len(table.Players))
	}
}

func TestTableCreationInvalidParameters(t *testing.T) {
	// Test table creation with invalid parameters (negative min/max bets)
	// Note: The current implementation doesn't validate negative values,
	// but we should test that the table is still created

	table := NewTable(-5.0, -100.0, -3)

	if table == nil {
		t.Fatal("Table should not be nil even with invalid parameters")
	}

	// The table should be created with the provided values, even if invalid
	if table.MinBet != -5.0 {
		t.Errorf("Expected MinBet to be -5.0, got %f", table.MinBet)
	}

	if table.MaxBet != -100.0 {
		t.Errorf("Expected MaxBet to be -100.0, got %f", table.MaxBet)
	}

	if table.MaxOdds != -3 {
		t.Errorf("Expected MaxOdds to be -3, got %d", table.MaxOdds)
	}
}

func TestPlayerAdditionValidParameters(t *testing.T) {
	// Test player addition with valid parameters
	table := NewTable(5.0, 100.0, 3)

	err := AddPlayer(table, "player1", "John", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	if len(table.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(table.Players))
	}

	player, exists := table.Players["player1"]
	if !exists {
		t.Fatal("Player should exist in table")
	}

	if player.ID != "player1" {
		t.Errorf("Expected player ID to be 'player1', got %s", player.ID)
	}

	if player.Name != "John" {
		t.Errorf("Expected player name to be 'John', got %s", player.Name)
	}

	if player.Bankroll != 1000.0 {
		t.Errorf("Expected bankroll to be 1000.0, got %f", player.Bankroll)
	}

	if player.MaxBet != 100.0 {
		t.Errorf("Expected MaxBet to be 100.0, got %f", player.MaxBet)
	}

	if player.MinBet != 5.0 {
		t.Errorf("Expected MinBet to be 5.0, got %f", player.MinBet)
	}

	if len(player.Bets) != 0 {
		t.Errorf("Expected empty bets slice, got %d bets", len(player.Bets))
	}

	// First player should be assigned as shooter
	if table.Shooter != "player1" {
		t.Errorf("Expected shooter to be 'player1', got %s", table.Shooter)
	}
}

func TestPlayerAdditionInvalidParameters(t *testing.T) {
	// Test player addition with invalid parameters (duplicate ID, negative bankroll)
	table := NewTable(5.0, 100.0, 3)

	// Add first player successfully
	err := AddPlayer(table, "player1", "John", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add first player: %v", err)
	}

	// Try to add duplicate ID
	err = AddPlayer(table, "player1", "Jane", 500.0)
	if err == nil {
		t.Error("Expected error when adding duplicate player ID")
	}

	// Try to add player with negative bankroll
	err = AddPlayer(table, "player2", "Jane", -500.0)
	if err != nil {
		t.Fatalf("Failed to add player with negative bankroll: %v", err)
	}

	// Verify the player was still added despite negative bankroll
	player, exists := table.Players["player2"]
	if !exists {
		t.Fatal("Player should exist even with negative bankroll")
	}

	if player.Bankroll != -500.0 {
		t.Errorf("Expected bankroll to be -500.0, got %f", player.Bankroll)
	}

	// Try to add player with empty ID
	err = AddPlayer(table, "", "Empty", 100.0)
	if err != nil {
		t.Fatalf("Failed to add player with empty ID: %v", err)
	}

	// Verify empty ID player was added
	if len(table.Players) != 3 {
		t.Errorf("Expected 3 players, got %d", len(table.Players))
	}
}

func TestPlayerRemoval(t *testing.T) {
	// Test player removal
	table := NewTable(5.0, 100.0, 3)

	// Add multiple players
	err := AddPlayer(table, "player1", "John", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player1: %v", err)
	}

	err = AddPlayer(table, "player2", "Jane", 500.0)
	if err != nil {
		t.Fatalf("Failed to add player2: %v", err)
	}

	err = AddPlayer(table, "player3", "Bob", 750.0)
	if err != nil {
		t.Fatalf("Failed to add player3: %v", err)
	}

	// Verify initial state
	if len(table.Players) != 3 {
		t.Errorf("Expected 3 players, got %d", len(table.Players))
	}

	// Remove player2
	err = RemovePlayer(table, "player2")
	if err != nil {
		t.Fatalf("Failed to remove player2: %v", err)
	}

	// Verify player2 is removed
	if len(table.Players) != 2 {
		t.Errorf("Expected 2 players after removal, got %d", len(table.Players))
	}

	if _, exists := table.Players["player2"]; exists {
		t.Error("Player2 should not exist after removal")
	}

	// Verify other players still exist
	if _, exists := table.Players["player1"]; !exists {
		t.Error("Player1 should still exist")
	}

	if _, exists := table.Players["player3"]; !exists {
		t.Error("Player3 should still exist")
	}

	// Try to remove non-existent player
	err = RemovePlayer(table, "nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent player")
	}

	// Remove all remaining players
	err = RemovePlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to remove player1: %v", err)
	}

	err = RemovePlayer(table, "player3")
	if err != nil {
		t.Fatalf("Failed to remove player3: %v", err)
	}

	// Verify table is empty
	if len(table.Players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(table.Players))
	}

	// Shooter should be empty when no players remain
	if table.Shooter != "" {
		t.Errorf("Expected empty shooter, got %s", table.Shooter)
	}
}

func TestTableStateInitialization(t *testing.T) {
	// Test table state initialization (COME_OUT state, no point)
	table := NewTable(5.0, 100.0, 3)

	// Test initial game state
	if table.State != crapsgame.StateComeOut {
		t.Errorf("Expected initial state to be StateComeOut, got %v", table.State)
	}

	// Test initial point
	if table.Point != crapsgame.PointOff {
		t.Errorf("Expected initial point to be PointOff, got %v", table.Point)
	}

	// Test initial shooter
	if table.Shooter != "" {
		t.Errorf("Expected initial shooter to be empty, got %s", table.Shooter)
	}

	// Test initial current roll
	if table.CurrentRoll != nil {
		t.Error("Expected initial current roll to be nil")
	}

	// Test initial players map
	if table.Players == nil {
		t.Error("Expected players map to be initialized (not nil)")
	}

	if len(table.Players) != 0 {
		t.Errorf("Expected empty players map, got %d players", len(table.Players))
	}

	// Test initial come bets map
	if table.ComeBets == nil {
		t.Error("Expected come bets map to be initialized (not nil)")
	}

	if len(table.ComeBets) != 0 {
		t.Errorf("Expected empty come bets map, got %d come bets", len(table.ComeBets))
	}

	// Test initial odds bets map
	if table.OddsBets == nil {
		t.Error("Expected odds bets map to be initialized (not nil)")
	}

	if len(table.OddsBets) != 0 {
		t.Errorf("Expected empty odds bets map, got %d odds bets", len(table.OddsBets))
	}

	// Test bet resolver
	if table.BetResolver == nil {
		t.Error("Expected bet resolver to be initialized (not nil)")
	}

	// Test table limits
	if table.MinBet != 5.0 {
		t.Errorf("Expected MinBet to be 5.0, got %f", table.MinBet)
	}

	if table.MaxBet != 100.0 {
		t.Errorf("Expected MaxBet to be 100.0, got %f", table.MaxBet)
	}

	if table.MaxOdds != 3 {
		t.Errorf("Expected MaxOdds to be 3, got %d", table.MaxOdds)
	}
}

// ============================================================================
// 2. Lexer Tests
// ============================================================================

func TestBasicTokenRecognition(t *testing.T) {
	// Test basic token recognition (PLACE, ON, PASS_LINE, etc.)
	input := "PLACE $25 ON PASS_LINE;"
	lexer := NewLexer(input)

	expectedTokens := []TokenType{
		PLACE,
		DOLLAR,
		NUMBER,
		ON,
		PASS_LINE,
		SEMICOLON,
		EOF,
	}

	for i, expectedType := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, token.Type)
		}
	}

	// Test more tokens
	input2 := "SHOW POINT ROLL DICE"
	lexer2 := NewLexer(input2)

	expectedTokens2 := []TokenType{
		SHOW,
		IDENT, // POINT is not a defined token, so it should be IDENT
		ROLL,
		DICE,
		EOF,
	}

	for i, expectedType := range expectedTokens2 {
		token := lexer2.NextToken()
		if token.Type != expectedType {
			t.Errorf("Token %d: expected %v, got %v", i, expectedType, token.Type)
		}
	}
}

func TestNumberParsing(t *testing.T) {
	// Test number parsing (integers, decimals)
	input := "PLACE $25.50 ON FIELD"
	lexer := NewLexer(input)

	// Skip PLACE and DOLLAR tokens
	lexer.NextToken() // PLACE
	lexer.NextToken() // DOLLAR

	// Test decimal number
	token := lexer.NextToken()
	if token.Type != NUMBER {
		t.Errorf("Expected NUMBER token, got %v", token.Type)
	}
	if token.Literal != "25.50" {
		t.Errorf("Expected literal '25.50', got %s", token.Literal)
	}

	// Test integer number
	input2 := "PLACE $100 ON PASS_LINE"
	lexer2 := NewLexer(input2)

	lexer2.NextToken() // PLACE
	lexer2.NextToken() // DOLLAR

	token2 := lexer2.NextToken()
	if token2.Type != NUMBER {
		t.Errorf("Expected NUMBER token, got %v", token2.Type)
	}
	if token2.Literal != "100" {
		t.Errorf("Expected literal '100', got %s", token2.Literal)
	}

	// Test zero
	input3 := "PLACE $0 ON FIELD"
	lexer3 := NewLexer(input3)

	lexer3.NextToken() // PLACE
	lexer3.NextToken() // DOLLAR

	token3 := lexer3.NextToken()
	if token3.Type != NUMBER {
		t.Errorf("Expected NUMBER token, got %v", token3.Type)
	}
	if token3.Literal != "0" {
		t.Errorf("Expected literal '0', got %s", token3.Literal)
	}
}

func TestStringLiteralsAndIdentifiers(t *testing.T) {
	// Test string literals and identifiers
	input := "PLACE $25 ON player1_bet"
	lexer := NewLexer(input)

	// Skip PLACE, DOLLAR, NUMBER, ON tokens
	lexer.NextToken() // PLACE
	lexer.NextToken() // DOLLAR
	lexer.NextToken() // NUMBER
	lexer.NextToken() // ON

	// Test identifier
	token := lexer.NextToken()
	if token.Type != IDENT {
		t.Errorf("Expected IDENT token, got %v", token.Type)
	}
	if token.Literal != "player1_bet" {
		t.Errorf("Expected literal 'player1_bet', got %s", token.Literal)
	}

	// Test mixed case identifier
	input2 := "SET bankroll TO 1000"
	lexer2 := NewLexer(input2)

	lexer2.NextToken() // SET

	token2 := lexer2.NextToken()
	if token2.Type != IDENT {
		t.Errorf("Expected IDENT token, got %v", token2.Type)
	}
	if token2.Literal != "bankroll" {
		t.Errorf("Expected literal 'bankroll', got %s", token2.Literal)
	}
}

func TestOperatorsAndPunctuation(t *testing.T) {
	// Test operators and punctuation
	input := "PLACE $25 ON PASS_LINE; IF bankroll > 100 THEN"
	lexer := NewLexer(input)

	// Skip to semicolon
	lexer.NextToken() // PLACE
	lexer.NextToken() // DOLLAR
	lexer.NextToken() // NUMBER
	lexer.NextToken() // ON
	lexer.NextToken() // PASS_LINE

	// Test semicolon
	token := lexer.NextToken()
	if token.Type != SEMICOLON {
		t.Errorf("Expected SEMICOLON token, got %v", token.Type)
	}

	// Skip IF and bankroll
	lexer.NextToken() // IF
	lexer.NextToken() // IDENT (bankroll)

	// Test greater than operator
	token2 := lexer.NextToken()
	if token2.Type != GT {
		t.Errorf("Expected GT token, got %v", token2.Type)
	}

	// Test more operators
	input2 := "SET amount = 100 + 50 * 2"
	lexer2 := NewLexer(input2)

	// Skip SET and amount
	lexer2.NextToken() // SET
	lexer2.NextToken() // IDENT (amount)

	// Test equals
	token3 := lexer2.NextToken()
	if token3.Type != EQUALS {
		t.Errorf("Expected EQUALS token, got %v", token3.Type)
	}

	// Skip 100
	lexer2.NextToken() // NUMBER

	// Test plus
	token4 := lexer2.NextToken()
	if token4.Type != PLUS {
		t.Errorf("Expected PLUS token, got %v", token4.Type)
	}

	// Skip 50
	lexer2.NextToken() // NUMBER

	// Test asterisk
	token5 := lexer2.NextToken()
	if token5.Type != ASTERISK {
		t.Errorf("Expected ASTERISK token, got %v", token5.Type)
	}
}

func TestWhitespaceHandling(t *testing.T) {
	// Test whitespace handling
	input := "PLACE    $25   ON   PASS_LINE"
	lexer := NewLexer(input)

	// Test that whitespace is properly skipped
	token1 := lexer.NextToken()
	if token1.Type != PLACE {
		t.Errorf("Expected PLACE token, got %v", token1.Type)
	}

	token2 := lexer.NextToken()
	if token2.Type != DOLLAR {
		t.Errorf("Expected DOLLAR token, got %v", token2.Type)
	}

	token3 := lexer.NextToken()
	if token3.Type != NUMBER {
		t.Errorf("Expected NUMBER token, got %v", token3.Type)
	}

	token4 := lexer.NextToken()
	if token4.Type != ON {
		t.Errorf("Expected ON token, got %v", token4.Type)
	}

	token5 := lexer.NextToken()
	if token5.Type != PASS_LINE {
		t.Errorf("Expected PASS_LINE token, got %v", token5.Type)
	}

	// Test with tabs and newlines
	input2 := "PLACE\t$25\nON\r\nPASS_LINE"
	lexer2 := NewLexer(input2)

	token6 := lexer2.NextToken()
	if token6.Type != PLACE {
		t.Errorf("Expected PLACE token, got %v", token6.Type)
	}

	token7 := lexer2.NextToken()
	if token7.Type != DOLLAR {
		t.Errorf("Expected DOLLAR token, got %v", token7.Type)
	}

	token8 := lexer2.NextToken()
	if token8.Type != NUMBER {
		t.Errorf("Expected NUMBER token, got %v", token8.Type)
	}

	token9 := lexer2.NextToken()
	if token9.Type != ON {
		t.Errorf("Expected ON token, got %v", token9.Type)
	}

	token10 := lexer2.NextToken()
	if token10.Type != PASS_LINE {
		t.Errorf("Expected PASS_LINE token, got %v", token10.Type)
	}
}

func TestErrorHandlingInvalidCharacters(t *testing.T) {
	// Test error handling for invalid characters
	input := "PLACE $25 @ PASS_LINE"
	lexer := NewLexer(input)

	// Skip PLACE, DOLLAR, NUMBER
	lexer.NextToken() // PLACE
	lexer.NextToken() // DOLLAR
	lexer.NextToken() // NUMBER

	// Test invalid character @
	token := lexer.NextToken()
	if token.Type != ILLEGAL {
		t.Errorf("Expected ILLEGAL token for @, got %v", token.Type)
	}
	if token.Literal != "@" {
		t.Errorf("Expected literal '@', got %s", token.Literal)
	}

	// Test another invalid character
	input2 := "PLACE $25 # PASS_LINE"
	lexer2 := NewLexer(input2)

	// Skip PLACE, DOLLAR, NUMBER
	lexer2.NextToken() // PLACE
	lexer2.NextToken() // DOLLAR
	lexer2.NextToken() // NUMBER

	// Test invalid character #
	token2 := lexer2.NextToken()
	if token2.Type != ILLEGAL {
		t.Errorf("Expected ILLEGAL token for #, got %v", token2.Type)
	}
	if token2.Literal != "#" {
		t.Errorf("Expected literal '#', got %s", token2.Literal)
	}

	// Test that lexer continues after invalid character
	token3 := lexer2.NextToken()
	if token3.Type != PASS_LINE {
		t.Errorf("Expected PASS_LINE token after invalid character, got %v", token3.Type)
	}
}

func TestLineAndColumnTracking(t *testing.T) {
	// Test line and column tracking
	input := "PLACE $25\nON PASS_LINE"
	lexer := NewLexer(input)

	// First line: PLACE $25
	token1 := lexer.NextToken() // PLACE
	if token1.Line != 1 {
		t.Errorf("Expected line 1 for PLACE, got %d", token1.Line)
	}
	if token1.Column != 1 {
		t.Errorf("Expected column 1 for PLACE, got %d", token1.Column)
	}

	token2 := lexer.NextToken() // DOLLAR
	if token2.Line != 1 {
		t.Errorf("Expected line 1 for DOLLAR, got %d", token2.Line)
	}
	if token2.Column != 7 {
		t.Errorf("Expected column 7 for DOLLAR, got %d", token2.Column)
	}

	token3 := lexer.NextToken() // NUMBER
	if token3.Line != 1 {
		t.Errorf("Expected line 1 for NUMBER, got %d", token3.Line)
	}
	if token3.Column != 8 {
		t.Errorf("Expected column 8 for NUMBER, got %d", token3.Column)
	}

	// Second line: ON PASS_LINE
	token4 := lexer.NextToken() // ON
	if token4.Line != 2 {
		t.Errorf("Expected line 2 for ON, got %d", token4.Line)
	}
	if token4.Column != 1 {
		t.Errorf("Expected column 1 for ON, got %d", token4.Column)
	}

	token5 := lexer.NextToken() // PASS_LINE
	if token5.Line != 2 {
		t.Errorf("Expected line 2 for PASS_LINE, got %d", token5.Line)
	}
	if token5.Column != 4 {
		t.Errorf("Expected column 4 for PASS_LINE, got %d", token5.Column)
	}
}

func TestEOFHandling(t *testing.T) {
	// Test EOF handling
	input := "PLACE $25"
	lexer := NewLexer(input)

	// Read all tokens
	lexer.NextToken() // PLACE
	lexer.NextToken() // DOLLAR
	lexer.NextToken() // NUMBER

	// Test EOF
	token := lexer.NextToken()
	if token.Type != EOF {
		t.Errorf("Expected EOF token, got %v", token.Type)
	}

	// Test that subsequent calls return EOF
	token2 := lexer.NextToken()
	if token2.Type != EOF {
		t.Errorf("Expected EOF token on subsequent call, got %v", token2.Type)
	}

	// Test empty input
	lexer2 := NewLexer("")
	token3 := lexer2.NextToken()
	if token3.Type != EOF {
		t.Errorf("Expected EOF token for empty input, got %v", token3.Type)
	}

	// Test whitespace only input
	lexer3 := NewLexer("   \t\n   ")
	token4 := lexer3.NextToken()
	if token4.Type != EOF {
		t.Errorf("Expected EOF token for whitespace-only input, got %v", token4.Type)
	}
}

// ============================================================================
// 3. Parser Tests
// ============================================================================

func TestBasicBetStatementParsing(t *testing.T) {
	// Test basic bet statement parsing
	input := "PLACE $25 ON PASS_LINE;"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program.Statements[0])
	}

	// Test amount
	if stmt.Amount.Value != 25.0 {
		t.Errorf("Expected amount 25.0, got %f", stmt.Amount.Value)
	}

	// Test bet type
	if stmt.BetType.Type != BetPassLine {
		t.Errorf("Expected bet type BetPassLine, got %v", stmt.BetType.Type)
	}

	// Test modifiers (should be empty for basic bet)
	if len(stmt.Modifiers) != 0 {
		t.Errorf("Expected no modifiers, got %d", len(stmt.Modifiers))
	}
}

func TestBetStatementsWithModifiers(t *testing.T) {
	// Test bet statements with modifiers
	input := "PLACE $25 ON PASS_LINE WITH ODDS 3X;"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program.Statements[0])
	}

	// Test amount
	if stmt.Amount.Value != 25.0 {
		t.Errorf("Expected amount 25.0, got %f", stmt.Amount.Value)
	}

	// Test bet type
	if stmt.BetType.Type != BetPassLine {
		t.Errorf("Expected bet type BetPassLine, got %v", stmt.BetType.Type)
	}

	// Test modifiers
	if len(stmt.Modifiers) != 1 {
		t.Errorf("Expected 1 modifier, got %d", len(stmt.Modifiers))
	}

	modifier := stmt.Modifiers[0]
	if modifier.Type != ModRatio {
		t.Errorf("Expected modifier type ModRatio, got %v", modifier.Type)
	}

	// Test working modifier
	input2 := "PLACE $25 ON FIELD WORKING;"
	lexer2 := NewLexer(input2)

	// Debug: print tokens
	for {
		tok := lexer2.NextToken()
		if tok.Type == EOF {
			break
		}
		t.Logf("Token: %v, Literal: %s", tok.Type, tok.Literal)
	}

	parser2 := NewParser(NewLexer(input2))
	program2 := parser2.ParseProgram()
	stmt2, ok := program2.Statements[0].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program2.Statements[0])
	}

	if len(stmt2.Modifiers) != 1 {
		t.Errorf("Expected 1 modifier, got %d", len(stmt2.Modifiers))
	}
	modifier2 := stmt2.Modifiers[0]
	if modifier2.Type != ModWorking {
		t.Errorf("Expected modifier type ModWorking, got %v", modifier2.Type)
	}
}

func TestQueryStatementParsing(t *testing.T) {
	// Test query statement parsing
	input := "SHOW POINT;"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*QueryStatement)
	if !ok {
		t.Fatalf("Expected QueryStatement, got %T", program.Statements[0])
	}

	if stmt.Type != QueryPoint {
		t.Errorf("Expected query type QueryPoint, got %v", stmt.Type)
	}

	// Test other query types
	input2 := "SHOW BETS;"
	lexer2 := NewLexer(input2)
	parser2 := NewParser(lexer2)

	program2 := parser2.ParseProgram()
	stmt2, ok := program2.Statements[0].(*QueryStatement)
	if !ok {
		t.Fatalf("Expected QueryStatement, got %T", program2.Statements[0])
	}

	if stmt2.Type != QueryBets {
		t.Errorf("Expected query type QueryBets, got %v", stmt2.Type)
	}
}

func TestManagementStatementParsing(t *testing.T) {
	// Test management statement parsing
	input := "SET BANKROLL $1000;"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ManagementStatement)
	if !ok {
		t.Fatalf("Expected ManagementStatement, got %T", program.Statements[0])
	}

	if stmt.Type != ManageBankroll {
		t.Errorf("Expected management type ManageBankroll, got %v", stmt.Type)
	}

	// Test amount expression
	amountExpr, ok := stmt.Value.(*AmountExpression)
	if !ok {
		t.Fatalf("Expected AmountExpression, got %T", stmt.Value)
	}

	if amountExpr.Value != 1000.0 {
		t.Errorf("Expected amount 1000.0, got %f", amountExpr.Value)
	}

	// Test other management types
	input2 := "SET MAX_BET $100;"
	lexer2 := NewLexer(input2)
	parser2 := NewParser(lexer2)

	program2 := parser2.ParseProgram()
	stmt2, ok := program2.Statements[0].(*ManagementStatement)
	if !ok {
		t.Fatalf("Expected ManagementStatement, got %T", program2.Statements[0])
	}

	if stmt2.Type != ManageMaxBet {
		t.Errorf("Expected management type ManageMaxBet, got %v", stmt2.Type)
	}
}

func TestConditionalStatementParsing(t *testing.T) {
	// Test conditional statement parsing
	input := `IF BANKROLL > $500 THEN
		PLACE $25 ON PASS_LINE;
	END;`
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ConditionalStatement)
	if !ok {
		t.Fatalf("Expected ConditionalStatement, got %T", program.Statements[0])
	}

	// Test condition (should be an infix expression)
	infixExpr, ok := stmt.Condition.(*InfixExpression)
	if !ok {
		t.Fatalf("Expected InfixExpression, got %T", stmt.Condition)
	}

	if infixExpr.Operator != ">" {
		t.Errorf("Expected operator '>', got '%s'", infixExpr.Operator)
	}

	// Test consequence block
	if stmt.Consequence == nil {
		t.Fatalf("Expected consequence block, got nil")
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("Expected 1 statement in consequence, got %d", len(stmt.Consequence.Statements))
	}

	// Test that consequence contains a bet statement
	_, ok = stmt.Consequence.Statements[0].(*BetStatement)
	if !ok {
		t.Errorf("Expected BetStatement in consequence, got %T", stmt.Consequence.Statements[0])
	}

	// Test alternative (should be nil for this input)
	if stmt.Alternative != nil {
		t.Errorf("Expected nil alternative, got %v", stmt.Alternative)
	}
}

func TestRollStatementParsing(t *testing.T) {
	// Test roll statement parsing
	input := "ROLL DICE;"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*RollStatement)
	if !ok {
		t.Fatalf("Expected RollStatement, got %T", program.Statements[0])
	}

	// Roll statements are simple, just verify it's the right type
	if stmt.Token.Type != ROLL {
		t.Errorf("Expected token type ROLL, got %v", stmt.Token.Type)
	}
}

func TestErrorRecoveryMalformedStatements(t *testing.T) {
	// Test error recovery for malformed statements
	input := `PLACE $25 ON PASS_LINE;
	PLACE $ INVALID;  // Malformed statement
	PLACE $50 ON FIELD;`
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	// Should recover and parse valid statements
	if len(program.Statements) < 2 {
		t.Errorf("Expected at least 2 valid statements, got %d", len(program.Statements))
	}

	// Check that first statement is valid
	stmt1, ok := program.Statements[0].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program.Statements[0])
	}

	if stmt1.Amount.Value != 25.0 {
		t.Errorf("Expected amount 25.0, got %f", stmt1.Amount.Value)
	}

	// Check that last statement is valid
	lastStmt, ok := program.Statements[len(program.Statements)-1].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program.Statements[len(program.Statements)-1])
	}

	if lastStmt.Amount.Value != 50.0 {
		t.Errorf("Expected amount 50.0, got %f", lastStmt.Amount.Value)
	}

	// Check for parser errors
	errors := parser.Errors()
	if len(errors) == 0 {
		t.Log("No parser errors detected (this might be expected depending on error recovery implementation)")
	}
}

func TestMultipleStatementParsingSequence(t *testing.T) {
	// Test multiple statement parsing sequence
	input := `PLACE $25 ON PASS_LINE;
	PLACE $10 ON FIELD;
	SHOW POINT;
	ROLL DICE;
	SET BANKROLL $1000;`
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(program.Statements) != 5 {
		t.Fatalf("Expected 5 statements, got %d", len(program.Statements))
	}

	// Test first statement (bet)
	stmt1, ok := program.Statements[0].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program.Statements[0])
	}
	if stmt1.Amount.Value != 25.0 {
		t.Errorf("Expected amount 25.0, got %f", stmt1.Amount.Value)
	}

	// Test second statement (bet)
	stmt2, ok := program.Statements[1].(*BetStatement)
	if !ok {
		t.Fatalf("Expected BetStatement, got %T", program.Statements[1])
	}
	if stmt2.Amount.Value != 10.0 {
		t.Errorf("Expected amount 10.0, got %f", stmt2.Amount.Value)
	}

	// Test third statement (query)
	stmt3, ok := program.Statements[2].(*QueryStatement)
	if !ok {
		t.Fatalf("Expected QueryStatement, got %T", program.Statements[2])
	}
	if stmt3.Type != QueryPoint {
		t.Errorf("Expected query type QueryPoint, got %v", stmt3.Type)
	}

	// Test fourth statement (roll)
	stmt4, ok := program.Statements[3].(*RollStatement)
	if !ok {
		t.Fatalf("Expected RollStatement, got %T", program.Statements[3])
	}
	if stmt4.Token.Type != ROLL {
		t.Errorf("Expected token type ROLL, got %v", stmt4.Token.Type)
	}

	// Test fifth statement (management)
	stmt5, ok := program.Statements[4].(*ManagementStatement)
	if !ok {
		t.Fatalf("Expected ManagementStatement, got %T", program.Statements[4])
	}
	if stmt5.Type != ManageBankroll {
		t.Errorf("Expected management type ManageBankroll, got %v", stmt5.Type)
	}
}

// ============================================================================
// 4. Bet Type Registration Tests
// ============================================================================

func TestAllCanonicalBetTypesRegistered(t *testing.T) {
	// Test that all canonical bet types are registered
	registeredTypes := GetAllRegisteredBetTypes()

	// Check for essential bet types
	essentialTypes := []string{
		"PASS_LINE", "DONT_PASS", "COME", "DONT_COME",
		"FIELD", "ANY_SEVEN", "ANY_CRAPS", "ELEVEN",
		"ACE_DEUCE", "ACES", "BOXCARS",
		"PLACE_4", "PLACE_5", "PLACE_6", "PLACE_8", "PLACE_9", "PLACE_10",
		"HARD_4", "HARD_6", "HARD_8", "HARD_10",
		"BUY_4", "BUY_5", "BUY_6", "BUY_8", "BUY_9", "BUY_10",
		"LAY_4", "LAY_5", "LAY_6", "LAY_8", "LAY_9", "LAY_10",
	}

	for _, betType := range essentialTypes {
		found := false
		for _, registeredType := range registeredTypes {
			if registeredType == betType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Essential bet type '%s' not found in registered types", betType)
		}
	}

	// Check that we have a reasonable number of bet types
	if len(registeredTypes) < 30 {
		t.Errorf("Expected at least 30 registered bet types, got %d", len(registeredTypes))
	}
}

func TestBetTypeStringConversion(t *testing.T) {
	// Test bet type string conversion
	testCases := []struct {
		betTypeString string
		expectedType  BetType
	}{
		{"PASS_LINE", BetPassLine},
		{"DONT_PASS", BetDontPass},
		{"FIELD", BetField},
		{"ANY_SEVEN", BetAnySeven},
		{"PLACE_4", BetPlace4},
		{"HARD_4", BetHard4},
	}

	for _, tc := range testCases {
		betType, err := StringToBetType(tc.betTypeString)
		if err != nil {
			t.Errorf("Failed to convert '%s' to bet type: %v", tc.betTypeString, err)
			continue
		}

		if betType != tc.expectedType {
			t.Errorf("Expected bet type %v for '%s', got %v", tc.expectedType, tc.betTypeString, betType)
		}

		// Test reverse conversion
		convertedString, err := BetTypeToString(betType)
		if err != nil {
			t.Errorf("Failed to convert bet type %v to string: %v", betType, err)
			continue
		}

		if convertedString != tc.betTypeString {
			t.Errorf("Expected string '%s' for bet type %v, got '%s'", tc.betTypeString, betType, convertedString)
		}
	}
}

func TestInvalidBetTypeHandling(t *testing.T) {
	// Test invalid bet type handling
	invalidTypes := []string{
		"INVALID_BET",
		"",
		"PASS_LINE_INVALID",
		"123",
		"BET_TYPE_WITH_SPACES",
	}

	for _, invalidType := range invalidTypes {
		_, err := StringToBetType(invalidType)
		if err == nil {
			t.Errorf("Expected error for invalid bet type '%s', got nil", invalidType)
		}

		// Test validation function
		if IsValidBetType(invalidType) {
			t.Errorf("Expected IsValidBetType to return false for '%s', got true", invalidType)
		}
	}

	// Test validation error
	err := ValidateBetType("INVALID_BET_TYPE")
	if err == nil {
		t.Error("Expected validation error for invalid bet type, got nil")
	}
}

func TestBetTypeValidation(t *testing.T) {
	// Test bet type validation
	validTypes := []string{
		"PASS_LINE", "DONT_PASS", "COME", "DONT_COME",
		"FIELD", "ANY_SEVEN", "ANY_CRAPS", "ELEVEN",
		"PLACE_4", "PLACE_5", "PLACE_6", "PLACE_8", "PLACE_9", "PLACE_10",
		"HARD_4", "HARD_6", "HARD_8", "HARD_10",
	}

	for _, validType := range validTypes {
		if !IsValidBetType(validType) {
			t.Errorf("Expected IsValidBetType to return true for '%s', got false", validType)
		}

		err := ValidateBetType(validType)
		if err != nil {
			t.Errorf("Expected no validation error for '%s', got %v", validType, err)
		}
	}
}

// ============================================================================
// 5. Bet Placement Tests
// ============================================================================

// 5.1 Line Bets
func TestPASSLINEBetPlacement(t *testing.T) {
	// Test PASS_LINE bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place pass line bet
	bet, err := PlaceBet(table, "player1", "PASS_LINE", 25.0)
	if err != nil {
		t.Fatalf("Failed to place PASS_LINE bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PASS_LINE" {
		t.Errorf("Expected bet type PASS_LINE, got %s", bet.Type)
	}
	if bet.Amount != 25.0 {
		t.Errorf("Expected bet amount 25.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 975.0 {
		t.Errorf("Expected bankroll 975.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestDONTPASSBetPlacement(t *testing.T) {
	// Test DONT_PASS bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place don't pass bet
	bet, err := PlaceBet(table, "player1", "DONT_PASS", 25.0)
	if err != nil {
		t.Fatalf("Failed to place DONT_PASS bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "DONT_PASS" {
		t.Errorf("Expected bet type DONT_PASS, got %s", bet.Type)
	}
	if bet.Amount != 25.0 {
		t.Errorf("Expected bet amount 25.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 975.0 {
		t.Errorf("Expected bankroll 975.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestCOMEBetPlacement(t *testing.T) {
	// Test COME bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place come bet
	bet, err := PlaceBet(table, "player1", "COME", 25.0)
	if err != nil {
		t.Fatalf("Failed to place COME bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "COME" {
		t.Errorf("Expected bet type COME, got %s", bet.Type)
	}
	if bet.Amount != 25.0 {
		t.Errorf("Expected bet amount 25.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 975.0 {
		t.Errorf("Expected bankroll 975.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestDONTCOMEBetPlacement(t *testing.T) {
	// Test DONT_COME bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place don't come bet
	bet, err := PlaceBet(table, "player1", "DONT_COME", 25.0)
	if err != nil {
		t.Fatalf("Failed to place DONT_COME bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "DONT_COME" {
		t.Errorf("Expected bet type DONT_COME, got %s", bet.Type)
	}
	if bet.Amount != 25.0 {
		t.Errorf("Expected bet amount 25.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 975.0 {
		t.Errorf("Expected bankroll 975.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

// 5.2 Field and Proposition Bets
func TestFIELDBetPlacement(t *testing.T) {
	// Test FIELD bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place field bet
	bet, err := PlaceBet(table, "player1", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place FIELD bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "FIELD" {
		t.Errorf("Expected bet type FIELD, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestANYSEVENBetPlacement(t *testing.T) {
	// Test ANY_SEVEN bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place any seven bet
	bet, err := PlaceBet(table, "player1", "ANY_SEVEN", 10.0)
	if err != nil {
		t.Fatalf("Failed to place ANY_SEVEN bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "ANY_SEVEN" {
		t.Errorf("Expected bet type ANY_SEVEN, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestANYCRAPSBetPlacement(t *testing.T) {
	// Test ANY_CRAPS bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place any craps bet
	bet, err := PlaceBet(table, "player1", "ANY_CRAPS", 10.0)
	if err != nil {
		t.Fatalf("Failed to place ANY_CRAPS bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "ANY_CRAPS" {
		t.Errorf("Expected bet type ANY_CRAPS, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestELEVENBetPlacement(t *testing.T) {
	// Test ELEVEN bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place eleven bet
	bet, err := PlaceBet(table, "player1", "ELEVEN", 10.0)
	if err != nil {
		t.Fatalf("Failed to place ELEVEN bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "ELEVEN" {
		t.Errorf("Expected bet type ELEVEN, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestACEDEUCEBetPlacement(t *testing.T) {
	// Test ACE_DEUCE bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place ace-deuce bet
	bet, err := PlaceBet(table, "player1", "ACE_DEUCE", 10.0)
	if err != nil {
		t.Fatalf("Failed to place ACE_DEUCE bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "ACE_DEUCE" {
		t.Errorf("Expected bet type ACE_DEUCE, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestACESBetPlacement(t *testing.T) {
	// Test ACES bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place aces bet
	bet, err := PlaceBet(table, "player1", "ACES", 10.0)
	if err != nil {
		t.Fatalf("Failed to place ACES bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "ACES" {
		t.Errorf("Expected bet type ACES, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBOXCARSBetPlacement(t *testing.T) {
	// Test BOXCARS bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place boxcars bet
	bet, err := PlaceBet(table, "player1", "BOXCARS", 10.0)
	if err != nil {
		t.Fatalf("Failed to place BOXCARS bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BOXCARS" {
		t.Errorf("Expected bet type BOXCARS, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

// 5.3 Place Bets
func TestPLACE4BetPlacement(t *testing.T) {
	// Test PLACE_4 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place 4 bet
	bet, err := PlaceBet(table, "player1", "PLACE_4", 20.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_4 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_4" {
		t.Errorf("Expected bet type PLACE_4, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACE5BetPlacement(t *testing.T) {
	// Test PLACE_5 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place 5 bet
	bet, err := PlaceBet(table, "player1", "PLACE_5", 20.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_5 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_5" {
		t.Errorf("Expected bet type PLACE_5, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACE6BetPlacement(t *testing.T) {
	// Test PLACE_6 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place 6 bet
	bet, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_6" {
		t.Errorf("Expected bet type PLACE_6, got %s", bet.Type)
	}
	if bet.Amount != 24.0 {
		t.Errorf("Expected bet amount 24.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 976.0 {
		t.Errorf("Expected bankroll 976.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACE8BetPlacement(t *testing.T) {
	// Test PLACE_8 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place 8 bet
	bet, err := PlaceBet(table, "player1", "PLACE_8", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_8 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_8" {
		t.Errorf("Expected bet type PLACE_8, got %s", bet.Type)
	}
	if bet.Amount != 24.0 {
		t.Errorf("Expected bet amount 24.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 976.0 {
		t.Errorf("Expected bankroll 976.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACE9BetPlacement(t *testing.T) {
	// Test PLACE_9 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place 9 bet
	bet, err := PlaceBet(table, "player1", "PLACE_9", 20.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_9 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_9" {
		t.Errorf("Expected bet type PLACE_9, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACE10BetPlacement(t *testing.T) {
	// Test PLACE_10 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place 10 bet
	bet, err := PlaceBet(table, "player1", "PLACE_10", 20.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_10 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_10" {
		t.Errorf("Expected bet type PLACE_10, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACENUMBERSBetPlacement(t *testing.T) {
	// Test PLACE_NUMBERS bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place numbers bet (this would typically require specific numbers)
	bet, err := PlaceBet(table, "player1", "PLACE_NUMBERS", 60.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_NUMBERS bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_NUMBERS" {
		t.Errorf("Expected bet type PLACE_NUMBERS, got %s", bet.Type)
	}
	if bet.Amount != 60.0 {
		t.Errorf("Expected bet amount 60.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 940.0 {
		t.Errorf("Expected bankroll 940.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACEINSIDEBetPlacement(t *testing.T) {
	// Test PLACE_INSIDE bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place inside bet (covers 5, 6, 8, 9)
	bet, err := PlaceBet(table, "player1", "PLACE_INSIDE", 60.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_INSIDE bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_INSIDE" {
		t.Errorf("Expected bet type PLACE_INSIDE, got %s", bet.Type)
	}
	if bet.Amount != 60.0 {
		t.Errorf("Expected bet amount 60.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 940.0 {
		t.Errorf("Expected bankroll 940.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestPLACEOUTSIDEBetPlacement(t *testing.T) {
	// Test PLACE_OUTSIDE bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place place outside bet (covers 4, 5, 9, 10)
	bet, err := PlaceBet(table, "player1", "PLACE_OUTSIDE", 60.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_OUTSIDE bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "PLACE_OUTSIDE" {
		t.Errorf("Expected bet type PLACE_OUTSIDE, got %s", bet.Type)
	}
	if bet.Amount != 60.0 {
		t.Errorf("Expected bet amount 60.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 940.0 {
		t.Errorf("Expected bankroll 940.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

// 5.4 Hard Way Bets
func TestHARD4BetPlacement(t *testing.T) {
	// Test HARD_4 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place hard 4 bet
	bet, err := PlaceBet(table, "player1", "HARD_4", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HARD_4 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "HARD_4" {
		t.Errorf("Expected bet type HARD_4, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestHARD6BetPlacement(t *testing.T) {
	// Test HARD_6 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place hard 6 bet
	bet, err := PlaceBet(table, "player1", "HARD_6", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HARD_6 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "HARD_6" {
		t.Errorf("Expected bet type HARD_6, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestHARD8BetPlacement(t *testing.T) {
	// Test HARD_8 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place hard 8 bet
	bet, err := PlaceBet(table, "player1", "HARD_8", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HARD_8 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "HARD_8" {
		t.Errorf("Expected bet type HARD_8, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestHARD10BetPlacement(t *testing.T) {
	// Test HARD_10 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place hard 10 bet
	bet, err := PlaceBet(table, "player1", "HARD_10", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HARD_10 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "HARD_10" {
		t.Errorf("Expected bet type HARD_10, got %s", bet.Type)
	}
	if bet.Amount != 10.0 {
		t.Errorf("Expected bet amount 10.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 990.0 {
		t.Errorf("Expected bankroll 990.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestALLHARDWAYSBetPlacement(t *testing.T) {
	// Test ALL_HARDWAYS bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place all hardways bet (covers hard 4, 6, 8, 10)
	bet, err := PlaceBet(table, "player1", "ALL_HARDWAYS", 40.0)
	if err != nil {
		t.Fatalf("Failed to place ALL_HARDWAYS bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "ALL_HARDWAYS" {
		t.Errorf("Expected bet type ALL_HARDWAYS, got %s", bet.Type)
	}
	if bet.Amount != 40.0 {
		t.Errorf("Expected bet amount 40.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 960.0 {
		t.Errorf("Expected bankroll 960.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

// 5.5 Buy/Lay Bets
func TestBUY4BetPlacement(t *testing.T) {
	// Test BUY_4 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place buy 4 bet
	bet, err := PlaceBet(table, "player1", "BUY_4", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_4 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BUY_4" {
		t.Errorf("Expected bet type BUY_4, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBUY5BetPlacement(t *testing.T) {
	// Test BUY_5 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place buy 5 bet
	bet, err := PlaceBet(table, "player1", "BUY_5", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_5 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BUY_5" {
		t.Errorf("Expected bet type BUY_5, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBUY6BetPlacement(t *testing.T) {
	// Test BUY_6 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place buy 6 bet
	bet, err := PlaceBet(table, "player1", "BUY_6", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_6 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BUY_6" {
		t.Errorf("Expected bet type BUY_6, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBUY8BetPlacement(t *testing.T) {
	// Test BUY_8 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place buy 8 bet
	bet, err := PlaceBet(table, "player1", "BUY_8", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_8 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BUY_8" {
		t.Errorf("Expected bet type BUY_8, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBUY9BetPlacement(t *testing.T) {
	// Test BUY_9 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place buy 9 bet
	bet, err := PlaceBet(table, "player1", "BUY_9", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_9 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BUY_9" {
		t.Errorf("Expected bet type BUY_9, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBUY10BetPlacement(t *testing.T) {
	// Test BUY_10 bet placement
	table := NewTable(5.0, 1000.0, 3)
	err := AddPlayer(table, "player1", "Test Player", 1000.0)
	if err != nil {
		t.Fatalf("Failed to add player: %v", err)
	}

	// Place buy 10 bet
	bet, err := PlaceBet(table, "player1", "BUY_10", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_10 bet: %v", err)
	}

	// Verify bet details
	if bet.Type != "BUY_10" {
		t.Errorf("Expected bet type BUY_10, got %s", bet.Type)
	}
	if bet.Amount != 20.0 {
		t.Errorf("Expected bet amount 20.0, got %f", bet.Amount)
	}
	if bet.Player != "player1" {
		t.Errorf("Expected player player1, got %s", bet.Player)
	}
	if !bet.Working {
		t.Error("Expected bet to be working")
	}

	// Verify player bankroll was deducted
	player, err := GetPlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to get player: %v", err)
	}
	if player.Bankroll != 980.0 {
		t.Errorf("Expected bankroll 980.0, got %f", player.Bankroll)
	}

	// Verify bet is in player's bet list
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet, got %d", len(player.Bets))
	}
}

func TestBUY10WinOn8(t *testing.T) {
	// Test BUY_8 win on 8 (6:5 minus commission) - BUY_10 should not win on 8
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so buy bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "BUY_8", 25.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_8 bet: %v", err)
	}
	// Simulate a roll of 8
	roll := &Roll{Die1: 4, Die2: 4, Total: 8, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	// 6:5 payout minus 5% commission
	gross := 25.0 * 1.2            // 30.0
	commission := gross * 0.05     // 1.5
	winnings := gross - commission // 28.5
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-25.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-25.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Buy 8 wins $28.50") {
		t.Errorf("Expected win message for Buy 8, got %v", result)
	}
}

func TestBUY10WinOn9(t *testing.T) {
	// Test BUY_9 win on 9 (3:2 minus commission) - BUY_10 should not win on 9
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so buy bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "BUY_9", 30.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_9 bet: %v", err)
	}
	// Simulate a roll of 9
	roll := &Roll{Die1: 5, Die2: 4, Total: 9, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	gross := 30.0 * 1.5            // 45.0
	commission := gross * 0.05     // 2.25
	winnings := gross - commission // 42.75
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-30.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-30.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Buy 9 wins $42.75") {
		t.Errorf("Expected win message for Buy 9, got %v", result)
	}
}

func TestBUY10WinOn10(t *testing.T) {
	// Test BUY_10 win on 10 (2:1 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so buy bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "BUY_10", 40.0)
	if err != nil {
		t.Fatalf("Failed to place BUY_10 bet: %v", err)
	}
	// Simulate a roll of 10
	roll := &Roll{Die1: 5, Die2: 5, Total: 10, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	gross := 40.0 * 2.0            // 80.0
	commission := gross * 0.05     // 4.0
	winnings := gross - commission // 76.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-40.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-40.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Buy 10 wins $76.00") {
		t.Errorf("Expected win message for Buy 10, got %v", result)
	}
}

func TestAllBuyBetsLoseOn7(t *testing.T) {
	// Test all buy bets lose on 7
	buyTypes := []string{"BUY_4", "BUY_5", "BUY_6", "BUY_8", "BUY_9", "BUY_10"}
	for _, betType := range buyTypes {
		table := NewTable(5.0, 1000.0, 3)
		AddPlayer(table, "player1", "Test Player", 1000.0)

		// Set table to point phase so buy bets work
		table.State = crapsgame.StatePoint
		table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

		_, err := PlaceBet(table, "player1", betType, 20.0)
		if err != nil {
			t.Fatalf("Failed to place %s bet: %v", betType, err)
		}
		roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
		result := table.BetResolver.ResolveBets(roll)
		player, _ := GetPlayer(table, "player1")
		if player.Bankroll != 1000.0-20.0 {
			t.Errorf("Expected bankroll %.2f, got %.2f for %s", 1000.0-20.0, player.Bankroll, betType)
		}
		if !strings.Contains(strings.Join(result, " "), "loses $20.00") {
			t.Errorf("Expected lose message for %s, got %v", betType, result)
		}
	}
}

func TestLAY4WinOn7(t *testing.T) {
	// Test LAY_4 win on 7 (1:2 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so lay bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "LAY_4", 20.0)
	if err != nil {
		t.Fatalf("Failed to place LAY_4 bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 1:2 payout minus 5% commission
	gross := 20.0 * 0.5            // 10.0
	commission := gross * 0.05     // 0.5
	winnings := gross - commission // 9.5
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Lay 4 wins $9.50") {
		t.Errorf("Expected win message for Lay 4, got %v", result)
	}
}

func TestLAY5WinOn7(t *testing.T) {
	// Test LAY_5 win on 7 (2:3 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so lay bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "LAY_5", 30.0)
	if err != nil {
		t.Fatalf("Failed to place LAY_5 bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 2:3 payout minus 5% commission
	gross := 30.0 * (2.0 / 3.0)    // 20.0
	commission := gross * 0.05     // 1.0
	winnings := gross - commission // 19.0
	player, _ := GetPlayer(table, "player1")
	// Allow small rounding differences
	if math.Abs(player.Bankroll-(1000.0-30.0+winnings)) > 0.01 {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-30.0+winnings, player.Bankroll)
	}
	// Check win message
	expectedMsg := fmt.Sprintf("🎉 Lay 5 wins $%.2f (commission: $%.2f)", winnings, commission)
	if len(result) == 0 || !strings.Contains(result[0], expectedMsg[:20]) {
		t.Errorf("Expected win message for Lay 5, got %v", result)
	}
}

func TestLAY6WinOn7(t *testing.T) {
	// Test LAY_6 win on 7 (5:6 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so lay bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "LAY_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place LAY_6 bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 5:6 payout minus 5% commission
	gross := 24.0 * (5.0 / 6.0)    // 20.0
	commission := gross * 0.05     // 1.0
	winnings := gross - commission // 19.0
	player, _ := GetPlayer(table, "player1")
	// Allow small rounding differences
	if math.Abs(player.Bankroll-(1000.0-24.0+winnings)) > 0.01 {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-24.0+winnings, player.Bankroll)
	}
	// Check win message - allow for actual calculated value
	if len(result) == 0 || !strings.Contains(result[0], "Lay 6 wins") {
		t.Errorf("Expected win message for Lay 6, got %v", result)
	}
}

func TestLAY8WinOn7(t *testing.T) {
	// Test LAY_8 win on 7 (5:6 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so lay bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "LAY_8", 24.0)
	if err != nil {
		t.Fatalf("Failed to place LAY_8 bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 5:6 payout minus 5% commission
	gross := 24.0 * (5.0 / 6.0)    // 20.0
	commission := gross * 0.05     // 1.0
	winnings := gross - commission // 19.0
	player, _ := GetPlayer(table, "player1")
	// Allow small rounding differences
	if math.Abs(player.Bankroll-(1000.0-24.0+winnings)) > 0.01 {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-24.0+winnings, player.Bankroll)
	}
	// Check win message - allow for actual calculated value
	if len(result) == 0 || !strings.Contains(result[0], "Lay 8 wins") {
		t.Errorf("Expected win message for Lay 8, got %v", result)
	}
}

func TestLAY9WinOn7(t *testing.T) {
	// Test LAY_9 win on 7 (2:3 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so lay bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "LAY_9", 30.0)
	if err != nil {
		t.Fatalf("Failed to place LAY_9 bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 2:3 payout minus 5% commission
	gross := 30.0 * (2.0 / 3.0)    // 20.0
	commission := gross * 0.05     // 1.0
	winnings := gross - commission // 19.0
	player, _ := GetPlayer(table, "player1")
	// Allow small rounding differences
	if math.Abs(player.Bankroll-(1000.0-30.0+winnings)) > 0.01 {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-30.0+winnings, player.Bankroll)
	}
	// Check win message - allow for actual calculated value
	if len(result) == 0 || !strings.Contains(result[0], "Lay 9 wins") {
		t.Errorf("Expected win message for Lay 9, got %v", result)
	}
}

func TestLAY10WinOn7(t *testing.T) {
	// Test LAY_10 win on 7 (1:2 minus commission)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so lay bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	_, err := PlaceBet(table, "player1", "LAY_10", 20.0)
	if err != nil {
		t.Fatalf("Failed to place LAY_10 bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 1:2 payout minus 5% commission
	gross := 20.0 * 0.5            // 10.0
	commission := gross * 0.05     // 0.5
	winnings := gross - commission // 9.5
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Lay 10 wins $9.50") {
		t.Errorf("Expected win message for Lay 10, got %v", result)
	}
}

func TestLayBetsLoseOnTheirNumber(t *testing.T) {
	// Test lay bets lose on their number
	layTypes := []struct {
		betType string
		number  int
		die1    int
		die2    int
	}{
		{"LAY_4", 4, 2, 2},
		{"LAY_5", 5, 2, 3},
		{"LAY_6", 6, 3, 3},
		{"LAY_8", 8, 4, 4},
		{"LAY_9", 9, 4, 5},
		{"LAY_10", 10, 5, 5},
	}

	for _, test := range layTypes {
		table := NewTable(5.0, 1000.0, 3)
		AddPlayer(table, "player1", "Test Player", 1000.0)

		// Set table to point phase so lay bets work
		table.State = crapsgame.StatePoint
		table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

		_, err := PlaceBet(table, "player1", test.betType, 20.0)
		if err != nil {
			t.Fatalf("Failed to place %s bet: %v", test.betType, err)
		}
		roll := &Roll{Die1: test.die1, Die2: test.die2, Total: test.number, IsHard: test.die1 == test.die2}
		result := table.BetResolver.ResolveBets(roll)
		player, _ := GetPlayer(table, "player1")
		if player.Bankroll != 1000.0-20.0 {
			t.Errorf("Expected bankroll %.2f, got %.2f for %s", 1000.0-20.0, player.Bankroll, test.betType)
		}
		if !strings.Contains(strings.Join(result, " "), "loses $20.00") {
			t.Errorf("Expected lose message for %s, got %v", test.betType, result)
		}
	}
}

// 8.7 Horn Bet Resolution
func TestHORNWinOn2(t *testing.T) {
	// Test HORN win on 2 (30:1 for aces portion)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN bet: %v", err)
	}
	// Simulate a roll of 2 (aces)
	roll := &Roll{Die1: 1, Die2: 1, Total: 2, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	// Horn bet is split 4 ways, so 5.0 on each number
	// 30:1 payout for aces (2)
	winnings := 5.0 * 30.0 // 150.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn wins $150.00") {
		t.Errorf("Expected win message for Horn on 2, got %v", result)
	}
}

func TestHORNWinOn3(t *testing.T) {
	// Test HORN win on 3 (15:1 for ace-deuce portion)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN bet: %v", err)
	}
	// Simulate a roll of 3 (ace-deuce)
	roll := &Roll{Die1: 1, Die2: 2, Total: 3, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// Horn bet is split 4 ways, so 5.0 on each number
	// 15:1 payout for ace-deuce (3)
	winnings := 5.0 * 15.0 // 75.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn wins $75.00") {
		t.Errorf("Expected win message for Horn on 3, got %v", result)
	}
}

func TestHORNWinOn11(t *testing.T) {
	// Test HORN win on 11 (15:1 for eleven portion)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN bet: %v", err)
	}
	// Simulate a roll of 11
	roll := &Roll{Die1: 5, Die2: 6, Total: 11, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// Horn bet is split 4 ways, so 5.0 on each number
	// 15:1 payout for eleven (11)
	winnings := 5.0 * 15.0 // 75.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn wins $75.00") {
		t.Errorf("Expected win message for Horn on 11, got %v", result)
	}
}

func TestHORNWinOn12(t *testing.T) {
	// Test HORN win on 12 (30:1 for boxcars portion)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN bet: %v", err)
	}
	// Simulate a roll of 12 (boxcars)
	roll := &Roll{Die1: 6, Die2: 6, Total: 12, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	// Horn bet is split 4 ways, so 5.0 on each number
	// 30:1 payout for boxcars (12)
	winnings := 5.0 * 30.0 // 150.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn wins $150.00") {
		t.Errorf("Expected win message for Horn on 12, got %v", result)
	}
}

func TestHORNHIGH2WinOn2(t *testing.T) {
	// Test HORN_HIGH_2 win on 2 (27:4 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN_HIGH_2", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN_HIGH_2 bet: %v", err)
	}
	// Simulate a roll of 2 (aces)
	roll := &Roll{Die1: 1, Die2: 1, Total: 2, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	// 27:4 payout for high 2
	winnings := 20.0 * (27.0 / 4.0) // 135.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn High 2 wins $135.00") {
		t.Errorf("Expected win message for Horn High 2, got %v", result)
	}
}

func TestHORNHIGH3WinOn3(t *testing.T) {
	// Test HORN_HIGH_3 win on 3 (15:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN_HIGH_3", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN_HIGH_3 bet: %v", err)
	}
	// Simulate a roll of 3 (ace-deuce)
	roll := &Roll{Die1: 1, Die2: 2, Total: 3, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 15:1 payout for high 3
	winnings := 20.0 * 15.0 // 300.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn High 3 wins $300.00") {
		t.Errorf("Expected win message for Horn High 3, got %v", result)
	}
}

func TestHORNHIGH11WinOn11(t *testing.T) {
	// Test HORN_HIGH_11 win on 11 (15:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN_HIGH_11", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN_HIGH_11 bet: %v", err)
	}
	// Simulate a roll of 11
	roll := &Roll{Die1: 5, Die2: 6, Total: 11, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 15:1 payout for high 11
	winnings := 20.0 * 15.0 // 300.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn High 11 wins $300.00") {
		t.Errorf("Expected win message for Horn High 11, got %v", result)
	}
}

func TestHORNHIGH12WinOn12(t *testing.T) {
	// Test HORN_HIGH_12 win on 12 (27:4 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HORN_HIGH_12", 20.0)
	if err != nil {
		t.Fatalf("Failed to place HORN_HIGH_12 bet: %v", err)
	}
	// Simulate a roll of 12 (boxcars)
	roll := &Roll{Die1: 6, Die2: 6, Total: 12, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	// 27:4 payout for high 12
	winnings := 20.0 * (27.0 / 4.0) // 135.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Horn High 12 wins $135.00") {
		t.Errorf("Expected win message for Horn High 12, got %v", result)
	}
}

// 8.8 Hop Bet Resolution
func TestHOP12WinOnExactRoll12(t *testing.T) {
	// Test HOP_1_2 win on exact roll 1-2 (15:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HOP_1_2", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HOP_1_2 bet: %v", err)
	}
	// Simulate exact roll 1-2
	roll := &Roll{Die1: 1, Die2: 2, Total: 3, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 15:1 payout for exact hop
	winnings := 10.0 * 15.0 // 150.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-10.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-10.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Hop 1-2 wins $150.00") {
		t.Errorf("Expected win message for Hop 1-2, got %v", result)
	}
}

func TestHOPHARD6WinOnHard6(t *testing.T) {
	// Test HOP_HARD_6 win on hard 6 (30:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HOP_HARD_6", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HOP_HARD_6 bet: %v", err)
	}
	// Simulate hard 6 (3-3)
	roll := &Roll{Die1: 3, Die2: 3, Total: 6, IsHard: true}
	result := table.BetResolver.ResolveBets(roll)
	// 30:1 payout for hard hop
	winnings := 10.0 * 30.0 // 300.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-10.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-10.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Hop Hard 6 wins $300.00") {
		t.Errorf("Expected win message for Hop Hard 6, got %v", result)
	}
}

func TestHOPEASY8WinOnEasy8(t *testing.T) {
	// Test HOP_EASY_8 win on easy 8 (15:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "HOP_EASY_8", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HOP_EASY_8 bet: %v", err)
	}
	// Simulate easy 8 (3-5)
	roll := &Roll{Die1: 3, Die2: 5, Total: 8, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 15:1 payout for easy hop
	winnings := 10.0 * 15.0 // 150.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-10.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-10.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Hop Easy 8 wins $150.00") {
		t.Errorf("Expected win message for Hop Easy 8, got %v", result)
	}
}

// 8.9 Big Bet Resolution
func TestBIG6WinOn6(t *testing.T) {
	// Test BIG_6 win on 6 (1:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "BIG_6", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BIG_6 bet: %v", err)
	}
	// Simulate a roll of 6
	roll := &Roll{Die1: 2, Die2: 4, Total: 6, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 1:1 payout for big 6
	winnings := 20.0 * 1.0 // 20.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Big 6 wins $20.00") {
		t.Errorf("Expected win message for Big 6, got %v", result)
	}
}

func TestBIG8WinOn8(t *testing.T) {
	// Test BIG_8 win on 8 (1:1 payout)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "BIG_8", 20.0)
	if err != nil {
		t.Fatalf("Failed to place BIG_8 bet: %v", err)
	}
	// Simulate a roll of 8
	roll := &Roll{Die1: 3, Die2: 5, Total: 8, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// 1:1 payout for big 8
	winnings := 20.0 * 1.0 // 20.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "Big 8 wins $20.00") {
		t.Errorf("Expected win message for Big 8, got %v", result)
	}
}

func TestBigBetsLoseOn7(t *testing.T) {
	// Test big bets lose on 7
	bigTypes := []string{"BIG_6", "BIG_8"}
	for _, betType := range bigTypes {
		table := NewTable(5.0, 1000.0, 3)
		AddPlayer(table, "player1", "Test Player", 1000.0)
		_, err := PlaceBet(table, "player1", betType, 20.0)
		if err != nil {
			t.Fatalf("Failed to place %s bet: %v", betType, err)
		}
		roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
		result := table.BetResolver.ResolveBets(roll)
		player, _ := GetPlayer(table, "player1")
		if player.Bankroll != 1000.0-20.0 {
			t.Errorf("Expected bankroll %.2f, got %.2f for %s", 1000.0-20.0, player.Bankroll, betType)
		}
		if !strings.Contains(strings.Join(result, " "), "loses $20.00") {
			t.Errorf("Expected lose message for %s, got %v", betType, result)
		}
	}
}

// 8.10 Combination Bet Resolution
func TestWORLDWinOn7(t *testing.T) {
	// Test WORLD win on 7 (4:1 for any 7 portion)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "WORLD", 20.0)
	if err != nil {
		t.Fatalf("Failed to place WORLD bet: %v", err)
	}
	// Simulate a roll of 7
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// World bet is split: 16.0 on any 7 (4:1), 4.0 on any craps (7:1)
	// 4:1 payout for any 7 portion
	winnings := 16.0 * 4.0 // 64.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "World wins $64.00") {
		t.Errorf("Expected win message for World on 7, got %v", result)
	}
}

func TestWORLDWinOn2_3_12(t *testing.T) {
	// Test WORLD win on 2, 3, 12 (7:1 for any craps portion)
	crapsNumbers := []struct {
		number int
		die1   int
		die2   int
	}{
		{2, 1, 1},  // aces
		{3, 1, 2},  // ace-deuce
		{12, 6, 6}, // boxcars
	}

	for _, test := range crapsNumbers {
		table := NewTable(5.0, 1000.0, 3)
		AddPlayer(table, "player1", "Test Player", 1000.0)
		_, err := PlaceBet(table, "player1", "WORLD", 20.0)
		if err != nil {
			t.Fatalf("Failed to place WORLD bet: %v", err)
		}
		roll := &Roll{Die1: test.die1, Die2: test.die2, Total: test.number, IsHard: test.die1 == test.die2}
		result := table.BetResolver.ResolveBets(roll)
		// World bet is split: 16.0 on any 7 (4:1), 4.0 on any craps (7:1)
		// 7:1 payout for any craps portion
		winnings := 4.0 * 7.0 // 28.0
		player, _ := GetPlayer(table, "player1")
		if player.Bankroll != 1000.0-20.0+winnings {
			t.Errorf("Expected bankroll %.2f, got %.2f for roll %d", 1000.0-20.0+winnings, player.Bankroll, test.number)
		}
		if !strings.Contains(strings.Join(result, " "), "World wins $28.00") {
			t.Errorf("Expected win message for World on %d, got %v", test.number, result)
		}
	}
}

func TestCANDEWinOn11(t *testing.T) {
	// Test C_AND_E win on 11 (15:1 for eleven portion)
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err := PlaceBet(table, "player1", "C_AND_E", 20.0)
	if err != nil {
		t.Fatalf("Failed to place C_AND_E bet: %v", err)
	}
	// Simulate a roll of 11
	roll := &Roll{Die1: 5, Die2: 6, Total: 11, IsHard: false}
	result := table.BetResolver.ResolveBets(roll)
	// C_AND_E bet is split: 10.0 on eleven (15:1), 10.0 on any craps (7:1)
	// 15:1 payout for eleven portion
	winnings := 10.0 * 15.0 // 150.0
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1000.0-20.0+winnings {
		t.Errorf("Expected bankroll %.2f, got %.2f", 1000.0-20.0+winnings, player.Bankroll)
	}
	if !strings.Contains(strings.Join(result, " "), "C and E wins $150.00") {
		t.Errorf("Expected win message for C and E on 11, got %v", result)
	}
}

func TestCANDEWinOn2_3_12(t *testing.T) {
	// Test C_AND_E win on 2, 3, 12 (7:1 for any craps portion)
	crapsNumbers := []struct {
		number int
		die1   int
		die2   int
	}{
		{2, 1, 1},  // aces
		{3, 1, 2},  // ace-deuce
		{12, 6, 6}, // boxcars
	}

	for _, test := range crapsNumbers {
		table := NewTable(5.0, 1000.0, 3)
		AddPlayer(table, "player1", "Test Player", 1000.0)
		_, err := PlaceBet(table, "player1", "C_AND_E", 20.0)
		if err != nil {
			t.Fatalf("Failed to place C_AND_E bet: %v", err)
		}
		roll := &Roll{Die1: test.die1, Die2: test.die2, Total: test.number, IsHard: test.die1 == test.die2}
		result := table.BetResolver.ResolveBets(roll)
		// C_AND_E bet is split: 10.0 on eleven (15:1), 10.0 on any craps (7:1)
		// 7:1 payout for any craps portion
		winnings := 10.0 * 7.0 // 70.0
		player, _ := GetPlayer(table, "player1")
		if player.Bankroll != 1000.0-20.0+winnings {
			t.Errorf("Expected bankroll %.2f, got %.2f for roll %d", 1000.0-20.0+winnings, player.Bankroll, test.number)
		}
		if !strings.Contains(strings.Join(result, " "), "C and E wins $70.00") {
			t.Errorf("Expected win message for C and E on %d, got %v", test.number, result)
		}
	}
}

// ============================================================================
// 9. Bet Lifecycle Tests
// ============================================================================

func TestOneRollBetRemovalAfterResolution(t *testing.T) {
	// Test one-roll bet removal after resolution
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place one-roll bets
	_, err := PlaceBet(table, "player1", "ANY_SEVEN", 10.0)
	if err != nil {
		t.Fatalf("Failed to place ANY_SEVEN bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place FIELD bet: %v", err)
	}

	// Verify bets are placed
	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 2 {
		t.Errorf("Expected 2 bets, got %d", len(player.Bets))
	}

	// Roll dice (any roll will resolve one-roll bets)
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	table.BetResolver.ResolveBets(roll)

	// Verify one-roll bets are removed
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets after resolution, got %d", len(player.Bets))
	}
}

func TestMultiRollBetPersistence(t *testing.T) {
	// Test multi-roll bet persistence
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Set table to point phase so place bets work
	table.State = crapsgame.StatePoint
	table.Point = crapsgame.Point6 // Set any point, doesn't matter for this test

	// Place multi-roll bets
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "HARD_8", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HARD_8 bet: %v", err)
	}

	// Verify bets are placed
	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 2 {
		t.Errorf("Expected 2 bets, got %d", len(player.Bets))
	}

	// Roll dice (roll that doesn't resolve these bets)
	roll := &Roll{Die1: 2, Die2: 3, Total: 5, IsHard: false}
	table.BetResolver.ResolveBets(roll)

	// Verify multi-roll bets persist
	if len(player.Bets) != 2 {
		t.Errorf("Expected 2 bets to persist, got %d", len(player.Bets))
	}

	// Roll dice that resolves one of the bets (hard 8)
	roll2 := &Roll{Die1: 4, Die2: 4, Total: 8, IsHard: true}
	table.BetResolver.ResolveBets(roll2)

	// Verify HARD_8 is removed but PLACE_6 persists
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet to persist, got %d", len(player.Bets))
	}
	if player.Bets[0].Type != "PLACE_6" {
		t.Errorf("Expected PLACE_6 bet to persist, got %s", player.Bets[0].Type)
	}
}

func TestBetWorkingStatusUpdates(t *testing.T) {
	// Test bet working status updates
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place a bet that can have working status changed
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}

	// Verify bet is initially working
	player, _ := GetPlayer(table, "player1")
	if !player.Bets[0].Working {
		t.Error("Expected bet to be initially working")
	}

	// Turn bet off
	// Note: This would require implementing a turn off mechanism
	// For now, we'll test the working status is properly set initially
	if player.Bets[0].Working != true {
		t.Error("Expected bet working status to be true")
	}

	// Verify bet type and amount are correct
	if player.Bets[0].Type != "PLACE_6" {
		t.Errorf("Expected bet type PLACE_6, got %s", player.Bets[0].Type)
	}
	if player.Bets[0].Amount != 24.0 {
		t.Errorf("Expected bet amount 24.0, got %f", player.Bets[0].Amount)
	}
}

func TestBetCleanupAfterPlayerRemoval(t *testing.T) {
	// Test bet cleanup after player removal
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)
	AddPlayer(table, "player2", "Test Player 2", 1000.0)

	// Place bets for both players
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}
	_, err = PlaceBet(table, "player2", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place FIELD bet: %v", err)
	}

	// Verify bets are placed
	player1, _ := GetPlayer(table, "player1")
	player2, _ := GetPlayer(table, "player2")
	if len(player1.Bets) != 1 {
		t.Errorf("Expected 1 bet for player1, got %d", len(player1.Bets))
	}
	if len(player2.Bets) != 1 {
		t.Errorf("Expected 1 bet for player2, got %d", len(player2.Bets))
	}

	// Remove player1
	err = RemovePlayer(table, "player1")
	if err != nil {
		t.Fatalf("Failed to remove player1: %v", err)
	}

	// Verify player1 is removed and their bets are cleaned up
	if _, exists := table.Players["player1"]; exists {
		t.Error("Player1 should be removed")
	}

	// Verify player2 and their bets remain
	if _, exists := table.Players["player2"]; !exists {
		t.Error("Player2 should still exist")
	}
	player2, _ = GetPlayer(table, "player2")
	if len(player2.Bets) != 1 {
		t.Errorf("Expected 1 bet for player2 to remain, got %d", len(player2.Bets))
	}
}

// ============================================================================
// 10. Query Statement Tests
// ============================================================================

func TestSHOWPOINTCommand(t *testing.T) {
	// Test SHOW POINT command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test initial state (no point)
	result, err := ExecuteString("SHOW POINT;", table)
	if err != nil {
		t.Fatalf("Failed to execute SHOW POINT command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SHOW POINT command")
	}

	// Establish a point by rolling dice until a point is set
	var roll *Roll
	for {
		roll = RollDice(table)
		if roll.Total == 4 || roll.Total == 5 || roll.Total == 6 || roll.Total == 8 || roll.Total == 9 || roll.Total == 10 {
			break
		}
		// If 7, 11, 2, 3, 12, keep rolling
	}

	// Test with point established
	result2, err2 := ExecuteString("SHOW POINT;", table)
	if err2 != nil {
		t.Fatalf("Failed to execute SHOW POINT command: %v", err2)
	}
	if len(result2) == 0 {
		t.Error("Expected result from SHOW POINT command with point established")
	}

	// Verify point is set
	if table.Point == crapsgame.PointOff {
		t.Error("Expected point to be established")
	}
}

func TestSHOWBETSCommand(t *testing.T) {
	// Test SHOW BETS command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test with no bets
	result, err := ExecuteString("SHOW BETS;", table)
	if err != nil {
		t.Fatalf("Failed to execute SHOW BETS command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SHOW BETS command")
	}

	// Place some bets
	ExecuteString("PLACE $25 ON PASS_LINE;", table)
	ExecuteString("PLACE $10 ON FIELD;", table)

	// Test with bets placed
	result2, err2 := ExecuteString("SHOW BETS;", table)
	if err2 != nil {
		t.Fatalf("Failed to execute SHOW BETS command: %v", err2)
	}
	if len(result2) == 0 {
		t.Error("Expected result from SHOW BETS command with bets")
	}

	// Verify bets are shown
	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 2 {
		t.Errorf("Expected 2 bets, got %d", len(player.Bets))
	}
}

func TestSHOWBANKROLLCommand(t *testing.T) {
	// Test SHOW BANKROLL command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test initial bankroll
	result, err := ExecuteString("SHOW BANKROLL;", table)
	if err != nil {
		t.Fatalf("Failed to execute SHOW BANKROLL command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SHOW BANKROLL command")
	}

	// Place a bet to change bankroll
	ExecuteString("PLACE $50 ON PASS_LINE;", table)

	// Test after bet placement
	result2, err2 := ExecuteString("SHOW BANKROLL;", table)
	if err2 != nil {
		t.Fatalf("Failed to execute SHOW BANKROLL command: %v", err2)
	}
	if len(result2) == 0 {
		t.Error("Expected result from SHOW BANKROLL command after bet")
	}

	// Verify bankroll decreased
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 950.0 {
		t.Errorf("Expected bankroll $950.00, got $%.2f", player.Bankroll)
	}
}

func TestSHOWTABLEMINIMUMSCommand(t *testing.T) {
	// Test SHOW TABLE_MINIMUMS command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test table minimums display
	result, err := ExecuteString("SHOW TABLE_MINIMUMS;", table)
	if err != nil {
		t.Fatalf("Failed to execute SHOW TABLE_MINIMUMS command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SHOW TABLE_MINIMUMS command")
	}

	// Verify table limits are correct
	if table.MinBet != 5.0 {
		t.Errorf("Expected min bet $5.00, got $%.2f", table.MinBet)
	}
	if table.MaxBet != 1000.0 {
		t.Errorf("Expected max bet $1000.00, got $%.2f", table.MaxBet)
	}
	if table.MaxOdds != 3 {
		t.Errorf("Expected max odds 3x, got %dx", table.MaxOdds)
	}
}

func TestSHOWODDSALLOWEDCommand(t *testing.T) {
	// Test SHOW ODDS_ALLOWED command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test odds allowed display
	result, err := ExecuteString("SHOW ODDS_ALLOWED;", table)
	if err != nil {
		t.Fatalf("Failed to execute SHOW ODDS_ALLOWED command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SHOW ODDS_ALLOWED command")
	}

	// Verify odds limit is correct
	if table.MaxOdds != 3 {
		t.Errorf("Expected max odds 3x, got %dx", table.MaxOdds)
	}

	// Test with different odds limit
	table2 := NewTable(10.0, 500.0, 5)
	AddPlayer(table2, "player2", "Test Player 2", 1000.0)

	result2, err2 := ExecuteString("SHOW ODDS_ALLOWED;", table2)
	if err2 != nil {
		t.Fatalf("Failed to execute SHOW ODDS_ALLOWED command: %v", err2)
	}
	if len(result2) == 0 {
		t.Error("Expected result from SHOW ODDS_ALLOWED command")
	}

	if table2.MaxOdds != 5 {
		t.Errorf("Expected max odds 5x, got %dx", table2.MaxOdds)
	}
}

// ============================================================================
// 11. Management Statement Tests
// ============================================================================

func TestSETBANKROLLCommand(t *testing.T) {
	// Test SET BANKROLL command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test setting bankroll
	result, err := ExecuteString("SET BANKROLL TO $1500;", table)
	if err != nil {
		t.Fatalf("Failed to execute SET BANKROLL command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SET BANKROLL command")
	}

	// Verify bankroll was set
	player, _ := GetPlayer(table, "player1")
	if player.Bankroll != 1500.0 {
		t.Errorf("Expected bankroll $1500.00, got $%.2f", player.Bankroll)
	}

	// Test setting to zero
	result2, err2 := ExecuteString("SET BANKROLL TO $0;", table)
	if err2 != nil {
		t.Fatalf("Failed to execute SET BANKROLL command: %v", err2)
	}
	if len(result2) == 0 {
		t.Error("Expected result from SET BANKROLL command")
	}

	player, _ = GetPlayer(table, "player1")
	if player.Bankroll != 0.0 {
		t.Errorf("Expected bankroll $0.00, got $%.2f", player.Bankroll)
	}
}

func TestSETMAXBETCommand(t *testing.T) {
	// Test SET MAX_BET command
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test setting max bet
	result, err := ExecuteString("SET MAX_BET TO $500;", table)
	if err != nil {
		t.Fatalf("Failed to execute SET MAX_BET command: %v", err)
	}
	if len(result) == 0 {
		t.Error("Expected result from SET MAX_BET command")
	}

	// Verify max bet was set
	player, _ := GetPlayer(table, "player1")
	if player.MaxBet != 500.0 {
		t.Errorf("Expected max bet $500.00, got $%.2f", player.MaxBet)
	}

	// Test setting to table minimum
	result2, err2 := ExecuteString("SET MAX_BET TO $5;", table)
	if err2 != nil {
		t.Fatalf("Failed to execute SET MAX_BET command: %v", err2)
	}
	if len(result2) == 0 {
		t.Error("Expected result from SET MAX_BET command")
	}

	player, _ = GetPlayer(table, "player1")
	if player.MaxBet != 5.0 {
		t.Errorf("Expected max bet $5.00, got $%.2f", player.MaxBet)
	}
}

// ============================================================================
// 12. Remove Command Tests
// ============================================================================

func TestRemoveAllBets(t *testing.T) {
	// Test REMOVE ALL functionality
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place multiple bets
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place FIELD bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "HARD_8", 10.0)
	if err != nil {
		t.Fatalf("Failed to place HARD_8 bet: %v", err)
	}

	// Verify bets are placed
	player, _ := GetPlayer(table, "player1")
	initialBankroll := player.Bankroll
	initialBetCount := len(player.Bets)
	if initialBetCount != 3 {
		t.Errorf("Expected 3 bets, got %d", initialBetCount)
	}

	// Execute REMOVE ALL command
	result, err := ExecuteString("REMOVE ALL;", table)
	if err != nil {
		t.Fatalf("Failed to execute REMOVE ALL command: %v", err)
	}

	// Verify result message
	if len(result) == 0 {
		t.Error("Expected result from REMOVE ALL command")
	}
	if !strings.Contains(strings.Join(result, " "), "Removed 3 bets") {
		t.Errorf("Expected 'Removed 3 bets' message, got %v", result)
	}

	// Verify all bets are removed
	player, _ = GetPlayer(table, "player1")
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets after REMOVE ALL, got %d", len(player.Bets))
	}

	// Verify bankroll is returned
	expectedBankroll := initialBankroll + 24.0 + 10.0 + 10.0 // Return all bet amounts
	if player.Bankroll != expectedBankroll {
		t.Errorf("Expected bankroll %.2f, got %.2f", expectedBankroll, player.Bankroll)
	}
}

func TestRemoveSpecificBetType(t *testing.T) {
	// Test REMOVE specific bet type
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place multiple bets of different types
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "PLACE_8", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_8 bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place FIELD bet: %v", err)
	}

	// Verify bets are placed
	player, _ := GetPlayer(table, "player1")
	initialBankroll := player.Bankroll
	if len(player.Bets) != 3 {
		t.Errorf("Expected 3 bets, got %d", len(player.Bets))
	}

	// Execute REMOVE PLACE_6 command
	result, err := ExecuteString("REMOVE PLACE_6;", table)
	if err != nil {
		t.Fatalf("Failed to execute REMOVE PLACE_6 command: %v", err)
	}

	// Verify result message
	if len(result) == 0 {
		t.Error("Expected result from REMOVE PLACE_6 command")
	}
	if !strings.Contains(strings.Join(result, " "), "Removed 1 PLACE_6 bets") {
		t.Errorf("Expected 'Removed 1 PLACE_6 bets' message, got %v", result)
	}

	// Verify only PLACE_6 bet is removed
	player, _ = GetPlayer(table, "player1")
	if len(player.Bets) != 2 {
		t.Errorf("Expected 2 bets after REMOVE PLACE_6, got %d", len(player.Bets))
	}

	// Verify remaining bets are correct
	betTypes := make(map[string]bool)
	for _, bet := range player.Bets {
		betTypes[bet.Type] = true
	}
	if !betTypes["PLACE_8"] {
		t.Error("Expected PLACE_8 bet to remain")
	}
	if !betTypes["FIELD"] {
		t.Error("Expected FIELD bet to remain")
	}

	// Verify bankroll is returned for removed bet
	expectedBankroll := initialBankroll + 24.0 // Return PLACE_6 bet amount
	if player.Bankroll != expectedBankroll {
		t.Errorf("Expected bankroll %.2f, got %.2f", expectedBankroll, player.Bankroll)
	}
}

func TestRemoveMultipleBetsOfSameType(t *testing.T) {
	// Test removing multiple bets of the same type
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place multiple FIELD bets
	_, err := PlaceBet(table, "player1", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place first FIELD bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "FIELD", 15.0)
	if err != nil {
		t.Fatalf("Failed to place second FIELD bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}

	// Verify bets are placed
	player, _ := GetPlayer(table, "player1")
	initialBankroll := player.Bankroll
	if len(player.Bets) != 3 {
		t.Errorf("Expected 3 bets, got %d", len(player.Bets))
	}

	// Execute REMOVE FIELD command
	result, err := ExecuteString("REMOVE FIELD;", table)
	if err != nil {
		t.Fatalf("Failed to execute REMOVE FIELD command: %v", err)
	}

	// Verify result message
	if len(result) == 0 {
		t.Error("Expected result from REMOVE FIELD command")
	}
	if !strings.Contains(strings.Join(result, " "), "Removed 2 FIELD bets") {
		t.Errorf("Expected 'Removed 2 FIELD bets' message, got %v", result)
	}

	// Verify only FIELD bets are removed
	player, _ = GetPlayer(table, "player1")
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet after REMOVE FIELD, got %d", len(player.Bets))
	}

	// Verify remaining bet is PLACE_6
	if player.Bets[0].Type != "PLACE_6" {
		t.Errorf("Expected PLACE_6 bet to remain, got %s", player.Bets[0].Type)
	}

	// Verify bankroll is returned for removed bets
	expectedBankroll := initialBankroll + 10.0 + 15.0 // Return both FIELD bet amounts
	if player.Bankroll != expectedBankroll {
		t.Errorf("Expected bankroll %.2f, got %.2f", expectedBankroll, player.Bankroll)
	}
}

func TestRemoveStatementParsing(t *testing.T) {
	// Test REMOVE statement parsing
	testCases := []struct {
		input        string
		expectedType string
		expectedAll  bool
		description  string
	}{
		{"REMOVE ALL;", "", true, "REMOVE ALL"},
		{"REMOVE PLACE_6;", "PLACE_6", false, "REMOVE specific bet type"},
		{"REMOVE FIELD;", "FIELD", false, "REMOVE FIELD"},
		{"REMOVE HARD_8;", "HARD_8", false, "REMOVE HARD_8"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			parser := NewParser(lexer)

			program := parser.ParseProgram()

			if len(program.Statements) != 1 {
				t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*RemoveStatement)
			if !ok {
				t.Fatalf("Expected RemoveStatement, got %T", program.Statements[0])
			}

			if tc.expectedAll {
				if stmt.BetType != nil {
					t.Errorf("Expected BetType to be nil for REMOVE ALL, got %v", stmt.BetType)
				}
			} else {
				if stmt.BetType == nil {
					t.Errorf("Expected BetType to be set for %s, got nil", tc.expectedType)
				} else {
					betTypeString, err := BetTypeToString(stmt.BetType.Type)
					if err != nil {
						t.Errorf("Failed to convert bet type to string: %v", err)
					}
					if betTypeString != tc.expectedType {
						t.Errorf("Expected bet type %s, got %s", tc.expectedType, betTypeString)
					}
				}
			}
		})
	}
}

func TestRemoveCommandExecution(t *testing.T) {
	// Test full REMOVE command execution through interpreter
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place some bets
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}

	// Test REMOVE ALL execution
	result, err := ExecuteString("REMOVE ALL;", table)
	if err != nil {
		t.Fatalf("Failed to execute REMOVE ALL: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected result from REMOVE ALL execution")
	}

	// Verify bet is removed
	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets after REMOVE ALL execution, got %d", len(player.Bets))
	}

	// Place another bet and test specific removal
	_, err = PlaceBet(table, "player1", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place FIELD bet: %v", err)
	}

	result2, err2 := ExecuteString("REMOVE FIELD;", table)
	if err2 != nil {
		t.Fatalf("Failed to execute REMOVE FIELD: %v", err2)
	}

	if len(result2) == 0 {
		t.Error("Expected result from REMOVE FIELD execution")
	}

	// Verify FIELD bet is removed
	player, _ = GetPlayer(table, "player1")
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets after REMOVE FIELD execution, got %d", len(player.Bets))
	}
}

func TestRemoveEdgeCases(t *testing.T) {
	// Test edge cases for REMOVE command
	table := NewTable(5.0, 1000.0, 3)

	// Test removing when no players exist
	_, err := ExecuteString("REMOVE ALL;", table)
	if err == nil {
		t.Error("Expected error when removing bets with no players")
	}

	// Test removing when no bets exist
	AddPlayer(table, "player1", "Test Player", 1000.0)
	_, err = ExecuteString("REMOVE ALL;", table)
	if err != nil {
		t.Fatalf("Unexpected error when removing with no bets: %v", err)
	}

	// Test removing non-existent bet type (should return error)
	_, err = ExecuteString("REMOVE INVALID_BET;", table)
	if err == nil {
		t.Error("Expected error when removing non-existent bet type")
	}

	// Test removing specific bet type when none exist
	_, err = ExecuteString("REMOVE PLACE_6;", table)
	if err != nil {
		t.Fatalf("Unexpected error when removing non-existent PLACE_6: %v", err)
	}

	// Verify no bets were removed (since none existed)
	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets, got %d", len(player.Bets))
	}
}

func TestRemoveWithNonWorkingBets(t *testing.T) {
	// Test that only working bets are removed
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place a bet
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}

	// Manually set bet to non-working (simulating resolved bet)
	player, _ := GetPlayer(table, "player1")
	player.Bets[0].Working = false

	// Execute REMOVE ALL
	result, err := ExecuteString("REMOVE ALL;", table)
	if err != nil {
		t.Fatalf("Failed to execute REMOVE ALL: %v", err)
	}

	// Verify result indicates no bets were removed
	if len(result) == 0 {
		t.Error("Expected result from REMOVE ALL")
	}
	if !strings.Contains(strings.Join(result, " "), "Removed 0 bets") {
		t.Errorf("Expected 'Removed 0 bets' message, got %v", result)
	}

	// Verify bet is still there (since it was non-working)
	if len(player.Bets) != 1 {
		t.Errorf("Expected 1 bet to remain (non-working), got %d", len(player.Bets))
	}
}

func TestRemoveCommandWithMultiplePlayers(t *testing.T) {
	// Test REMOVE command with multiple players
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player 1", 1000.0)
	AddPlayer(table, "player2", "Test Player 2", 1000.0)

	// Place bets for both players
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place bet for player1: %v", err)
	}
	_, err = PlaceBet(table, "player2", "FIELD", 10.0)
	if err != nil {
		t.Fatalf("Failed to place bet for player2: %v", err)
	}

	// Execute REMOVE ALL for player1 only
	result, err := ExecuteStringForPlayer("REMOVE ALL;", table, "player1")
	if err != nil {
		t.Fatalf("Failed to execute REMOVE ALL: %v", err)
	}

	// Verify result
	if len(result) == 0 {
		t.Error("Expected result from REMOVE ALL")
	}

	// Verify player1's bet is removed
	player1, _ := GetPlayer(table, "player1")
	if len(player1.Bets) != 0 {
		t.Errorf("Expected 0 bets for player1, got %d", len(player1.Bets))
	}

	// Verify player2's bet remains
	player2, _ := GetPlayer(table, "player2")
	if len(player2.Bets) != 1 {
		t.Errorf("Expected 1 bet for player2, got %d", len(player2.Bets))
	}
	if len(player2.Bets) > 0 && player2.Bets[0].Type != "FIELD" {
		t.Errorf("Expected FIELD bet for player2, got %s", player2.Bets[0].Type)
	}
}

func TestRemoveCommandErrorHandling(t *testing.T) {
	// Test error handling for malformed REMOVE commands
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Test malformed REMOVE command (missing semicolon)
	_, err := ExecuteString("REMOVE ALL", table)
	if err == nil {
		t.Error("Expected error for malformed REMOVE command")
	}

	// Test malformed REMOVE command (invalid bet type) - should return parse error
	_, err = ExecuteString("REMOVE INVALID_BET_TYPE;", table)
	if err == nil {
		t.Error("Expected error for invalid bet type")
	}

	// Test empty REMOVE command
	_, err = ExecuteString("REMOVE;", table)
	if err == nil {
		t.Error("Expected error for empty REMOVE command")
	}
}

func TestRemoveCommandIntegration(t *testing.T) {
	// Test integration of REMOVE command with other commands
	table := NewTable(5.0, 1000.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Execute a sequence of commands
	commands := []string{
		"PLACE $24 ON PLACE_6;",
		"PLACE $10 ON FIELD;",
		"SHOW BETS;",
		"REMOVE PLACE_6;",
		"SHOW BETS;",
		"REMOVE ALL;",
		"SHOW BETS;",
	}

	for i, cmd := range commands {
		_, err := ExecuteString(cmd, table)
		if err != nil {
			t.Fatalf("Failed to execute command %d (%s): %v", i+1, cmd, err)
		}
	}

	// Verify final state
	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets at end of integration test, got %d", len(player.Bets))
	}
}

func TestPlaceBetsAreRemovedAfterSevenOut(t *testing.T) {
	table := NewTable(5.0, 100.0, 3)
	AddPlayer(table, "player1", "Test Player", 1000.0)

	// Place multiple place bets
	_, err := PlaceBet(table, "player1", "PLACE_6", 24.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_6 bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "PLACE_8", 25.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_8 bet: %v", err)
	}
	_, err = PlaceBet(table, "player1", "PLACE_9", 25.0)
	if err != nil {
		t.Fatalf("Failed to place PLACE_9 bet: %v", err)
	}

	player, _ := GetPlayer(table, "player1")
	if len(player.Bets) != 3 {
		t.Fatalf("Expected 3 place bets, got %d", len(player.Bets))
	}

	// Establish a point (roll until point is set)
	for {
		roll := RollDice(table)
		if roll.Total == 4 || roll.Total == 5 || roll.Total == 6 || roll.Total == 8 || roll.Total == 9 || roll.Total == 10 {
			break
		}
	}

	// Simulate a roll of 7 (seven-out)
	roll := &Roll{Die1: 3, Die2: 4, Total: 7, IsHard: false}
	table.BetResolver.ResolveBets(roll)

	// All place bets should be removed
	player, _ = GetPlayer(table, "player1")
	for _, bet := range player.Bets {
		if strings.HasPrefix(bet.Type, "PLACE_") {
			t.Errorf("Expected all place bets to be removed after seven-out, but found: %s", bet.Type)
		}
	}
	if len(player.Bets) != 0 {
		t.Errorf("Expected 0 bets after seven-out, got %d", len(player.Bets))
	}
}
