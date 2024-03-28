package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"strconv"
)

type Comparison struct {
	repository *repository.TestingRepository
	parameters *parameters
}

type parameters struct {
	kind     string
	left     *IdentifierExpression
	right    *IdentifierExpression
	position entity.Position
}

func NewComparison(repository *repository.TestingRepository) *Comparison {
	return &Comparison{repository: repository}
}

func (c *Comparison) Compare(params *BinaryExpression) (*entity.Relation, error) {
	c.parameters = &parameters{
		kind:     params.kind,
		left:     params.left.(*IdentifierExpression),
		right:    params.right.(*IdentifierExpression),
		position: params.position,
	}

	result := &entity.Relation{}
	attr := model.Attribute{}

	attributeLeft, err := attr.ExtractAttribute(c.parameters.left.value, c.parameters.position)
	if err != nil {
		return nil, err
	}

	attributeRight, err := attr.ExtractAttribute(c.parameters.right.value, c.parameters.position)
	if err != nil {
		return nil, err
	}

	if c.parameters.left.kind == model.ATTRIBUTE.String() && c.parameters.right.kind == model.ATTRIBUTE.String() {
		if err := c.twoAttributesCompare(attributeLeft, attributeRight, result); err != nil {
			return nil, err
		}
	} else if c.parameters.left.kind == model.ATTRIBUTE.String() {
		if err := c.oneAttributeCompare(attributeLeft, result); err != nil {
			return nil, err
		}
	} else if c.parameters.right.kind == model.ATTRIBUTE.String() {
		if err := c.oneAttributeCompare(attributeRight, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Comparison) twoAttributesCompare(attributeLeft, attributeRight model.ComplexAttribute, result *entity.Relation) error {
	relation1, err := c.repository.GetRelation(attributeLeft.Relation)
	if err != nil {
		err.(*entity.CustomError).Position = c.parameters.position
		return err
	}

	relation2, err := c.repository.GetRelation(attributeRight.Relation)
	if err != nil {
		err.(*entity.CustomError).Position = c.parameters.position
		return err
	}

	for row1 := range *relation1 {
		if c.incorrectAttribute(row1, attributeLeft.Attribute) {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   fmt.Sprintf("incorrect attribute %s", attributeLeft.Attribute),
				Position:  c.parameters.position,
			}
		}

		for row2 := range *relation2 {
			if c.incorrectAttribute(row2, attributeRight.Attribute) {
				return &entity.CustomError{
					ErrorType: entity.ResponseTypes["CE"],
					Message:   fmt.Sprintf("incorrect attribute %s", attributeRight.Attribute),
					Position:  c.parameters.position,
				}
			}

			valuesLeft := (*row1)[attributeLeft.Attribute]
			valuesRight := (*row2)[attributeRight.Attribute]
			for _, valLeft := range valuesLeft {
				for _, valRight := range valuesRight {
					c.parameters.left.value = valLeft
					c.parameters.right.value = valRight
					if isTrue, err := c.valueComparator(); isTrue && err == nil {
						(*result)[row1] = struct{}{}
						(*result)[row2] = struct{}{}
					}

					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (c *Comparison) oneAttributeCompare(attribute model.ComplexAttribute, result *entity.Relation) error {
	relation, err := c.repository.GetRelation(attribute.Relation)
	if err != nil {
		err.(*entity.CustomError).Position = c.parameters.position
		return err
	}

	for row := range *relation {
		if c.incorrectAttribute(row, attribute.Attribute) {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   fmt.Sprintf("incorrect attribute %s", attribute.Attribute),
				Position:  c.parameters.position,
			}
		}

		values := (*row)[attribute.Attribute]
		for _, value := range values {
			c.parameters.left.value = value
			if isTrue, err := c.valueComparator(); isTrue && err == nil {
				(*result)[row] = struct{}{}
			} else if err != nil {
				return err
			}
		}
	}
	return nil
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

	//date, err := c.dateComparator()
	//if err != nil {
	//	return false, err
	//}

	return numeric || str /*|| date*/, nil
}

func (c *Comparison) numericComparator() (bool, error) {
	oldValNum, err := strconv.ParseFloat(c.parameters.left.value, 10)
	if err != nil {
		return false, nil
	}

	newValNum, err := strconv.ParseFloat(c.parameters.right.value, 10)
	if err != nil {
		return false, nil
	}

	switch c.parameters.kind {
	case "=":
		return oldValNum == newValNum, nil
	case "≠":
		return oldValNum != newValNum, nil
	case "≤":
		return oldValNum <= newValNum, nil
	case "≥":
		return oldValNum >= newValNum, nil
	case "<":
		return oldValNum < newValNum, nil
	case ">":
		return oldValNum > newValNum, nil
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

//func (*Comparison) dateComparator(parameters BinaryExpression) bool {
//	if entity.IsQuoted(parameters.valLeft) {
//		parameters.valLeft = removeQuotes(parameters.valLeft)
//	}
//
//	if entity.IsQuoted(parameters.valRight) {
//		parameters.valRight = removeQuotes(parameters.valRight)
//	}
//
//	oldValDate, err := time.Parse(time.DateTime, parameters.valLeft)
//	if err != nil {
//		return false
//	}
//
//	newValDate, err := time.Parse(time.DateTime, parameters.valRight)
//	if err != nil {
//		return false
//	}
//
//	switch parameters.operator {
//	case "=":
//		return oldValDate.Equal(newValDate)
//	case "!=":
//		return !oldValDate.Equal(newValDate)
//	case "<=":
//		return oldValDate.Before(newValDate) || oldValDate.Equal(newValDate)
//	case ">=":
//		return oldValDate.After(newValDate) || oldValDate.Equal(newValDate)
//	case "<":
//		return oldValDate.Before(newValDate)
//	case ">":
//		return oldValDate.After(newValDate)
//	}
//
//	return false
//}