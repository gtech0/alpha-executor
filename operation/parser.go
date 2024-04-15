package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"fmt"
	"log"
	"slices"
)

type Parser struct {
	tokens []*model.Token
}

func NewParser(tokens []*model.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) next() model.Token {
	if len(p.tokens) == 0 {
		panic("No tokens in parser")
	}

	prev := p.tokens[0]
	p.tokens = p.tokens[1:]
	return *prev
}

func (p *Parser) peek() model.Token {
	if len(p.tokens) == 0 {
		return model.Token{}
	}
	return *p.tokens[0]
}

func (p *Parser) expect(lexType model.LexType) model.Token {
	prev := p.next()
	if prev.Type == model.EOF || prev.Type != lexType {
		panic(fmt.Sprintf("Unexpected token type %s", lexType.String()))
	}

	return prev
}

func (p *Parser) ParseFullExpression() Expression {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Expression {
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
		left = &BinaryExpression{
			kind:     operator,
			left:     left,
			right:    right,
			position: p.peek().Position,
		}
	}

	return left
}

func (p *Parser) parseDisjunction() Expression {
	left := p.parseConjunction()

	for p.peek().Type == model.DISJUNCTION {
		operator := p.next().Value
		right := p.parseConjunction()
		left = &BinaryExpression{
			kind:     operator,
			left:     left,
			right:    right,
			position: p.peek().Position,
		}
	}

	return left
}

func (p *Parser) parseConjunction() Expression {
	left := p.parseComparison()

	for p.peek().Type == model.CONJUNCTION {
		operator := p.next().Value
		right := p.parseComparison()
		left = &BinaryExpression{
			kind:     operator,
			left:     left,
			right:    right,
			position: p.peek().Position,
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
		left = &BinaryExpression{
			kind:     operator,
			left:     left,
			right:    right,
			position: p.peek().Position,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	parsedType := p.peek().Type
	parsedValue := p.peek().Value
	position := p.peek().Position
	switch parsedType {
	case model.ATTRIBUTE, model.FREE_RELATION, model.BIND_RELATION, model.CONSTANT, model.INTEGER, model.DATE:
		token := p.next()
		return &IdentifierExpression{parsedType.String(), token.Value, token.Position}
	case model.EXISTS, model.FOR_ALL:
		p.next()
		return &BinaryExpression{parsedType.String(), p.parsePrimary(), p.parseComparison(), position}
	case model.NEGATION:
		p.next()
		return &UnaryExpression{parsedType.String(), p.parseComparison(), position}
	case model.LEFT_PARENTHESIS:
		p.next()
		value := p.parseExpression()
		p.expect(model.RIGHT_PARENTHESIS)
		return value
	case model.GET:
		p.next()
		variable := p.parsePrimary()
		rows, relations := p.parseRowNumAndRelation()
		expression := p.parseExpression()
		sort := p.parseSort()
		return &GetExpression{parsedType.String(), variable, rows, relations,
			expression, sort, position}
	case model.COMMA:
		p.next()
		return p.parsePrimary()
	case model.RANGE:
		p.next()
		return &RangeExpression{parsedType.String(), p.parsePrimary(), p.parsePrimary(), position}
	case model.HOLD:
		p.next()
		variable := p.parsePrimary()
		_, relations := p.parseRowNumAndRelation()
		return &HoldExpression{parsedType.String(), variable, relations,
			p.parseExpression(), position}
	case model.RELEASE, model.UPDATE, model.DELETE:
		p.next()
		return &OperationExpression{parsedType.String(), p.parsePrimary(), position}
	case model.PUT:
		p.next()
		variable := p.parsePrimary()
		_, relations := p.parseRowNumAndRelation()
		return &PutExpression{parsedType.String(), variable, relations, position}
	case model.LOGIC_START:
		p.next()
		return p.parseExpression()
	case model.DOWN, model.UP:
		p.next()
		return &UnaryExpression{parsedType.String(), p.parsePrimary(), position}
	default:
		panic(fmt.Sprintf("Unhandled type %s with value %s", parsedType.String(), parsedValue))
	}
}

func (p *Parser) parseRowNumAndRelation() (Expression, []Expression) {
	if p.peek().Type != model.LEFT_PARENTHESIS {
		panic(fmt.Sprintf("Unexpected type %s", p.peek().Type))
	}

	row := p.parseRelations()
	if len(row) > 0 && row[0].GetKind() != model.INTEGER.String() {
		return &IdentifierExpression{model.NULL.String(), model.NULL.String(), p.peek().Position}, row
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

func (p *Parser) parseSort() Expression {
	if len(p.tokens) == 0 {
		return &UnaryExpression{
			kind:       model.NULL.String(),
			expression: nil,
			position:   entity.Position{},
		}
	}

	return p.parsePrimary()
}
