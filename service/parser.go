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
		log.Fatal(fmt.Sprintf("Unexpected token type %s", lexType.String()))
	}

	return prev
}

func (p *Parser) ParseExpression() Expression {
	switch p.peek().Type {
	case model.GET, model.RANGE, model.HOLD, model.RELEASE, model.UPDATE, model.DELETE, model.PUT:
		return p.parsePrimary()
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
			kind:     model.BINARY.String(),
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
			kind:     model.BINARY.String(),
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
			kind:     model.BINARY.String(),
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
			kind:     model.BINARY.String(),
			left:     left,
			right:    right,
			operator: operator,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	parsedType := p.peek().Type
	switch parsedType {
	case model.ATTRIBUTE:
		return Identifier{model.ATTRIBUTE.String(), p.next().Value}
	case model.RELATION:
		return Identifier{model.RELATION.String(), p.next().Value}
	case model.CONSTANT:
		return Identifier{model.CONSTANT.String(), p.next().Value}
	case model.INTEGER:
		return Identifier{model.INTEGER.String(), p.next().Value}
	case model.NEGATION:
		p.next()
		return UnaryExpression{model.NEGATION.String(), p.parsePrimary(), "¬"}
	case model.EXIST:
		p.next()
		return BinaryExpression{model.EXIST.String(), p.parsePrimary(), p.parseComparison(), "∃"}
	case model.FOR_ALL:
		p.next()
		return BinaryExpression{model.FOR_ALL.String(), p.parsePrimary(), p.parseComparison(), "∀"}
	case model.LEFT_PARENTHESIS:
		p.next()
		value := p.ParseExpression()
		p.expect(model.RIGHT_PARENTHESIS)
		return value
	case model.GET:
		p.next()
		variable := p.parsePrimary()
		rows, relations := p.parseRowNumAndRelation()
		return Get{model.GET.String(), variable, rows, relations, p.ParseExpression()}
	case model.COMMA:
		p.next()
		return p.parsePrimary()
	case model.RANGE:
		p.next()
		return Range{model.RANGE.String(), p.parsePrimary(), p.parsePrimary()}
	case model.HOLD:
		p.next()
		variable := p.parsePrimary()
		_, relations := p.parseRowNumAndRelation()
		return Hold{model.HOLD.String(), variable, relations, p.ParseExpression()}
	case model.RELEASE, model.UPDATE, model.DELETE:
		kind := parsedType.String()
		p.next()
		return SimpleOperation{kind, p.parsePrimary()}
	case model.PUT:
		p.next()
		variable := p.parsePrimary()
		_, relations := p.parseRowNumAndRelation()
		return Put{model.PUT.String(), variable, relations}
	case model.LOGIC_START:
		p.next()
		return p.ParseExpression()
	default:
		panic(fmt.Sprintf("Unhandled default case %s", parsedType.String()))
	}
}

func (p *Parser) parseRowNumAndRelation() (Expression, []Expression) {
	if p.peek().Type != model.LEFT_PARENTHESIS {
		panic(fmt.Sprintf("Unexpected type %s", p.peek().Type))
	}

	row := p.parseRelations()
	if len(row) > 0 && row[0].GetKind() != model.INTEGER.String() {
		return Identifier{model.INTEGER.String(), "0"}, row
	} else if len(row) == 0 {
		log.Fatal(fmt.Sprintf("No parameters detected on %d:%d", p.peek().Position.Line, p.peek().Position.Column))
	}

	relations := p.parseRelations()
	return row[0], relations
}

func (p *Parser) parseRelations() []Expression {
	p.next()
	relations := []Expression{p.parsePrimary()}
	for p.peek().Type == model.COMMA {
		relations = append(relations, p.parsePrimary())
	}
	p.expect(model.RIGHT_PARENTHESIS)
	return relations
}
