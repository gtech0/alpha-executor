package service

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"log"
	"strconv"
)

type Comparison struct {
	repository *repository.TestingRepository
}

func NewComparison(repository *repository.TestingRepository) *Comparison {
	return &Comparison{repository: repository}
}

func (c *Comparison) Compare(params BinaryExpression) (*entity.Relation, error) {
	result := &entity.Relation{}
	attr := model.Attribute{}

	attributeLeft, err := attr.ExtractAttribute(params.left.(IdentifierExpression).value)
	if err != nil {
		return nil, err
	}

	attributeRight, err := attr.ExtractAttribute(params.right.(IdentifierExpression).value)
	if err != nil {
		return nil, err
	}

	relation1, err := c.repository.GetRelation(attributeLeft.Relation)
	if err != nil {
		return nil, err
	}

	relation2, err := c.repository.GetRelation(attributeRight.Relation)
	if err != nil {
		return nil, err
	}

	for row1 := range *relation1 {
		if c.incorrectAttribute(row1, attributeLeft.Attribute) {
			return nil, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   fmt.Sprintf("incorrect attribute %s", attributeLeft.Attribute),
			}
		}

		for row2 := range *relation2 {
			if c.incorrectAttribute(row2, attributeRight.Attribute) {
				return nil, &entity.CustomError{
					ErrorType: entity.ResponseTypes["CE"],
					Message:   fmt.Sprintf("incorrect attribute %s", attributeRight.Attribute),
				}
			}
			valuesLeft := (*row1)[attributeLeft.Attribute]
			valuesRight := (*row2)[attributeRight.Attribute]
			c.checkIfAttributeAndCompare(row1, row2, valuesLeft, valuesRight, result, params)
		}
	}

	//attributes := []string{attributeLeft.Attribute, attributeRight.Attribute}
	//for row := range *relation.Right {
	//	for _, attribute := range attributes {
	//		if c.incorrectAttribute(row, attribute) {
	//			return &entity.CustomError{
	//				ErrorType: entity.ResponseTypes["CE"],
	//				Message:   fmt.Sprintf("incorrect attribute %s", attribute),
	//			}
	//		}
	//	}
	//	valuesLeft := (*row)[attributeLeft.Attribute]
	//	valuesRight := (*row)[attributeRight.Attribute]
	//	c.checkIfAttributeAndCompare(row, valuesLeft, valuesRight, result, params)
	//}

	//c.repository.AddRelation(relation.Left, relation.Right)
	return result, nil
}

func (*Comparison) incorrectAttribute(row *entity.RowMap, attribute string) bool {
	_, exists := (*row)[attribute]
	return !exists
}

func (c *Comparison) checkIfAttributeAndCompare(
	row1, row2 *entity.RowMap,
	valuesLeft, valuesRight []string,
	result *entity.Relation,
	params BinaryExpression,
) {
	if len(valuesLeft) != 0 && len(valuesRight) != 0 {
		for _, valLeft := range valuesLeft {
			for _, valRight := range valuesRight {
				//TODO: add interface methods to change values
				params.left.(*IdentifierExpression).value = valLeft
				params.right.(*IdentifierExpression).value = valRight
				c.valueComparator(params, row1, row2, result)
			}
		}
	} else if len(valuesLeft) != 0 {
		for _, valLeft := range valuesLeft {
			params.left.(*IdentifierExpression).value = valLeft
			c.valueComparator(params, row1, row2, result)
		}
	} else if len(valuesRight) != 0 {
		for _, valRight := range valuesRight {
			params.right.(*IdentifierExpression).value = valRight
			c.valueComparator(params, row1, row2, result)
		}
	} else {
		c.valueComparator(params, row1, row2, result)
	}
}

func (c *Comparison) valueComparator(params BinaryExpression, row1, row2 *entity.RowMap, result *entity.Relation) {
	//comparatorTokens := []string{">", "<", "≥", "≤", "=", "≠"}
	if c.numericComparator(params) || c.stringComparator(params) {
		(*result)[row1] = struct{}{}
		(*result)[row2] = struct{}{}
	}
}

//func (c *Comparison) compareTypes(params BinaryExpression) bool {
//	return c.numericComparator(params) || c.stringComparator(params) //|| c.dateComparator(params)
//}

func (*Comparison) numericComparator(params BinaryExpression) bool {
	oldValNum, err := strconv.ParseFloat(params.left.(IdentifierExpression).value, 10)
	if err != nil {
		return false
	}

	newValNum, err := strconv.ParseFloat(params.right.(IdentifierExpression).value, 10)
	if err != nil {
		return false
	}

	switch params.kind {
	case "=":
		return oldValNum == newValNum
	case "≠":
		return oldValNum != newValNum
	case "≤":
		return oldValNum <= newValNum
	case "≥":
		return oldValNum >= newValNum
	case "<":
		return oldValNum < newValNum
	case ">":
		return oldValNum > newValNum
	default:
		log.Fatal(fmt.Sprintf("Unknown operator %s", params.kind))
		return false
	}
}

func (*Comparison) stringComparator(params BinaryExpression) bool {
	oldVal := params.left.(IdentifierExpression).value
	newVal := params.right.(IdentifierExpression).value
	switch params.kind {
	case "=":
		return oldVal == newVal
	case "≠":
		return oldVal != newVal
	case "≤":
		return oldVal <= newVal
	case "≥":
		return oldVal >= newVal
	case "<":
		return oldVal < newVal
	case ">":
		return oldVal > newVal
	default:
		log.Fatal(fmt.Sprintf("Unknown operator %s", params.kind))
		return false
	}
}

//func (*Comparison) dateComparator(params BinaryExpression) bool {
//	if entity.IsQuoted(params.valLeft) {
//		params.valLeft = removeQuotes(params.valLeft)
//	}
//
//	if entity.IsQuoted(params.valRight) {
//		params.valRight = removeQuotes(params.valRight)
//	}
//
//	oldValDate, err := time.Parse(time.DateTime, params.valLeft)
//	if err != nil {
//		return false
//	}
//
//	newValDate, err := time.Parse(time.DateTime, params.valRight)
//	if err != nil {
//		return false
//	}
//
//	switch params.operator {
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
