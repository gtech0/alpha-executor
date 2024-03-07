package service

import (
	"alpha-executor/model"
	"fmt"
	"log"
)

type Parser struct {
	tokens []model.Token
}

func NewParser(tokens []model.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Next() model.Token {
	prev := p.tokens[0]
	p.tokens = p.tokens[1:]
	return prev
}

func (p *Parser) Current() model.Token {
	return p.tokens[0]
}

func (p *Parser) Expect(lexType model.LexType) model.Token {
	prev := p.Next()
	if prev.Type == model.EOF || prev.Type != lexType {
		log.Fatal(fmt.Sprint("Unexpected token type ", lexType))
	}

	return prev
}

//func (p *Parser) ParseExpression() model.LexType {
//
//}
