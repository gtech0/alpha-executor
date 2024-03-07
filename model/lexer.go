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
	quant  quant
	result []Token
}

type quant struct {
	token      Token
	difference int
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
		quant:  quant{Token{}, 0},
		result: make([]Token, 0),
	}
}

func (l *Lexer) Lex() []Token {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.result
			}

			panic(err)
		}

		if l.quant.token.Value != "" && l.quant.difference > 0 {
			l.result = append(l.result, l.quant.token)
			l.quant = quant{Token{}, 0}
		} else if l.quant.token.Value != "" {
			l.quant.difference++
		}

		l.pos.Column++
		switch r {
		case '\n':
			l.resetPosition()
		case '=':
			l.result = append(l.result, Token{EQUALS, "="})
			break
		case '<':
			l.result = append(l.result, Token{LESS_THAN, "<"})
			break
		case '≤':
			l.result = append(l.result, Token{LESS_THAN_EQUALS, "≤"})
			break
		case '>':
			l.result = append(l.result, Token{GREATER_THAN, ">"})
			break
		case '≥':
			l.result = append(l.result, Token{GREATER_THAN_EQUALS, "≥"})
			break
		case '∃':
			l.quant = quant{Token{EXIST, "∃"}, 0}
			continue
		case '∀':
			l.quant = quant{Token{FOR_ALL, "∀"}, 0}
			continue
		case '¬':
			l.result = append(l.result, Token{NOT, "¬"})
			break
		case '∨':
			l.result = append(l.result, Token{OR, "∨"})
			break
		case '∧':
			l.result = append(l.result, Token{AND, "∧"})
			break
		case '(':
			l.result = append(l.result, Token{LEFT_PARENTHESIS, "("})
			break
		case ')':
			l.result = append(l.result, Token{RIGHT_PARENTHESIS, ")"})
			break
		case ':':
			l.result = append(l.result, Token{LOGIC_START, ":"})
			break
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				l.backup()
				lit := l.lexInt()
				l.result = append(l.result, Token{INT, lit})
				break
			} else if unicode.IsLetter(r) {
				l.backup()
				lit, period := l.lexStr()
				if period {
					l.result = append(l.result, Token{ATTRIBUTE, lit})
					break
				}

				operations := []string{"GET", "RANGE", "HOLD", "RELEASE", "UPDATE"}
				if slices.Contains(operations, lit) {
					l.result = append(l.result, Token{OPERATION, lit})
					break
				}
				l.result = append(l.result, Token{RELATION, lit})
				break
			} else if r == '\'' {
				l.backup()
				lit, period := l.lexStr()
				if period {
					l.result = append(l.result, Token{ILLEGAL, lit})
					break
				}
				l.result = append(l.result, Token{CONST, lit})
				break
			} else {
				l.result = append(l.result, Token{ILLEGAL, string(r)})
				break
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
			l.backup()
			return lit
		}
	}
}

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
			l.backup()
			return lit, period
		}
	}

}
