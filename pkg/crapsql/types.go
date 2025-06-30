package crapsql

import (
	"fmt"
	"strconv"
)

// Token represents a lexical token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	SEMICOLON

	// Literals
	IDENT  // bet types, keywords
	DOLLAR // $
	NUMBER // 25, 100, etc.

	// Keywords
	PLACE
	ON
	WITH
	IF
	THEN
	ELSE
	END
	SET
	SHOW
	DEFINE

	AS
	EXECUTE
	APPLY
	TO
	REMOVE
	ALL
	TURN
	OFF
	ON_KEYWORD
	BY
	START_BET
	ON_LOSS
	ON_WIN
	MAX_BET
	MULTIPLY
	RESET
	TAKE_DOWN
	WORKING_KEYWORD
	ROLL
	DICE

	// Bet types
	PASS_LINE
	DONT_PASS
	COME
	DONT_COME
	FIELD
	ANY_SEVEN
	ANY_CRAPS
	ELEVEN
	ACE_DEUCE
	ACES
	BOXCARS
	PLACE_4
	PLACE_5
	PLACE_6
	PLACE_8
	PLACE_9
	PLACE_10
	PLACE_NUMBERS
	PLACE_INSIDE
	PLACE_OUTSIDE
	HARD_4
	HARD_6
	HARD_8
	HARD_10
	ALL_HARDWAYS
	PASS_ODDS
	DONT_PASS_ODDS
	BUY_4
	BUY_10
	LAY_4
	LAY_10
	BIG_6
	BIG_8
	HOP
	HOP_HARD_6
	HOP_EASY_8
	WORLD
	C_AND_E
	HORN
	HORN_HIGH_11
	HORN_HIGH_ACE_DEUCE

	// Missing bet type tokens from canonical definitions
	// Buy bets
	BUY_5
	BUY_6
	BUY_8
	BUY_9

	// Lay bets
	LAY_5
	LAY_6
	LAY_8
	LAY_9

	// Place-to-lose bets
	PLACE_TO_LOSE_4
	PLACE_TO_LOSE_5
	PLACE_TO_LOSE_6
	PLACE_TO_LOSE_8
	PLACE_TO_LOSE_9
	PLACE_TO_LOSE_10

	// Horn high bets
	HORN_HIGH_2
	HORN_HIGH_3
	HORN_HIGH_12

	// Hop bets (all combinations)
	HOP_1_2
	HOP_1_3
	HOP_1_4
	HOP_1_5
	HOP_1_6
	HOP_2_3
	HOP_2_4
	HOP_2_5
	HOP_2_6
	HOP_3_4
	HOP_3_5
	HOP_3_6
	HOP_4_5
	HOP_4_6
	HOP_5_6

	// Odds bets (specific types)
	COME_ODDS
	DONT_COME_ODDS

	// Modifiers
	WORKING
	OFF_MODIFIER
	PRESS
	ODDS
	ONE_ROLL
	MAX
	AMOUNT
	RATIO

	// Operators
	EQUALS
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	COMMA
	COLON
	PLUS
	MINUS
	ASTERISK
	SLASH
	BANG
	LT
	GT
	EQ
	NOT_EQ
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// AST Node interface
type Node interface {
	TokenLiteral() string
}

// Statement interface
type Statement interface {
	Node
	statementNode()
}

// Expression interface
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// BetStatement represents a PLACE bet command
type BetStatement struct {
	Token     Token
	Amount    *AmountExpression
	BetType   *BetTypeExpression
	Modifiers []*ModifierExpression
}

func (bs *BetStatement) statementNode()       {}
func (bs *BetStatement) TokenLiteral() string { return bs.Token.Literal }

// AmountExpression represents a dollar amount
type AmountExpression struct {
	Token Token
	Value float64
}

func (ae *AmountExpression) expressionNode()      {}
func (ae *AmountExpression) TokenLiteral() string { return ae.Token.Literal }

// BetTypeExpression represents a bet type
type BetTypeExpression struct {
	Token Token
	Type  BetType
	Args  []Expression // For place numbers, hop combinations, etc.
}

func (bte *BetTypeExpression) expressionNode()      {}
func (bte *BetTypeExpression) TokenLiteral() string { return bte.Token.Literal }

// ModifierExpression represents bet modifiers
type ModifierExpression struct {
	Token Token
	Type  ModifierType
	Value Expression // For ratios, amounts, etc.
}

func (me *ModifierExpression) expressionNode()      {}
func (me *ModifierExpression) TokenLiteral() string { return me.Token.Literal }

// IdentifierExpression represents an identifier
type IdentifierExpression struct {
	Token Token
	Value string
}

func (ie *IdentifierExpression) expressionNode()      {}
func (ie *IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }

// NumberExpression represents a number literal
type NumberExpression struct {
	Token Token
	Value float64
}

func (ne *NumberExpression) expressionNode()      {}
func (ne *NumberExpression) TokenLiteral() string { return ne.Token.Literal }

// InfixExpression represents an infix operation
type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// ConditionalStatement represents IF/THEN/ELSE blocks
type ConditionalStatement struct {
	Token       Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (cs *ConditionalStatement) statementNode()       {}
func (cs *ConditionalStatement) TokenLiteral() string { return cs.Token.Literal }

// BlockStatement represents a block of statements
type BlockStatement struct {
	Token      Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// QueryStatement represents SHOW commands
type QueryStatement struct {
	Token Token
	Type  QueryType
}

func (qs *QueryStatement) statementNode()       {}
func (qs *QueryStatement) TokenLiteral() string { return qs.Token.Literal }

// ManagementStatement represents SET commands
type ManagementStatement struct {
	Token Token
	Type  ManagementType
	Value Expression
}

func (ms *ManagementStatement) statementNode()       {}
func (ms *ManagementStatement) TokenLiteral() string { return ms.Token.Literal }

// RemoveStatement represents REMOVE BET commands
type RemoveStatement struct {
	Token   Token
	BetType *BetTypeExpression
}

func (rs *RemoveStatement) statementNode()       {}
func (rs *RemoveStatement) TokenLiteral() string { return rs.Token.Literal }

// PressStatement represents PRESS commands
type PressStatement struct {
	Token   Token
	BetType *BetTypeExpression
	Amount  *AmountExpression
}

func (ps *PressStatement) statementNode()       {}
func (ps *PressStatement) TokenLiteral() string { return ps.Token.Literal }

// TurnStatement represents TURN ON/OFF commands
type TurnStatement struct {
	Token   Token
	Action  string // "ON" or "OFF"
	BetType *BetTypeExpression
}

func (ts *TurnStatement) statementNode()       {}
func (ts *TurnStatement) TokenLiteral() string { return ts.Token.Literal }

// RollStatement represents a ROLL DICE command
type RollStatement struct {
	Token Token
}

func (rs *RollStatement) statementNode()       {}
func (rs *RollStatement) TokenLiteral() string { return rs.Token.Literal }

// Bet types
type BetType int

const (
	// Line bets
	BetPassLine BetType = iota
	BetDontPass
	BetCome
	BetDontCome

	// Field and one-roll bets
	BetField
	BetAnySeven
	BetAnyCraps
	BetEleven
	BetAceDeuce
	BetAces
	BetBoxcars

	// Place bets
	BetPlace4
	BetPlace5
	BetPlace6
	BetPlace8
	BetPlace9
	BetPlace10
	BetPlaceNumbers
	BetPlaceInside
	BetPlaceOutside

	// Hard ways
	BetHard4
	BetHard6
	BetHard8
	BetHard10
	BetAllHardways

	// Odds
	BetPassOdds
	BetDontPassOdds

	// Buy/Lay
	BetBuy4
	BetBuy10
	BetLay4
	BetLay10

	// Big 6/8
	BetBig6
	BetBig8

	// Hop bets
	BetHop
	BetHopHard6
	BetHopEasy8

	// Proposition bets
	BetWorld
	BetCAndE
	BetHorn
	BetHornHigh11
	BetHornHighAceDeuce

	// Missing bet types from canonical definitions
	// Buy bets
	BetBuy5
	BetBuy6
	BetBuy8
	BetBuy9

	// Lay bets
	BetLay5
	BetLay6
	BetLay8
	BetLay9

	// Place-to-lose bets
	BetPlaceToLose4
	BetPlaceToLose5
	BetPlaceToLose6
	BetPlaceToLose8
	BetPlaceToLose9
	BetPlaceToLose10

	// Horn high bets
	BetHornHigh2
	BetHornHigh3
	BetHornHigh12

	// Hop bets (all combinations)
	BetHop12
	BetHop13
	BetHop14
	BetHop15
	BetHop16
	BetHop23
	BetHop24
	BetHop25
	BetHop26
	BetHop34
	BetHop35
	BetHop36
	BetHop45
	BetHop46
	BetHop56

	// Odds bets (specific types)
	BetComeOdds
	BetDontComeOdds
)

// Modifier types
type ModifierType int

const (
	ModWorking ModifierType = iota
	ModOff
	ModPress
	ModOneRoll
	ModMax
	ModAmount
	ModRatio
)

// Query types
type QueryType int

const (
	QueryPoint QueryType = iota
	QueryBets
	QueryBankroll
	QueryTableMinimums
	QueryOddsAllowed
)

// Management types
type ManagementType int

const (
	ManageBankroll ManagementType = iota
	ManageMaxBet
	ManageMinBet
	ManageWinGoal
	ManageLossLimit
	ManageSessionTime
)

// Error types
type ParseError struct {
	Message string
	Line    int
	Column  int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// Helper function to convert string to float64
func parseAmount(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// String returns the string representation of a TokenType
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case SEMICOLON:
		return "SEMICOLON"
	case IDENT:
		return "IDENT"
	case DOLLAR:
		return "DOLLAR"
	case NUMBER:
		return "NUMBER"
	case PLACE:
		return "PLACE"
	case ON:
		return "ON"
	case WITH:
		return "WITH"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ELSE:
		return "ELSE"
	case END:
		return "END"
	case SET:
		return "SET"
	case SHOW:
		return "SHOW"
	case DEFINE:
		return "DEFINE"
	case AS:
		return "AS"
	case EXECUTE:
		return "EXECUTE"
	case APPLY:
		return "APPLY"
	case TO:
		return "TO"
	case REMOVE:
		return "REMOVE"
	case ALL:
		return "ALL"
	case TURN:
		return "TURN"
	case OFF:
		return "OFF"
	case ON_KEYWORD:
		return "ON_KEYWORD"
	case BY:
		return "BY"
	case START_BET:
		return "START_BET"
	case ON_LOSS:
		return "ON_LOSS"
	case ON_WIN:
		return "ON_WIN"
	case MAX_BET:
		return "MAX_BET"
	case MULTIPLY:
		return "MULTIPLY"
	case RESET:
		return "RESET"
	case TAKE_DOWN:
		return "TAKE_DOWN"
	case WORKING_KEYWORD:
		return "WORKING_KEYWORD"
	case ROLL:
		return "ROLL"
	case DICE:
		return "DICE"
	case PASS_LINE:
		return "PASS_LINE"
	case DONT_PASS:
		return "DONT_PASS"
	case COME:
		return "COME"
	case DONT_COME:
		return "DONT_COME"
	case FIELD:
		return "FIELD"
	case ANY_SEVEN:
		return "ANY_SEVEN"
	case ANY_CRAPS:
		return "ANY_CRAPS"
	case ELEVEN:
		return "ELEVEN"
	case ACE_DEUCE:
		return "ACE_DEUCE"
	case ACES:
		return "ACES"
	case BOXCARS:
		return "BOXCARS"
	case PLACE_4:
		return "PLACE_4"
	case PLACE_5:
		return "PLACE_5"
	case PLACE_6:
		return "PLACE_6"
	case PLACE_8:
		return "PLACE_8"
	case PLACE_9:
		return "PLACE_9"
	case PLACE_10:
		return "PLACE_10"
	case PLACE_NUMBERS:
		return "PLACE_NUMBERS"
	case PLACE_INSIDE:
		return "PLACE_INSIDE"
	case PLACE_OUTSIDE:
		return "PLACE_OUTSIDE"
	case HARD_4:
		return "HARD_4"
	case HARD_6:
		return "HARD_6"
	case HARD_8:
		return "HARD_8"
	case HARD_10:
		return "HARD_10"
	case ALL_HARDWAYS:
		return "ALL_HARDWAYS"
	case PASS_ODDS:
		return "PASS_ODDS"
	case DONT_PASS_ODDS:
		return "DONT_PASS_ODDS"
	case BUY_4:
		return "BUY_4"
	case BUY_10:
		return "BUY_10"
	case LAY_4:
		return "LAY_4"
	case LAY_10:
		return "LAY_10"
	case BIG_6:
		return "BIG_6"
	case BIG_8:
		return "BIG_8"
	case HOP:
		return "HOP"
	case HOP_HARD_6:
		return "HOP_HARD_6"
	case HOP_EASY_8:
		return "HOP_EASY_8"
	case WORLD:
		return "WORLD"
	case C_AND_E:
		return "C_AND_E"
	case HORN:
		return "HORN"
	case HORN_HIGH_11:
		return "HORN_HIGH_11"
	case HORN_HIGH_ACE_DEUCE:
		return "HORN_HIGH_ACE_DEUCE"
	case WORKING:
		return "WORKING"
	case OFF_MODIFIER:
		return "OFF_MODIFIER"
	case PRESS:
		return "PRESS"
	case ODDS:
		return "ODDS"
	case ONE_ROLL:
		return "ONE_ROLL"
	case MAX:
		return "MAX"
	case AMOUNT:
		return "AMOUNT"
	case RATIO:
		return "RATIO"
	case EQUALS:
		return "EQUALS"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case BANG:
		return "BANG"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case BUY_5:
		return "BUY_5"
	case BUY_6:
		return "BUY_6"
	case BUY_8:
		return "BUY_8"
	case BUY_9:
		return "BUY_9"
	case LAY_5:
		return "LAY_5"
	case LAY_6:
		return "LAY_6"
	case LAY_8:
		return "LAY_8"
	case LAY_9:
		return "LAY_9"
	case PLACE_TO_LOSE_4:
		return "PLACE_TO_LOSE_4"
	case PLACE_TO_LOSE_5:
		return "PLACE_TO_LOSE_5"
	case PLACE_TO_LOSE_6:
		return "PLACE_TO_LOSE_6"
	case PLACE_TO_LOSE_8:
		return "PLACE_TO_LOSE_8"
	case PLACE_TO_LOSE_9:
		return "PLACE_TO_LOSE_9"
	case PLACE_TO_LOSE_10:
		return "PLACE_TO_LOSE_10"
	case HORN_HIGH_2:
		return "HORN_HIGH_2"
	case HORN_HIGH_3:
		return "HORN_HIGH_3"
	case HORN_HIGH_12:
		return "HORN_HIGH_12"
	case HOP_1_2:
		return "HOP_1_2"
	case HOP_1_3:
		return "HOP_1_3"
	case HOP_1_4:
		return "HOP_1_4"
	case HOP_1_5:
		return "HOP_1_5"
	case HOP_1_6:
		return "HOP_1_6"
	case HOP_2_3:
		return "HOP_2_3"
	case HOP_2_4:
		return "HOP_2_4"
	case HOP_2_5:
		return "HOP_2_5"
	case HOP_2_6:
		return "HOP_2_6"
	case HOP_3_4:
		return "HOP_3_4"
	case HOP_3_5:
		return "HOP_3_5"
	case HOP_3_6:
		return "HOP_3_6"
	case HOP_4_5:
		return "HOP_4_5"
	case HOP_4_6:
		return "HOP_4_6"
	case HOP_5_6:
		return "HOP_5_6"
	case COME_ODDS:
		return "COME_ODDS"
	case DONT_COME_ODDS:
		return "DONT_COME_ODDS"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}
