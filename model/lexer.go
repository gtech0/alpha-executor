package model

import (
	"bufio"
	"io"
	"slices"
	"unicode"
)

type LexType int

const (
	EOF LexType = iota
	ILLEGAL
	OPERATION
	ATTRIBUTE
	CONST
	INT
	RELATION
	LOGIC_START

	EQUALS
	LESS_THAN
	LESS_THAN_EQUALS
	GREATER_THAN
	GREATER_THAN_EQUALS

	EXIST
	FOR_ALL

	NOT
	OR
	AND

	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
)

type Token struct {
	Type  LexType
	Value string
}

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, Token) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, Token{EOF, ""}
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		// update the column to the position of the newly read in rune
		l.pos.Column++

		switch r {
		case '\n':
			l.resetPosition()
		case '=':
			return l.pos, Token{EQUALS, "="}
		case '<':
			return l.pos, Token{LESS_THAN, "<"}
		case '≤':
			return l.pos, Token{LESS_THAN_EQUALS, "≤"}
		case '>':
			return l.pos, Token{GREATER_THAN, ">"}
		case '≥':
			return l.pos, Token{GREATER_THAN_EQUALS, "≥"}
		case '∃':
			return l.pos, Token{EXIST, "∃"}
		case '∀':
			return l.pos, Token{FOR_ALL, "∀"}
		case '¬':
			return l.pos, Token{NOT, "¬"}
		case '∨':
			return l.pos, Token{OR, "∨"}
		case '∧':
			return l.pos, Token{AND, "∧"}
		case '(':
			return l.pos, Token{LEFT_PARENTHESIS, "("}
		case ')':
			return l.pos, Token{RIGHT_PARENTHESIS, ")"}
		case ':':
			return l.pos, Token{LOGIC_START, ":"}
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, Token{INT, lit}
			} else if unicode.IsLetter(r) {
				// backup and let lexStr rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit, period := l.lexStr()
				if period {
					return startPos, Token{ATTRIBUTE, lit}
				}

				operations := []string{"GET", "RANGE", "HOLD", "RELEASE", "UPDATE"}
				if slices.Contains(operations, lit) {
					return startPos, Token{OPERATION, lit}
				}
				return startPos, Token{RELATION, lit}
			} else if r == '\'' {
				startPos := l.pos
				l.backup()
				lit, period := l.lexStr()
				if period {
					return startPos, Token{ILLEGAL, lit}
				}
				return startPos, Token{CONST, lit}
			} else {
				return l.pos, Token{ILLEGAL, string(r)}
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.Line++
	l.pos.Column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.Column--
}

// lexInt scans the input until the end of an integer
// and then returns the literal.
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.pos.Column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

// lexStr scans the input until the end of an identifier
// and then returns the literal.
func (l *Lexer) lexStr() (string, bool) {
	lit := ""
	quoteCount := 0
	special := []rune{'.', '-', '/'}
	period := false
	for {
		if quoteCount == 2 {
			return lit, period
		}

		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit, period
			}
		}

		l.pos.Column++
		if unicode.IsLetter(r) || unicode.IsDigit(r) || slices.Contains(special, r) {
			lit = lit + string(r)
			if r == '.' {
				period = true
			}
		} else if r == '\'' {
			quoteCount++
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit, period
		}
	}
}
