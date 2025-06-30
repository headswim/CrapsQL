package crapsql

import (
	"fmt"
	"strconv"
)

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.infixParseFns = make(map[TokenType]infixParseFn)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case PLACE:
		return p.parseBetStatement()
	case IF:
		return p.parseConditionalStatement()
	case SHOW:
		return p.parseQueryStatement()
	case SET:
		return p.parseManagementStatement()
	case REMOVE:
		return p.parseRemoveStatement()
	case PRESS:
		return p.parsePressStatement()
	case TURN:
		return p.parseTurnStatement()
	case ROLL:
		return p.parseRollStatement()
	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.curToken.Literal))
		// Use error recovery to skip to next statement
		return recoverFromParseError(p)
	}
}

func (p *Parser) parseBetStatement() *BetStatement {
	stmt := &BetStatement{Token: p.curToken}

	if !p.expectPeek(DOLLAR) {
		return nil
	}

	if !p.expectPeek(NUMBER) {
		return nil
	}

	// Parse amount
	amount := &AmountExpression{Token: p.curToken}
	val, err := parseAmount(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid amount: %s", p.curToken.Literal))
		return nil
	}
	amount.Value = val
	stmt.Amount = amount

	if !p.expectPeek(ON) {
		return nil
	}
	p.nextToken() // advance to bet type

	// Parse bet type
	stmt.BetType = p.parseBetTypeExpression()

	p.nextToken() // advance to next token after bet type

	// Parse optional modifiers (WITH or direct modifiers)
	modifiers := []*ModifierExpression{}
	if p.curToken.Type == WITH {
		p.nextToken() // advance to first modifier
		modifiers = append(modifiers, p.parseModifiers()...)
	} else {
		// Check for direct modifiers (e.g., WORKING, OFF, etc.)
		for isModifierToken(p.curToken.Type) {
			modifiers = append(modifiers, p.parseModifiers()...)
		}
	}
	stmt.Modifiers = modifiers

	if p.curToken.Type != SEMICOLON {
		p.addError("expected semicolon after bet statement")
		return nil
	}

	return stmt
}

// Helper to check if a token is a modifier
func isModifierToken(t TokenType) bool {
	switch t {
	case WORKING_KEYWORD, OFF_MODIFIER, PRESS, ODDS, ONE_ROLL, MAX, AMOUNT, RATIO:
		return true
	default:
		return false
	}
}

func (p *Parser) parseBetTypeExpression() *BetTypeExpression {
	expr := &BetTypeExpression{Token: p.curToken}

	// Map token type to bet type
	switch p.curToken.Type {
	case PASS_LINE:
		expr.Type = BetPassLine
	case DONT_PASS:
		expr.Type = BetDontPass
	case COME:
		expr.Type = BetCome
	case DONT_COME:
		expr.Type = BetDontCome
	case FIELD:
		expr.Type = BetField
	case ANY_SEVEN:
		expr.Type = BetAnySeven
	case ANY_CRAPS:
		expr.Type = BetAnyCraps
	case ELEVEN:
		expr.Type = BetEleven
	case ACE_DEUCE:
		expr.Type = BetAceDeuce
	case ACES:
		expr.Type = BetAces
	case BOXCARS:
		expr.Type = BetBoxcars
	case PLACE_4:
		expr.Type = BetPlace4
	case PLACE_5:
		expr.Type = BetPlace5
	case PLACE_6:
		expr.Type = BetPlace6
	case PLACE_8:
		expr.Type = BetPlace8
	case PLACE_9:
		expr.Type = BetPlace9
	case PLACE_10:
		expr.Type = BetPlace10
	case PLACE_NUMBERS:
		expr.Type = BetPlaceNumbers
		expr.Args = p.parsePlaceNumbers()
	case PLACE_INSIDE:
		expr.Type = BetPlaceInside
	case PLACE_OUTSIDE:
		expr.Type = BetPlaceOutside
	case HARD_4:
		expr.Type = BetHard4
	case HARD_6:
		expr.Type = BetHard6
	case HARD_8:
		expr.Type = BetHard8
	case HARD_10:
		expr.Type = BetHard10
	case ALL_HARDWAYS:
		expr.Type = BetAllHardways
	case PASS_ODDS:
		expr.Type = BetPassOdds
	case DONT_PASS_ODDS:
		expr.Type = BetDontPassOdds
	case BUY_4:
		expr.Type = BetBuy4
	case BUY_10:
		expr.Type = BetBuy10
	case LAY_4:
		expr.Type = BetLay4
	case LAY_10:
		expr.Type = BetLay10
	case BIG_6:
		expr.Type = BetBig6
	case BIG_8:
		expr.Type = BetBig8
	case HOP:
		expr.Type = BetHop
		expr.Args = p.parseHopCombination()
	case HOP_HARD_6:
		expr.Type = BetHopHard6
	case HOP_EASY_8:
		expr.Type = BetHopEasy8
	case WORLD:
		expr.Type = BetWorld
	case C_AND_E:
		expr.Type = BetCAndE
	case HORN:
		expr.Type = BetHorn
	case HORN_HIGH_11:
		expr.Type = BetHornHigh11
	case HORN_HIGH_ACE_DEUCE:
		expr.Type = BetHornHighAceDeuce
	// Missing bet type cases from canonical definitions
	// Buy bets
	case BUY_5:
		expr.Type = BetBuy5
	case BUY_6:
		expr.Type = BetBuy6
	case BUY_8:
		expr.Type = BetBuy8
	case BUY_9:
		expr.Type = BetBuy9
	// Lay bets
	case LAY_5:
		expr.Type = BetLay5
	case LAY_6:
		expr.Type = BetLay6
	case LAY_8:
		expr.Type = BetLay8
	case LAY_9:
		expr.Type = BetLay9
	// Place-to-lose bets
	case PLACE_TO_LOSE_4:
		expr.Type = BetPlaceToLose4
	case PLACE_TO_LOSE_5:
		expr.Type = BetPlaceToLose5
	case PLACE_TO_LOSE_6:
		expr.Type = BetPlaceToLose6
	case PLACE_TO_LOSE_8:
		expr.Type = BetPlaceToLose8
	case PLACE_TO_LOSE_9:
		expr.Type = BetPlaceToLose9
	case PLACE_TO_LOSE_10:
		expr.Type = BetPlaceToLose10
	// Horn high bets
	case HORN_HIGH_2:
		expr.Type = BetHornHigh2
	case HORN_HIGH_3:
		expr.Type = BetHornHigh3
	case HORN_HIGH_12:
		expr.Type = BetHornHigh12
	// Hop bets (all combinations)
	case HOP_1_2:
		expr.Type = BetHop12
	case HOP_1_3:
		expr.Type = BetHop13
	case HOP_1_4:
		expr.Type = BetHop14
	case HOP_1_5:
		expr.Type = BetHop15
	case HOP_1_6:
		expr.Type = BetHop16
	case HOP_2_3:
		expr.Type = BetHop23
	case HOP_2_4:
		expr.Type = BetHop24
	case HOP_2_5:
		expr.Type = BetHop25
	case HOP_2_6:
		expr.Type = BetHop26
	case HOP_3_4:
		expr.Type = BetHop34
	case HOP_3_5:
		expr.Type = BetHop35
	case HOP_3_6:
		expr.Type = BetHop36
	case HOP_4_5:
		expr.Type = BetHop45
	case HOP_4_6:
		expr.Type = BetHop46
	case HOP_5_6:
		expr.Type = BetHop56
	// Odds bets (specific types)
	case COME_ODDS:
		expr.Type = BetComeOdds
	case DONT_COME_ODDS:
		expr.Type = BetDontComeOdds
	default:
		p.addError(fmt.Sprintf("unknown bet type: %s", p.curToken.Literal))
		return nil
	}

	return expr
}

func (p *Parser) parsePlaceNumbers() []Expression {
	var numbers []Expression

	if !p.expectPeek(LPAREN) {
		return numbers
	}
	p.nextToken() // consume (

	// Valid place numbers in craps
	validPlaceNumbers := map[int]bool{4: true, 5: true, 6: true, 8: true, 9: true, 10: true}

	for !p.peekTokenIs(RPAREN) && !p.peekTokenIs(EOF) {
		p.nextToken()
		if p.curToken.Type != NUMBER {
			p.addError(fmt.Sprintf("expected number, got %s", p.curToken.Literal))
			return numbers
		}

		// Parse the number
		val, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError(fmt.Sprintf("invalid number: %s", p.curToken.Literal))
			return numbers
		}

		// Validate that it's a valid place number
		if !validPlaceNumbers[val] {
			p.addError(fmt.Sprintf("invalid place number: %s", p.curToken.Literal))
			return numbers
		}

		// Create NumberExpression for the valid place number
		expr := &NumberExpression{Token: p.curToken, Value: float64(val)}
		numbers = append(numbers, expr)

		if p.peekTokenIs(COMMA) {
			p.nextToken() // consume comma
		}
	}

	if !p.expectPeek(RPAREN) {
		return numbers
	}

	return numbers
}

func (p *Parser) parseHopCombination() []Expression {
	var combinations []Expression

	if !p.expectPeek(LPAREN) {
		return combinations
	}
	p.nextToken() // consume (

	// Valid hop combinations in craps (die1-die2 format)
	validHopCombinations := map[string]bool{
		"1-2": true, "1-3": true, "1-4": true, "1-5": true, "1-6": true,
		"2-3": true, "2-4": true, "2-5": true, "2-6": true,
		"3-4": true, "3-5": true, "3-6": true,
		"4-5": true, "4-6": true,
		"5-6": true,
	}

	for !p.peekTokenIs(RPAREN) && !p.peekTokenIs(EOF) {
		p.nextToken()
		if p.curToken.Type != NUMBER {
			p.addError(fmt.Sprintf("expected number, got %s", p.curToken.Literal))
			return combinations
		}

		// Parse first die value
		val1, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError(fmt.Sprintf("invalid number: %s", p.curToken.Literal))
			return combinations
		}

		// Validate first die is 1-6
		if val1 < 1 || val1 > 6 {
			p.addError(fmt.Sprintf("invalid die value: %s (must be 1-6)", p.curToken.Literal))
			return combinations
		}

		// Create NumberExpression for first die
		expr1 := &NumberExpression{Token: p.curToken, Value: float64(val1)}
		combinations = append(combinations, expr1)

		if !p.expectPeek(COMMA) {
			return combinations
		}
		p.nextToken() // consume comma

		if !p.expectPeek(NUMBER) {
			return combinations
		}
		p.nextToken() // consume second number

		// Parse second die value
		val2, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError(fmt.Sprintf("invalid number: %s", p.curToken.Literal))
			return combinations
		}

		// Validate second die is 1-6
		if val2 < 1 || val2 > 6 {
			p.addError(fmt.Sprintf("invalid die value: %s (must be 1-6)", p.curToken.Literal))
			return combinations
		}

		// Validate the combination is valid
		combination := fmt.Sprintf("%d-%d", val1, val2)
		if !validHopCombinations[combination] {
			p.addError(fmt.Sprintf("invalid hop combination: %d-%d", val1, val2))
			return combinations
		}

		// Create NumberExpression for second die
		expr2 := &NumberExpression{Token: p.curToken, Value: float64(val2)}
		combinations = append(combinations, expr2)

		if p.peekTokenIs(COMMA) {
			p.nextToken() // consume comma
		}
	}

	if !p.expectPeek(RPAREN) {
		return combinations
	}

	return combinations
}

func (p *Parser) parseModifiers() []*ModifierExpression {
	var modifiers []*ModifierExpression

	// Track modifier types to validate combinations
	usedModifiers := make(map[ModifierType]bool)

	for p.curToken.Type != SEMICOLON && p.curToken.Type != EOF {
		mod := &ModifierExpression{Token: p.curToken}

		switch p.curToken.Type {
		case WORKING_KEYWORD:
			mod.Type = ModWorking
			// Check for conflicting OFF modifier
			if usedModifiers[ModOff] {
				p.addError("cannot have both WORKING and OFF modifiers")
				return modifiers
			}
		case OFF_MODIFIER:
			mod.Type = ModOff
			// Check for conflicting WORKING modifier
			if usedModifiers[ModWorking] {
				p.addError("cannot have both WORKING and OFF modifiers")
				return modifiers
			}
		case PRESS:
			mod.Type = ModPress
			// PRESS can have a value (amount)
			if p.peekTokenIs(DOLLAR) || p.peekTokenIs(NUMBER) {
				p.nextToken() // consume $ or number
				mod.Value = p.parsePrimaryExpression()
			}
		case ODDS:
			mod.Type = ModRatio
			// ODDS can have a value (e.g., "3X")
			if p.peekTokenIs(NUMBER) {
				p.nextToken() // consume number
				firstNum := p.curToken.Literal
				if p.peekTokenIs(IDENT) {
					p.nextToken() // consume X or other identifier
					ratioStr := firstNum + ":" + p.curToken.Literal
					mod.Value = &IdentifierExpression{Token: p.curToken, Value: ratioStr}
				} else {
					// Just a number, treat as 1:number
					ratioStr := "1:" + firstNum
					mod.Value = &IdentifierExpression{Token: p.curToken, Value: ratioStr}
				}
			} else {
				p.addError("ODDS modifier requires a value")
				return modifiers
			}
		case ONE_ROLL:
			mod.Type = ModOneRoll
		case MAX:
			mod.Type = ModMax
		case AMOUNT:
			mod.Type = ModAmount
			// AMOUNT must have a value
			if p.peekTokenIs(DOLLAR) || p.peekTokenIs(NUMBER) {
				p.nextToken() // consume $ or number
				mod.Value = p.parsePrimaryExpression()
			} else {
				p.addError("AMOUNT modifier requires a value")
				return modifiers
			}
		case RATIO:
			mod.Type = ModRatio
			// RATIO must have a value (e.g., "2:1")
			if p.peekTokenIs(NUMBER) {
				p.nextToken() // consume first number
				firstNum := p.curToken.Literal
				if p.peekTokenIs(COLON) {
					p.nextToken() // consume :
					if p.peekTokenIs(NUMBER) {
						p.nextToken() // consume second number
						secondNum := p.curToken.Literal
						ratioStr := firstNum + ":" + secondNum
						mod.Value = &IdentifierExpression{Token: p.curToken, Value: ratioStr}
					} else {
						p.addError("RATIO modifier requires format 'number:number'")
						return modifiers
					}
				} else {
					p.addError("RATIO modifier requires format 'number:number'")
					return modifiers
				}
			} else {
				p.addError("RATIO modifier requires a value")
				return modifiers
			}
		default:
			p.addError(fmt.Sprintf("invalid modifier: %s", p.curToken.Literal))
			return modifiers
		}

		// Track this modifier type
		usedModifiers[mod.Type] = true
		modifiers = append(modifiers, mod)
		p.nextToken()
	}

	return modifiers
}

func (p *Parser) parseConditionalStatement() *ConditionalStatement {
	stmt := &ConditionalStatement{Token: p.curToken}

	p.nextToken() // consume IF

	// Parse full conditional expression with support for complex comparisons
	stmt.Condition = p.parsePrimaryExpression()
	p.nextToken() // advance to next token

	// Check if we have a comparison operator
	if p.curTokenIs(GT) || p.curTokenIs(LT) || p.curTokenIs(EQ) || p.curTokenIs(NOT_EQ) {
		operator := p.curToken.Literal
		p.nextToken() // consume operator
		right := p.parsePrimaryExpression()
		p.nextToken() // advance to next token

		stmt.Condition = &InfixExpression{
			Token:    p.curToken,
			Left:     stmt.Condition,
			Operator: operator,
			Right:    right,
		}
	}

	if !p.curTokenIs(THEN) {
		p.addError(fmt.Sprintf("expected THEN, got %s", p.curToken.Literal))
		return nil
	}
	p.nextToken() // consume THEN

	stmt.Consequence = p.parseBlockStatement()

	// Check for ELSE clause
	if p.peekTokenIs(ELSE) {
		p.nextToken() // consume ELSE
		stmt.Alternative = p.parseBlockStatement()
	}

	// Expect END token
	if p.peekTokenIs(END) {
		p.nextToken() // consume END
		if p.peekTokenIs(SEMICOLON) {
			p.nextToken() // consume semicolon
		}
	}

	return stmt
}

// parsePrimaryExpression parses primary expressions (identifiers, numbers, amounts)
func (p *Parser) parsePrimaryExpression() Expression {
	switch p.curToken.Type {
	case IDENT:
		expr := &IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
		return expr
	case NUMBER:
		val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		expr := &NumberExpression{Token: p.curToken, Value: val}
		return expr
	case DOLLAR:
		token := p.curToken
		p.nextToken() // consume $
		if p.curToken.Type != NUMBER {
			p.addError(fmt.Sprintf("expected number after $, got %s", p.curToken.Literal))
			return &NumberExpression{Token: token, Value: 0}
		}
		val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		expr := &AmountExpression{Token: p.curToken, Value: val}
		return expr
	default:
		p.addError(fmt.Sprintf("unexpected token in expression: %s", p.curToken.Literal))
		return &NumberExpression{Token: p.curToken, Value: 0}
	}
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	// Check if we have a brace to start a block
	if p.peekTokenIs(LBRACE) {
		p.nextToken() // consume {

		for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
			stmt := p.parseStatement()
			if stmt != nil {
				block.Statements = append(block.Statements, stmt)
			}
			p.nextToken()
		}

		if p.curTokenIs(EOF) {
			p.addError("unexpected end of input: missing closing brace '}'")
			return nil
		}
	} else {
		// Single statement without braces (for IF/THEN)
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

func (p *Parser) parseQueryStatement() *QueryStatement {
	stmt := &QueryStatement{Token: p.curToken}

	p.nextToken() // consume SHOW

	switch p.curToken.Type {
	case IDENT:
		switch p.curToken.Literal {
		case "POINT":
			stmt.Type = QueryPoint
		case "BETS":
			stmt.Type = QueryBets
		case "BANKROLL":
			stmt.Type = QueryBankroll
		case "TABLE_MINIMUMS":
			stmt.Type = QueryTableMinimums
		case "ODDS_ALLOWED":
			stmt.Type = QueryOddsAllowed
		default:
			p.addError(fmt.Sprintf("unknown query type: %s", p.curToken.Literal))
			return nil
		}
	default:
		p.addError(fmt.Sprintf("expected identifier, got %s", p.curToken.Literal))
		return nil
	}

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseManagementStatement() *ManagementStatement {
	stmt := &ManagementStatement{Token: p.curToken}

	p.nextToken() // consume SET

	// Parse management type - can be either IDENT or specific tokens
	switch p.curToken.Type {
	case IDENT:
		switch p.curToken.Literal {
		case "BANKROLL":
			stmt.Type = ManageBankroll
		case "MIN_BET":
			stmt.Type = ManageMinBet
		case "WIN_GOAL":
			stmt.Type = ManageWinGoal
		case "LOSS_LIMIT":
			stmt.Type = ManageLossLimit
		case "SESSION_TIME":
			stmt.Type = ManageSessionTime
		default:
			p.addError(fmt.Sprintf("unknown management type: %s", p.curToken.Literal))
			return nil
		}
	case MAX_BET:
		stmt.Type = ManageMaxBet
	default:
		p.addError(fmt.Sprintf("expected identifier or management type, got %s", p.curToken.Literal))
		return nil
	}

	p.nextToken() // consume management type

	// Handle optional "TO" keyword
	if p.curToken.Type == TO {
		p.nextToken() // consume TO
	}

	// Parse value with support for different types
	switch p.curToken.Type {
	case DOLLAR:
		p.nextToken() // consume $
		if p.curToken.Type != NUMBER {
			p.addError(fmt.Sprintf("expected number, got %s", p.curToken.Literal))
			return nil
		}
		stmt.Value = &AmountExpression{Token: p.curToken}
		val, err := parseAmount(p.curToken.Literal)
		if err != nil {
			p.addError(fmt.Sprintf("invalid amount: %s", p.curToken.Literal))
			return nil
		}
		stmt.Value.(*AmountExpression).Value = val
	case NUMBER:
		// Handle numeric values (like session time in minutes)
		stmt.Value = &NumberExpression{Token: p.curToken}
		val, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid number: %s", p.curToken.Literal))
			return nil
		}
		stmt.Value.(*NumberExpression).Value = val
	case IDENT:
		// Handle identifier values (like "ON", "OFF", etc.)
		stmt.Value = &IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	default:
		p.addError(fmt.Sprintf("expected $, number, or identifier, got %s", p.curToken.Literal))
		return nil
	}

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseRemoveStatement() *RemoveStatement {
	stmt := &RemoveStatement{Token: p.curToken}

	p.nextToken() // consume REMOVE

	// Check if the current token is ALL
	if p.curTokenIs(ALL) {
		// REMOVE ALL case - BetType remains nil
		// Don't advance past ALL, let expectPeek handle the semicolon
	} else {
		// REMOVE <bet_type> case - parse the bet type
		stmt.BetType = p.parseBetTypeExpression()
		// Do NOT advance past the bet type here
	}

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parsePressStatement() *PressStatement {
	stmt := &PressStatement{Token: p.curToken}

	p.nextToken() // consume PRESS

	// Parse bet type
	stmt.BetType = p.parseBetTypeExpression()

	if !p.expectPeek(BY) {
		return nil
	}
	p.nextToken() // consume BY

	if !p.expectPeek(DOLLAR) {
		return nil
	}
	p.nextToken() // consume $

	if !p.expectPeek(NUMBER) {
		return nil
	}
	p.nextToken() // consume number

	// Parse amount
	amount := &AmountExpression{Token: p.curToken}
	val, err := parseAmount(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid amount: %s", p.curToken.Literal))
		return nil
	}
	amount.Value = val
	stmt.Amount = amount

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseTurnStatement() *TurnStatement {
	stmt := &TurnStatement{Token: p.curToken}

	p.nextToken() // consume TURN

	switch p.curToken.Type {
	case ON_KEYWORD:
		stmt.Action = "ON"
	case OFF:
		stmt.Action = "OFF"
	default:
		p.addError(fmt.Sprintf("expected ON or OFF, got %s", p.curToken.Literal))
		return nil
	}

	p.nextToken() // consume ON/OFF

	// Parse bet type
	stmt.BetType = p.parseBetTypeExpression()

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	return stmt
}

// Helper methods
func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t.String(), p.peekToken.Type.String())
	p.errors = append(p.errors, msg)
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseRollStatement() *RollStatement {
	stmt := &RollStatement{Token: p.curToken}

	if !p.expectPeek(DICE) {
		return nil
	}
	if !p.expectPeek(SEMICOLON) {
		return nil
	}
	return stmt
}
