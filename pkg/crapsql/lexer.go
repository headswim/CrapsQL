package crapsql

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar() // initialize first character
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.column += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, l.column)
	case ':':
		tok = newToken(COLON, l.ch, l.line, l.column)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, l.column)
	case ',':
		tok = newToken(COMMA, l.ch, l.line, l.column)
	case '+':
		tok = newToken(PLUS, l.ch, l.line, l.column)
	case '-':
		tok = newToken(MINUS, l.ch, l.line, l.column)
	case '*':
		tok = newToken(ASTERISK, l.ch, l.line, l.column)
	case '/':
		tok = newToken(SLASH, l.ch, l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(BANG, l.ch, l.line, l.column)
		}
	case '<':
		tok = newToken(LT, l.ch, l.line, l.column)
	case '>':
		tok = newToken(GT, l.ch, l.line, l.column)
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(EQUALS, l.ch, l.line, l.column)
		}
	case '$':
		tok = newToken(DOLLAR, l.ch, l.line, l.column)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, l.column)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = NUMBER
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			return tok
		} else {
			// Handle illegal characters
			tok = Token{
				Type:    ILLEGAL,
				Literal: string(l.ch),
				Line:    l.line,
				Column:  l.column,
			}
			l.readChar()
			return tok
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// Handle decimal point
	if l.ch == '.' {
		l.readChar()
		// Read digits after decimal point
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func (l *Lexer) lookupIdent(ident string) TokenType {
	switch ident {
	case "PLACE":
		return PLACE
	case "ON":
		return ON
	case "WITH":
		return WITH
	case "IF":
		return IF
	case "THEN":
		return THEN
	case "ELSE":
		return ELSE
	case "END":
		return END
	case "SET":
		return SET
	case "SHOW":
		return SHOW
	case "DEFINE":
		return DEFINE
	case "AS":
		return AS
	case "EXECUTE":
		return EXECUTE
	case "APPLY":
		return APPLY
	case "TO":
		return TO
	case "REMOVE":
		return REMOVE
	case "ALL":
		return ALL
	case "TURN":
		return TURN
	case "BY":
		return BY
	case "START_BET":
		return START_BET
	case "ON_LOSS":
		return ON_LOSS
	case "ON_WIN":
		return ON_WIN
	case "MAX_BET":
		return MAX_BET
	case "MULTIPLY":
		return MULTIPLY
	case "RESET":
		return RESET
	case "TAKE_DOWN":
		return TAKE_DOWN
	case "WORKING":
		return WORKING_KEYWORD
	// Bet types
	case "PASS_LINE":
		return PASS_LINE
	case "DONT_PASS":
		return DONT_PASS
	case "COME":
		return COME
	case "DONT_COME":
		return DONT_COME
	case "FIELD":
		return FIELD
	case "ANY_SEVEN":
		return ANY_SEVEN
	case "ANY_CRAPS":
		return ANY_CRAPS
	case "ELEVEN":
		return ELEVEN
	case "ACE_DEUCE":
		return ACE_DEUCE
	case "ACES":
		return ACES
	case "BOXCARS":
		return BOXCARS
	case "PLACE_4":
		return PLACE_4
	case "PLACE_5":
		return PLACE_5
	case "PLACE_6":
		return PLACE_6
	case "PLACE_8":
		return PLACE_8
	case "PLACE_9":
		return PLACE_9
	case "PLACE_10":
		return PLACE_10
	case "PLACE_NUMBERS":
		return PLACE_NUMBERS
	case "PLACE_INSIDE":
		return PLACE_INSIDE
	case "PLACE_OUTSIDE":
		return PLACE_OUTSIDE
	case "HARD_4":
		return HARD_4
	case "HARD_6":
		return HARD_6
	case "HARD_8":
		return HARD_8
	case "HARD_10":
		return HARD_10
	case "ALL_HARDWAYS":
		return ALL_HARDWAYS
	case "PASS_ODDS":
		return PASS_ODDS
	case "DONT_PASS_ODDS":
		return DONT_PASS_ODDS
	case "BUY_4":
		return BUY_4
	case "BUY_10":
		return BUY_10
	case "LAY_4":
		return LAY_4
	case "LAY_10":
		return LAY_10
	case "BIG_6":
		return BIG_6
	case "BIG_8":
		return BIG_8
	case "HOP":
		return HOP
	case "HOP_HARD_6":
		return HOP_HARD_6
	case "HOP_EASY_8":
		return HOP_EASY_8
	case "WORLD":
		return WORLD
	case "C_AND_E":
		return C_AND_E
	case "HORN":
		return HORN
	case "HORN_HIGH_11":
		return HORN_HIGH_11
	case "HORN_HIGH_ACE_DEUCE":
		return HORN_HIGH_ACE_DEUCE
	// Missing bet type tokens from canonical definitions
	// Buy bets
	case "BUY_5":
		return BUY_5
	case "BUY_6":
		return BUY_6
	case "BUY_8":
		return BUY_8
	case "BUY_9":
		return BUY_9
	// Lay bets
	case "LAY_5":
		return LAY_5
	case "LAY_6":
		return LAY_6
	case "LAY_8":
		return LAY_8
	case "LAY_9":
		return LAY_9
	// Place-to-lose bets
	case "PLACE_TO_LOSE_4":
		return PLACE_TO_LOSE_4
	case "PLACE_TO_LOSE_5":
		return PLACE_TO_LOSE_5
	case "PLACE_TO_LOSE_6":
		return PLACE_TO_LOSE_6
	case "PLACE_TO_LOSE_8":
		return PLACE_TO_LOSE_8
	case "PLACE_TO_LOSE_9":
		return PLACE_TO_LOSE_9
	case "PLACE_TO_LOSE_10":
		return PLACE_TO_LOSE_10
	// Horn high bets
	case "HORN_HIGH_2":
		return HORN_HIGH_2
	case "HORN_HIGH_3":
		return HORN_HIGH_3
	case "HORN_HIGH_12":
		return HORN_HIGH_12
	// Hop bets (all combinations)
	case "HOP_1_2":
		return HOP_1_2
	case "HOP_1_3":
		return HOP_1_3
	case "HOP_1_4":
		return HOP_1_4
	case "HOP_1_5":
		return HOP_1_5
	case "HOP_1_6":
		return HOP_1_6
	case "HOP_2_3":
		return HOP_2_3
	case "HOP_2_4":
		return HOP_2_4
	case "HOP_2_5":
		return HOP_2_5
	case "HOP_2_6":
		return HOP_2_6
	case "HOP_3_4":
		return HOP_3_4
	case "HOP_3_5":
		return HOP_3_5
	case "HOP_3_6":
		return HOP_3_6
	case "HOP_4_5":
		return HOP_4_5
	case "HOP_4_6":
		return HOP_4_6
	case "HOP_5_6":
		return HOP_5_6
	// Odds bets (specific types)
	case "COME_ODDS":
		return COME_ODDS
	case "DONT_COME_ODDS":
		return DONT_COME_ODDS
	// Modifiers
	case "OFF":
		return OFF_MODIFIER
	case "PRESS":
		return PRESS
	case "ODDS":
		return ODDS
	case "ROLL":
		return ROLL
	case "DICE":
		return DICE
	case "ONE_ROLL":
		return ONE_ROLL
	case "MAX":
		return MAX
	case "AMOUNT":
		return AMOUNT
	case "RATIO":
		return RATIO
	default:
		return IDENT
	}
}
