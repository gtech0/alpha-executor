package service

import (
	"alpha-executor/model"
	"fmt"
	"io"
)

type Parser struct {
	tokens []model.Token
}

func (p *Parser) Next() model.Token {
	prev := p.tokens[0]
	p.tokens = p.tokens[1:]
	return prev
}

func (p *Parser) Current() model.Token {
	return p.tokens[0]
}

func (p *Parser) ParseStatement() model.LexType {
	return p.ParseExpression()
}

func (p *Parser) ParseExpression() model.LexType {
	return p.ParseComparisonExpression()
}

func (p *Parser) ParseComparisonExpression() model.LexType {
	return 0
}

func GenerateAST(reader io.Reader) {
	lexer := model.NewLexer(reader)
	output := Program{make([]model.Token, 0)}
	tempToken := model.Token{Value: ""}
	for {
		_, token := lexer.Lex()
		if token.Value == "∀" || token.Value == "∃" {
			tempToken = token
			continue
		}

		output.body = append(output.body, token)
		if tempToken.Value != "" {
			output.body = append(output.body, tempToken)
			tempToken.Value = ""
		}

		if token.Type == model.EOF {
			break
		}

		//fmt.Printf("%d:%d\t%s\t%d\n", pos.Line, pos.Column, token.Value, token.Type)
	}

	for _, token := range output.body {
		fmt.Printf("%s\t%d\n", token.Value, token.Type)
	}
}
