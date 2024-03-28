package model

import (
	"alpha-executor/entity"
	"bufio"
	"io"
	"slices"
	"unicode"
)

type LexType int

func (t LexType) String() string {
	return tokens[t]
}

const (
	EOF LexType = iota

	PROGRAM

	ILLEGAL
	ATTRIBUTE
	CONSTANT
	INTEGER
	NULL
	RELATION

	GET
	RANGE
	HOLD
	RELEASE
	UPDATE
	DELETE
	PUT

	EQUALS
	NOT_EQUALS
	LESS_THAN
	LESS_THAN_EQUALS
	GREATER_THAN
	GREATER_THAN_EQUALS

	EXIST
	FOR_ALL

	NEGATION
	CONJUNCTION
	DISJUNCTION
	IMPLICATION

	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
	COMMA
	LOGIC_START
)

var tokens = []string{
	EOF: "EOF",

	PROGRAM: "PROGRAM",

	ILLEGAL:   "ILLEGAL",
	ATTRIBUTE: "ATTRIBUTE",
	CONSTANT:  "CONSTANT",
	INTEGER:   "INTEGER",
	NULL:      "NULL",
	RELATION:  "RELATION",

	GET:     "GET",
	RANGE:   "RANGE",
	HOLD:    "HOLD",
	RELEASE: "RELEASE",
	UPDATE:  "UPDATE",
	DELETE:  "DELETE",
	PUT:     "PUT",

	EQUALS:              "=",
	NOT_EQUALS:          "≠",
	LESS_THAN:           "<",
	LESS_THAN_EQUALS:    "≤",
	GREATER_THAN:        ">",
	GREATER_THAN_EQUALS: "≥",

	EXIST:   "∃",
	FOR_ALL: "∀",

	NEGATION:    "¬",
	CONJUNCTION: "∧",
	DISJUNCTION: "∨",
	IMPLICATION: "→",

	LEFT_PARENTHESIS:  "(",
	RIGHT_PARENTHESIS: ")",
	COMMA:             ",",
	LOGIC_START:       ":",
}

type Token struct {
	Type     LexType
	Value    string
	Position entity.Position
}

type Lexer struct {
	pos     entity.Position
	reader  *bufio.Reader
	results [][]Token
}

func NewLexer(reader *bufio.Reader) *Lexer {
	return &Lexer{
		pos:     entity.Position{Line: 1, Column: 0},
		reader:  reader,
		results: make([][]Token, 0),
	}
}

func (l *Lexer) Lex() [][]Token {
	result := make([]Token, 0)
	for {
		r, _, err := l.reader.ReadRune()
		//fmt.Print(string(r) + "\n")
		if err != nil {
			if err == io.EOF {
				l.results = append(l.results, result)
				return l.results
			}

			panic(err)
		}

		l.pos.Column++

		if r == ';' {
			l.results = append(l.results, result)
			result = make([]Token, 0)
			continue
		}

		switch r {
		case '\n':
			l.nextLine()
		case '=':
			result = append(result, Token{EQUALS, EQUALS.String(), l.pos})
			break
		case '≠':
			result = append(result, Token{NOT_EQUALS, NOT_EQUALS.String(), l.pos})
			break
		case '<':
			result = append(result, Token{LESS_THAN, LESS_THAN.String(), l.pos})
			break
		case '≤':
			result = append(result, Token{LESS_THAN_EQUALS, LESS_THAN_EQUALS.String(), l.pos})
			break
		case '>':
			result = append(result, Token{GREATER_THAN, GREATER_THAN.String(), l.pos})
			break
		case '≥':
			result = append(result, Token{GREATER_THAN_EQUALS, GREATER_THAN_EQUALS.String(), l.pos})
			break
		case '∃':
			result = append(result, Token{EXIST, EXIST.String(), l.pos})
			break
		case '∀':
			result = append(result, Token{FOR_ALL, FOR_ALL.String(), l.pos})
			break
		case '¬':
			result = append(result, Token{NEGATION, NEGATION.String(), l.pos})
			break
		case '∧':
			result = append(result, Token{CONJUNCTION, CONJUNCTION.String(), l.pos})
			break
		case '∨':
			result = append(result, Token{DISJUNCTION, DISJUNCTION.String(), l.pos})
			break
		case '→':
			result = append(result, Token{IMPLICATION, IMPLICATION.String(), l.pos})
			break
		case '(':
			result = append(result, Token{LEFT_PARENTHESIS, LEFT_PARENTHESIS.String(), l.pos})
			break
		case ')':
			result = append(result, Token{RIGHT_PARENTHESIS, RIGHT_PARENTHESIS.String(), l.pos})
			break
		case ',':
			result = append(result, Token{COMMA, COMMA.String(), l.pos})
			break
		case ':':
			result = append(result, Token{LOGIC_START, LOGIC_START.String(), l.pos})
			break
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				l.backup()
				lit := l.lexInt()
				result = append(result, Token{INTEGER, lit, l.pos})
				break
			} else if unicode.IsLetter(r) {
				l.backup()
				lit, period := l.lexStr()
				if period {
					result = append(result, Token{ATTRIBUTE, lit, l.pos})
					break
				}

				switch lit {
				case "GET":
					result = append(result, Token{GET, lit, l.pos})
					break
				case "RANGE":
					result = append(result, Token{RANGE, lit, l.pos})
					break
				case "HOLD":
					result = append(result, Token{HOLD, lit, l.pos})
					break
				case "RELEASE":
					result = append(result, Token{RELEASE, lit, l.pos})
					break
				case "UPDATE":
					result = append(result, Token{UPDATE, lit, l.pos})
					break
				case "DELETE":
					result = append(result, Token{DELETE, lit, l.pos})
					break
				case "PUT":
					result = append(result, Token{PUT, lit, l.pos})
					break
				default:
					result = append(result, Token{RELATION, lit, l.pos})
					break
				}
			} else if r == '"' {
				l.backup()
				lit, period := l.lexStr()
				if period {
					result = append(result, Token{ILLEGAL, lit, l.pos})
					break
				}

				result = append(result, Token{CONSTANT, lit, l.pos})
				break
			} else {
				result = append(result, Token{ILLEGAL, string(r), l.pos})
				break
			}
		}
	}
}

func (l *Lexer) nextLine() {
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
	special := []rune{'.', '-', '/', '\\'}
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
		} else if r == '"' {
			quoteCount++
		} else {
			l.backup()
			return lit, period
		}
	}

}
