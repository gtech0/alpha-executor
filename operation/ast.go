package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"bufio"
)

type Expression interface {
	GetKind() string
}

type Program struct {
	kind string
	body []Expression
}

func (p *Program) GetKind() string {
	return p.kind
}

type BinaryExpression struct {
	kind     string
	left     Expression
	right    Expression
	position entity.Position
}

func (b *BinaryExpression) GetKind() string {
	return b.kind
}

type UnaryExpression struct {
	kind     string
	right    Expression
	position entity.Position
}

func (u *UnaryExpression) GetKind() string {
	return u.kind
}

type IdentifierExpression struct {
	kind     string
	value    string
	position entity.Position
}

func (i *IdentifierExpression) GetKind() string {
	return i.kind
}

type GetExpression struct {
	kind       string
	variable   Expression
	rows       Expression
	relations  []Expression
	expression Expression
	position   entity.Position
}

func (g *GetExpression) GetKind() string {
	return g.kind
}

type RangeExpression struct {
	kind     string
	relation Expression
	variable Expression
	position entity.Position
}

func (r *RangeExpression) GetKind() string {
	return r.kind
}

type HoldExpression struct {
	kind       string
	variable   Expression
	relations  []Expression
	expression Expression
	position   entity.Position
}

func (h *HoldExpression) GetKind() string {
	return h.kind
}

type OperationExpression struct {
	kind     string
	variable Expression
	position entity.Position
}

func (s *OperationExpression) GetKind() string {
	return s.kind
}

type PutExpression struct {
	kind      string
	variable  Expression
	relations []Expression
	position  entity.Position
}

func (p *PutExpression) GetKind() string {
	return p.kind
}

func GenerateAST(reader *bufio.Reader) Program {
	program := Program{model.PROGRAM.String(), make([]Expression, 0)}
	lexer := model.NewLexer(reader)
	output := lexer.Lex()
	for _, query := range output {
		if len(query) > 0 {
			//for _, token := range query {
			//	fmt.Printf("%d:%d\t%s\t%s\n", token.Position.Line, token.Position.Column, token.Type.String(), token.Value)
			//}
			parser := NewParser(query)
			program.body = append(program.body, parser.ParseExpression())
		}
	}

	return program
}
