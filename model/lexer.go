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
	ATTRIBUTE
	CONST
	INT
	RELATION
	LOGIC_START

	EQUALS
	NOT_EQUALS
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
	Type     LexType
	Value    string
	Position Position
}

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
	result []Token
	write  bool
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(reader),
		result: make([]Token, 0),
		write:  false,
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

		l.pos.Column++

		if r == ':' {
			l.write = true
			continue
		}

		if !l.write {
			continue
		}

		switch r {
		case '\n':
			l.resetPosition()
		case '=':
			l.result = append(l.result, Token{EQUALS, "=", l.pos})
			break
		case '≠':
			l.result = append(l.result, Token{NOT_EQUALS, "≠", l.pos})
			break
		case '<':
			l.result = append(l.result, Token{LESS_THAN, "<", l.pos})
			break
		case '≤':
			l.result = append(l.result, Token{LESS_THAN_EQUALS, "≤", l.pos})
			break
		case '>':
			l.result = append(l.result, Token{GREATER_THAN, ">", l.pos})
			break
		case '≥':
			l.result = append(l.result, Token{GREATER_THAN_EQUALS, "≥", l.pos})
			break
		case '∃':
			l.result = append(l.result, Token{EXIST, "∃", l.pos})
			break
		case '∀':
			l.result = append(l.result, Token{FOR_ALL, "∀", l.pos})
			break
		case '¬':
			l.result = append(l.result, Token{NOT, "¬", l.pos})
			break
		case '∨':
			l.result = append(l.result, Token{OR, "∨", l.pos})
			break
		case '∧':
			l.result = append(l.result, Token{AND, "∧", l.pos})
			break
		case '(':
			l.result = append(l.result, Token{LEFT_PARENTHESIS, "(", l.pos})
			break
		case ')':
			l.result = append(l.result, Token{RIGHT_PARENTHESIS, ")", l.pos})
			break
		case ':':
			l.result = append(l.result, Token{LOGIC_START, ":", l.pos})
			break
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				l.backup()
				lit := l.lexInt()
				l.result = append(l.result, Token{INT, lit, l.pos})
				break
			} else if unicode.IsLetter(r) {
				l.backup()
				lit, period := l.lexStr()
				if period {
					l.result = append(l.result, Token{ATTRIBUTE, lit, l.pos})
					break
				}

				l.result = append(l.result, Token{RELATION, lit, l.pos})
				break
			} else if r == '\'' {
				l.backup()
				lit, period := l.lexStr()
				if period {
					l.result = append(l.result, Token{ILLEGAL, lit, l.pos})
					break
				}
				l.result = append(l.result, Token{CONST, lit, l.pos})
				break
			} else {
				l.result = append(l.result, Token{ILLEGAL, string(r), l.pos})
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
