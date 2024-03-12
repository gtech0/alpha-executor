package service

import (
	"alpha-executor/model"
	"github.com/kr/pretty"
	"io"
)

type Expression interface {
	getKind() string
}

type BinaryExpression struct {
	kind     string
	left     Expression
	right    Expression
	operator string
}

func (b BinaryExpression) getKind() string {
	return b.kind
}

type Identifier struct {
	kind  string
	value string
}

func (b Identifier) getKind() string {
	return b.kind
}

type UnaryExpression struct {
	kind     string
	right    Expression
	operator string
}

func (b UnaryExpression) getKind() string {
	return b.kind
}

func GenerateAST(reader io.Reader) {
	lexer := model.NewLexer(reader)
	output := lexer.Lex()
	parser := NewParser(output)
	pretty.Print(parser.ParseExpression())
	//for _, token := range output {
	//	fmt.Printf("%d:%d\t%d\t%s\n", token.Position.Line, token.Position.Column, token.Type, token.Value)
	//}
}
