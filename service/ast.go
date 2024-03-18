package service

import (
	"alpha-executor/model"
	"bufio"
	"fmt"
	"github.com/kr/pretty"
)

type Kind int

type Expression interface {
	GetKind() Kind
}

const (
	PROGRAM Kind = iota
	BINARY
	ATTRIBUTE
	RELATION
	CONSTANT
	INTEGER
	FOR_ALL
	EXIST
	NEGATION
)

type Program struct {
	kind Kind
	body Expression
}

func (p Program) GetKind() Kind {
	return p.kind
}

type BinaryExpression struct {
	kind     Kind
	left     Expression
	right    Expression
	operator string
}

func (b BinaryExpression) GetKind() Kind {
	return b.kind
}

type UnaryExpression struct {
	kind     Kind
	right    Expression
	operator string
}

func (u UnaryExpression) GetKind() Kind {
	return u.kind
}

type Identifier struct {
	kind  Kind
	value string
}

func (i Identifier) GetKind() Kind {
	return i.kind
}

func GenerateAST(reader *bufio.Reader) {
	for {
		lexer := model.NewLexer(reader)
		output, currReader := lexer.Lex()
		if len(output) > 0 {
			parser := NewParser(output)
			program := Program{PROGRAM, parser.ParseExpression()}
			pretty.Print(program)
			fmt.Print("\n")
		}

		if currReader.Size() == 0 {
			break
		}
	}

	//for _, token := range output {
	//	fmt.Printf("%d:%d\t%d\t%s\n", token.Position.Line, token.Position.Column, token.Type, token.Value)
	//}
}
