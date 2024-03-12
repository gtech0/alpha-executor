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
		fmt.Print("WARNING: no tokens left\n")
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
	return p.parseDisjunction()
}

func (p *Parser) parseDisjunction() Expression {
	left := p.parseConjunction()

	for p.peek().Value == "∨" {
		operator := p.next().Value
		right := p.parseConjunction()
		left = BinaryExpression{
			kind:     "Binary expression",
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

func (p *Parser) parseConjunction() Expression {
	left := p.parseComparison()

	for p.peek().Value == "∧" {
		operator := p.next().Value
		right := p.parseComparison()
		left = BinaryExpression{
			kind:     "Binary expression",
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

//func (p *Parser) parseQuantifier() Expression {
//	left := p.parseComparison()
//
//	for p.peek().Value == "∃" || p.peek().Value == "∀" || p.peek().Value == "(" {
//		//operator := p.next().Value
//		right := p.parseComparison()
//		left = UnaryExpression{
//			kind:     "Unary expression",
//			right:    right,
//			operator: p.peek().Value,
//		}
//	}
//
//	return left
//}

func (p *Parser) parseComparison() Expression {
	left := p.parsePrimary()

	operators := []string{"=", "≠", "<", "≤", ">", "≥"}
	for slices.Contains(operators, p.peek().Value) {
		operator := p.next().Value
		right := p.parsePrimary()
		left = BinaryExpression{
			kind:     "Binary expression",
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
		return Identifier{"Attribute", p.next().Value}
	case model.RELATION:
		return Identifier{"Relation", p.next().Value}
	case model.CONST:
		return Identifier{"Constant", p.next().Value}
	case model.INT:
		return Identifier{"Integer", p.next().Value}
	case model.NOT:
		p.next()
		return UnaryExpression{"Negation", p.parsePrimary(), "¬"}
	case model.EXIST:
		p.next()
		return BinaryExpression{"Quantifier", p.parsePrimary(), p.parseComparison(), "∃"}
	case model.FOR_ALL:
		p.next()
		return BinaryExpression{"Quantifier", p.parsePrimary(), p.parseComparison(), "∀"}
	case model.LEFT_PARENTHESIS:
		p.next()
		value := p.ParseExpression()
		p.expect(model.RIGHT_PARENTHESIS)
		return value
	default:
		panic(fmt.Sprint("unhandled default case ", p.peek().Type))
	}
}
