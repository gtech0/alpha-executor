package service

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

type Comparison struct {
}

func (c *Comparison) compareOperands(
	relation *entity.Relation,
	params BinaryExpression,
	result *entity.Relation,
) error {
	attr := model.Attribute{}

	attributeLeft, err := attr.ExtractAttribute(params.left.(IdentifierExpression).token.Value)
	if err != nil {
		return err
	}

	attributeRight, err := attr.ExtractAttribute(params.right.(IdentifierExpression).token.Value)
	if err != nil {
		return err
	}

	attributes := []string{attributeLeft.Attribute, attributeRight.Attribute}
	for row := range *relation {
		for _, attribute := range attributes {
			if c.incorrectAttribute(row, attribute) {
				return &entity.CustomError{
					ErrorType: entity.ResponseTypes["CE"],
					Message:   fmt.Sprintf("incorrect attribute %s", attribute),
				}
			}
		}
		valuesLeft := (*row)[attributeLeft.Attribute]
		valuesRight := (*row)[attributeRight.Attribute]
		c.checkIfAttributeAndCompare(row, valuesLeft, valuesRight, result, params)
	}
	//TODO: add to relations

	return nil
}

func (*Comparison) incorrectAttribute(row *entity.RowMap, attribute string) bool {
	_, exists := (*row)[attribute]
	return !entity.IsQuoted(attribute) && !entity.IsNumeric(attribute) && !exists
}

func (c *Comparison) checkIfAttributeAndCompare(
	row *entity.RowMap,
	valuesLeft, valuesRight []string,
	result *entity.Relation,
	params BinaryExpression,
) {
	if len(valuesLeft) != 0 && len(valuesRight) != 0 {
		for _, valLeft := range valuesLeft {
			for _, valRight := range valuesRight {
				params.left.(*IdentifierExpression).token.Value = valLeft
				params.right.(*IdentifierExpression).token.Value = valRight
				c.valueComparator(params, row, result)
			}
		}
	} else if len(valuesLeft) != 0 {
		for _, valLeft := range valuesLeft {
			params.left.(*IdentifierExpression).token.Value = valLeft
			c.valueComparator(params, row, result)
		}
	} else if len(valuesRight) != 0 {
		for _, valRight := range valuesRight {
			params.right.(*IdentifierExpression).token.Value = valRight
			c.valueComparator(params, row, result)
		}
	} else {
		c.valueComparator(params, row, result)
	}
}

func (c *Comparison) valueComparator(params BinaryExpression, row *entity.RowMap, result *entity.Relation) {
	comparatorTokens := []string{">", "<", "≥", "≤", "=", "≠"}
	if slices.Contains(comparatorTokens, params.kind) && c.compare(params) {
		(*result)[row] = struct{}{}
	}
}

func (c *Comparison) compare(params BinaryExpression) bool {
	return c.numericComparator(params) || c.stringComparator(params) //|| c.dateComparator(params)
}

func (*Comparison) numericComparator(params BinaryExpression) bool {
	oldValNum, err := strconv.ParseFloat(params.left.(IdentifierExpression).token.Value, 10)
	if err != nil {
		return false
	}

	newValNum, err := strconv.ParseFloat(params.right.(IdentifierExpression).token.Value, 10)
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
	oldVal := params.left.(IdentifierExpression).token.Value
	newVal := params.right.(IdentifierExpression).token.Value
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

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}
