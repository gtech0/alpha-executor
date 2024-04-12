package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"strconv"
	"time"
)

type Comparison struct {
	repository *repository.TestingRepository
	parameters *parameters
}

type parameters struct {
	kind     string
	left     IdentifierExpression
	right    IdentifierExpression
	position entity.Position
}

func NewComparison(repository *repository.TestingRepository) *Comparison {
	return &Comparison{repository: repository}
}

func (c *Comparison) Compare(params *BinaryExpression) (bool, error) {
	c.parameters = &parameters{
		kind:     params.kind,
		left:     *params.left.(*IdentifierExpression),
		right:    *params.right.(*IdentifierExpression),
		position: params.position,
	}

	attr := model.Attribute{}

	attributeLeft, err := attr.ExtractAttribute(c.parameters.left.value, c.parameters.position)
	if err != nil {
		return false, err
	}

	attributeRight, err := attr.ExtractAttribute(c.parameters.right.value, c.parameters.position)
	if err != nil {
		return false, err
	}

	if c.parameters.left.kind == model.ATTRIBUTE.String() && c.parameters.right.kind == model.ATTRIBUTE.String() {
		return c.twoAttributesCompare(attributeLeft, attributeRight)
	} else if c.parameters.left.kind == model.ATTRIBUTE.String() {
		return c.oneAttributeCompare(attributeLeft)
	} else if c.parameters.right.kind == model.ATTRIBUTE.String() {
		return c.oneAttributeCompare(attributeRight)
	}

	return false, nil
}

func (c *Comparison) twoAttributesCompare(attributeLeft, attributeRight model.ComplexAttribute) (bool, error) {
	relation1, err := c.repository.GetRow(attributeLeft.Relation)
	if err != nil {
		err.(*entity.CustomError).Position = c.parameters.position
		return false, err
	}

	relation2, err := c.repository.GetRow(attributeRight.Relation)
	if err != nil {
		err.(*entity.CustomError).Position = c.parameters.position
		return false, err
	}

	if c.incorrectAttribute(relation1, attributeLeft.Attribute) {
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprintf("incorrect attribute %s", attributeLeft.Attribute),
			Position:  c.parameters.position,
		}
	}

	if c.incorrectAttribute(relation2, attributeRight.Attribute) {
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprintf("incorrect attribute %s", attributeRight.Attribute),
			Position:  c.parameters.position,
		}
	}

	valuesLeft := (*relation1)[attributeLeft.Attribute]
	valuesRight := (*relation2)[attributeRight.Attribute]
	for _, valLeft := range valuesLeft {
		for _, valRight := range valuesRight {
			c.parameters.left.value = valLeft
			c.parameters.right.value = valRight
			if isTrue, err := c.valueComparator(); isTrue && err == nil {
				return true, nil
			}

			if err != nil {
				return false, err
			}
		}
	}

	return false, nil
}

func (c *Comparison) oneAttributeCompare(attribute model.ComplexAttribute) (bool, error) {
	relation, err := c.repository.GetRow(attribute.Relation)
	if err != nil {
		err.(*entity.CustomError).Position = c.parameters.position
		return false, err
	}

	if c.incorrectAttribute(relation, attribute.Attribute) {
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprintf("incorrect attribute %s", attribute.Attribute),
			Position:  c.parameters.position,
		}
	}

	values := (*relation)[attribute.Attribute]
	for _, value := range values {
		c.parameters.left.value = value
		if isTrue, err := c.valueComparator(); isTrue && err == nil {
			return true, nil
		} else if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (*Comparison) incorrectAttribute(row *entity.RowMap, attribute string) bool {
	_, exists := (*row)[attribute]
	return !exists
}

func (c *Comparison) valueComparator() (bool, error) {
	numeric, err := c.numericComparator()
	if err != nil {
		return false, err
	}

	str, err := c.stringComparator()
	if err != nil {
		return false, err
	}

	date, err := c.dateComparator()
	if err != nil {
		return false, err
	}

	return numeric || str || date, nil
}

func (c *Comparison) numericComparator() (bool, error) {
	oldVal, err := strconv.ParseFloat(c.parameters.left.value, 10)
	if err != nil {
		return false, nil
	}

	newVal, err := strconv.ParseFloat(c.parameters.right.value, 10)
	if err != nil {
		return false, nil
	}

	switch c.parameters.kind {
	case "=":
		return oldVal == newVal, nil
	case "≠":
		return oldVal != newVal, nil
	case "≤":
		return oldVal <= newVal, nil
	case "≥":
		return oldVal >= newVal, nil
	case "<":
		return oldVal < newVal, nil
	case ">":
		return oldVal > newVal, nil
	default:
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprintf("Unknown operator %s", c.parameters.kind),
			Position:  c.parameters.position,
		}
	}
}

func (c *Comparison) stringComparator() (bool, error) {
	oldVal := c.parameters.left.value
	newVal := c.parameters.right.value
	switch c.parameters.kind {
	case "=":
		return oldVal == newVal, nil
	case "≠":
		return oldVal != newVal, nil
	case "≤":
		return oldVal <= newVal, nil
	case "≥":
		return oldVal >= newVal, nil
	case "<":
		return oldVal < newVal, nil
	case ">":
		return oldVal > newVal, nil
	default:
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprintf("Unknown operator %s", c.parameters.kind),
			Position:  c.parameters.position,
		}
	}
}

func (c *Comparison) dateComparator() (bool, error) {
	oldVal, err := time.Parse(time.DateTime, c.parameters.left.value)
	if err != nil {
		return false, nil
	}

	newVal, err := time.Parse(time.DateTime, c.parameters.right.value)
	if err != nil {
		return false, nil
	}

	switch c.parameters.kind {
	case "=":
		return oldVal.Equal(newVal), nil
	case "!=":
		return !oldVal.Equal(newVal), nil
	case "<=":
		return oldVal.Before(newVal) || oldVal.Equal(newVal), nil
	case ">=":
		return oldVal.After(newVal) || oldVal.Equal(newVal), nil
	case "<":
		return oldVal.Before(newVal), nil
	case ">":
		return oldVal.After(newVal), nil
	default:
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprintf("Unknown operator %s", c.parameters.kind),
			Position:  c.parameters.position,
		}
	}
}
