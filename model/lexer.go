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
	DATE
	NULL
	FREE_RELATION
	BIND_RELATION

	GET
	RANGE
	HOLD
	RELEASE
	UPDATE
	DELETE
	PUT

	DOWN
	UP

	ASSIGN

	EQUALS
	NOT_EQUALS
	LESS_THAN
	LESS_THAN_EQUALS
	GREATER_THAN
	GREATER_THAN_EQUALS

	EXISTS
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

	ILLEGAL:       "ILLEGAL",
	ATTRIBUTE:     "ATTRIBUTE",
	CONSTANT:      "CONSTANT",
	INTEGER:       "INTEGER",
	DATE:          "DATE",
	NULL:          "NULL",
	FREE_RELATION: "FREE_RELATION",
	BIND_RELATION: "BIND_RELATION",

	GET:     "GET",
	RANGE:   "RANGE",
	HOLD:    "HOLD",
	RELEASE: "RELEASE",
	UPDATE:  "UPDATE",
	DELETE:  "DELETE",
	PUT:     "PUT",

	DOWN: "DOWN",
	UP:   "UP",

	ASSIGN: "ASSIGN",

	EQUALS:              "=",
	NOT_EQUALS:          "≠",
	LESS_THAN:           "<",
	LESS_THAN_EQUALS:    "≤",
	GREATER_THAN:        ">",
	GREATER_THAN_EQUALS: "≥",

	EXISTS:  "∃", //какие-то кортежи удовлетворяют условию
	FOR_ALL: "∀", //     все

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
	results [][]*Token
}

func NewLexer(reader *bufio.Reader) *Lexer {
	return &Lexer{
		pos:     entity.Position{Line: 1, Column: 0},
		reader:  reader,
		results: make([][]*Token, 0),
	}
}

func (l *Lexer) Lex() [][]*Token {
	result := make([]*Token, 0)
	for {
		r, _, err := l.reader.ReadRune()
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
			result = make([]*Token, 0)
			continue
		}

		switch r {
		case '\n':
			l.nextLine()
		case '=':
			result = append(result, &Token{EQUALS, EQUALS.String(), l.pos})
			break
		case '≠':
			result = append(result, &Token{NOT_EQUALS, NOT_EQUALS.String(), l.pos})
			break
		case '<':
			result = append(result, &Token{LESS_THAN, LESS_THAN.String(), l.pos})
			break
		case '≤':
			result = append(result, &Token{LESS_THAN_EQUALS, LESS_THAN_EQUALS.String(), l.pos})
			break
		case '>':
			result = append(result, &Token{GREATER_THAN, GREATER_THAN.String(), l.pos})
			break
		case '≥':
			result = append(result, &Token{GREATER_THAN_EQUALS, GREATER_THAN_EQUALS.String(), l.pos})
			break
		case '∃':
			result = append(result, &Token{EXISTS, EXISTS.String(), l.pos})
			break
		case '∀':
			result = append(result, &Token{FOR_ALL, FOR_ALL.String(), l.pos})
			break
		case '¬':
			result = append(result, &Token{NEGATION, NEGATION.String(), l.pos})
			break
		case '∧':
			result = append(result, &Token{CONJUNCTION, CONJUNCTION.String(), l.pos})
			break
		case '∨':
			result = append(result, &Token{DISJUNCTION, DISJUNCTION.String(), l.pos})
			break
		case '→':
			result = append(result, &Token{IMPLICATION, IMPLICATION.String(), l.pos})
			break
		case '(':
			result = append(result, &Token{LEFT_PARENTHESIS, LEFT_PARENTHESIS.String(), l.pos})
			break
		case ')':
			result = append(result, &Token{RIGHT_PARENTHESIS, RIGHT_PARENTHESIS.String(), l.pos})
			break
		case ',':
			result = append(result, &Token{COMMA, COMMA.String(), l.pos})
			break
		case ':':
			result = append(result, &Token{LOGIC_START, LOGIC_START.String(), l.pos})
			break
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				l.backup()
				lit := l.lexInt()
				result = append(result, &Token{INTEGER, lit, l.pos})
				break
			} else if unicode.IsLetter(r) {
				l.backup()
				lit, dot, dash := l.lexStr()
				if dot == 1 && dash == 0 {
					result = append(result, &Token{ATTRIBUTE, lit, l.pos})
					break
				} else if dot > 1 && dash != 0 {
					result = append(result, &Token{ILLEGAL, lit, l.pos})
					break
				}

				switch lit {
				case "GET":
					result = append(result, &Token{GET, lit, l.pos})
					break
				case "RANGE":
					result = append(result, &Token{RANGE, lit, l.pos})
					break
				case "HOLD":
					result = append(result, &Token{HOLD, lit, l.pos})
					break
				case "RELEASE":
					result = append(result, &Token{RELEASE, lit, l.pos})
					break
				case "UPDATE":
					result = append(result, &Token{UPDATE, lit, l.pos})
					break
				case "DELETE":
					result = append(result, &Token{DELETE, lit, l.pos})
					break
				case "PUT":
					result = append(result, &Token{PUT, lit, l.pos})
					break
				case "DOWN":
					result = append(result, &Token{DOWN, lit, l.pos})
					break
				case "UP":
					result = append(result, &Token{UP, lit, l.pos})
					break
				default:
					if len(result) > 1 && result[len(result)-1].Type == EXISTS {
						result = append(result, &Token{BIND_RELATION, lit, l.pos})
						for _, tokens := range l.results {
							for _, token := range tokens {
								if token.Type == FREE_RELATION && token.Value == lit {
									token.Type = BIND_RELATION
									break
								}
							}
						}
						break
					}

					result = append(result, &Token{FREE_RELATION, lit, l.pos})
					break
				}
			} else if r == '"' {
				l.backup()
				lit, dot, dash := l.lexStr()
				if dot > 0 || (dash > 0 && dash != 2) {
					result = append(result, &Token{ILLEGAL, lit, l.pos})
					break
				}

				if dash == 2 {
					result = append(result, &Token{DATE, lit, l.pos})
					break
				}

				result = append(result, &Token{CONSTANT, lit, l.pos})
				break
			} else {
				result = append(result, &Token{ILLEGAL, string(r), l.pos})
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
	lit := ""
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
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

func (l *Lexer) lexStr() (string, int, int) {
	lit := ""
	special := []rune{'.', '-', '/', '\\'}

	quote := 0
	dot := 0
	dash := 0

	for {
		if quote == 2 {
			return lit, dot, dash
		}

		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit, dot, dash
			}
		}

		l.pos.Column++
		if unicode.IsLetter(r) || unicode.IsDigit(r) || slices.Contains(special, r) {
			lit = lit + string(r)
			if r == '.' {
				dot += 1
			}

			if r == '-' {
				dash += 1
			}
		} else if r == '"' {
			quote++
		} else {
			l.backup()
			return lit, dot, dash
		}
	}

}
