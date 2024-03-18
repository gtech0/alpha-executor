package service

import (
	"alpha-executor/model"
	"fmt"
	"log"
	"slices"
)

type Parser struct {
	tokens []model.Token
}

func NewParser(tokens []model.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) next() model.Token {
	prev := p.tokens[0]
	p.tokens = p.tokens[1:]
	return prev
}

func (p *Parser) peek() model.Token {
	if len(p.tokens) == 0 {
		return model.Token{}
	}
	return p.tokens[0]
}

func (p *Parser) expect(lexType model.LexType) model.Token {
	prev := p.next()
	if prev.Type == model.EOF || prev.Type != lexType {
		log.Fatal(fmt.Sprint("Unexpected token type ", lexType))
	}

	return prev
}

func (p *Parser) ParseExpression() Expression {
	switch p.peek().Type {
	default:
		return p.parseImplication()
	}
}

func (p *Parser) parseImplication() Expression {
	left := p.parseDisjunction()

	for p.peek().Type == model.IMPLICATION {
		operator := p.next().Value
		right := p.parseDisjunction()
		left = BinaryExpression{
			kind:     BINARY,
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

func (p *Parser) parseDisjunction() Expression {
	left := p.parseConjunction()

	for p.peek().Type == model.DISJUNCTION {
		operator := p.next().Value
		right := p.parseConjunction()
		left = BinaryExpression{
			kind:     BINARY,
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

func (p *Parser) parseConjunction() Expression {
	left := p.parseComparison()

	for p.peek().Type == model.CONJUNCTION {
		operator := p.next().Value
		right := p.parseComparison()
		left = BinaryExpression{
			kind:     BINARY,
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

func (p *Parser) parseComparison() Expression {
	left := p.parsePrimary()

	operators := []model.LexType{
		model.EQUALS,
		model.NOT_EQUALS,
		model.LESS_THAN,
		model.LESS_THAN_EQUALS,
		model.GREATER_THAN,
		model.GREATER_THAN_EQUALS,
	}

	for slices.Contains(operators, p.peek().Type) {
		operator := p.next().Value
		right := p.parsePrimary()
		left = BinaryExpression{
			kind:     BINARY,
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	switch p.peek().Type {
	case model.ATTRIBUTE:
		return Identifier{ATTRIBUTE, p.next().Value}
	case model.RELATION:
		return Identifier{RELATION, p.next().Value}
	case model.CONSTANT:
		return Identifier{CONSTANT, p.next().Value}
	case model.INTEGER:
		return Identifier{INTEGER, p.next().Value}
	case model.NEGATION:
		p.next()
		return UnaryExpression{NEGATION, p.parsePrimary(), "¬"}
	case model.EXIST:
		p.next()
		return BinaryExpression{EXIST, p.parsePrimary(), p.parseComparison(), "∃"}
	case model.FOR_ALL:
		p.next()
		return BinaryExpression{FOR_ALL, p.parsePrimary(), p.parseComparison(), "∀"}
	case model.LEFT_PARENTHESIS:
		p.next()
		value := p.ParseExpression()
		p.expect(model.RIGHT_PARENTHESIS)
		return value
	default:
		panic(fmt.Sprint("Unhandled default case ", p.peek().Type))
	}
}
