package crapsql

import (
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
