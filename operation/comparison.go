package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Comparison struct {
	valLeft  string
	valRight string
	operator string
}

func (c *Comparison) compareOperands(
	relation entity.Pair[string, *entity.Relation],
	params Comparison,
	result *entity.Relation,
	results *entity.Stack[any],
) error {
	attr := model.Attribute{}
	attributeLeft, err := attr.ExtractAttribute(params.valLeft)
	if err != nil {
		return err
	}

	attributeRight, err := attr.ExtractAttribute(params.valRight)
	if err != nil {
		return err
	}

	attributes := []string{attributeLeft.Attribute, attributeRight.Attribute}
	for row := range *relation.Right {
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
	(*results).Push(result)
	return nil
}

func assertTypeAndPop[T any](results *entity.Stack[any]) T {
	operand, exists := (*results).Pop()
	if !exists {
		log.Fatal("stack shouldn't be empty")
	}

	operandRight, asserted := operand.(T)
	if !asserted {
		log.Fatal("type assertion failed")
	}

	return operandRight
}

func (*Comparison) incorrectAttribute(row *entity.RowMap, attribute string) bool {
	_, exists := (*row)[attribute]
	return !entity.IsQuoted(attribute) && !entity.IsNumeric(attribute) && !exists
}

func (c *Comparison) checkIfAttributeAndCompare(
	row *entity.RowMap,
	valuesLeft, valuesRight []string,
	result *entity.Relation,
	params Comparison,
) {
	if len(valuesLeft) != 0 && len(valuesRight) != 0 {
		for _, valLeft := range valuesLeft {
			for _, valRight := range valuesRight {
				params.valLeft = valLeft
				params.valRight = valRight
				c.valueComparator(params, row, result)
			}
		}
	} else if len(valuesLeft) != 0 {
		for _, valLeft := range valuesLeft {
			params.valLeft = valLeft
			c.valueComparator(params, row, result)
		}
	} else if len(valuesRight) != 0 {
		for _, valRight := range valuesRight {
			params.valRight = valRight
			c.valueComparator(params, row, result)
		}
	} else {
		c.valueComparator(params, row, result)
	}
}

func (c *Comparison) valueComparator(params Comparison, row *entity.RowMap, result *entity.Relation) {
	comparatorTokens := []string{">", "<", ">=", "<=", "=", "!="}
	if slices.Contains(comparatorTokens, params.operator) && c.compare(params) {
		(*result)[row] = struct{}{}
	}
}

func (c *Comparison) compare(params Comparison) bool {
	return c.numericComparator(params) || c.dateComparator(params) || c.stringComparator(params)
}

func (*Comparison) numericComparator(params Comparison) bool {
	oldValNum, err := strconv.ParseFloat(params.valLeft, 10)
	if err != nil {
		return false
	}

	newValNum, err := strconv.ParseFloat(params.valRight, 10)
	if err != nil {
		return false
	}

	switch params.operator {
	case "=":
		return oldValNum == newValNum
	case "!=":
		return oldValNum != newValNum
	case "<=":
		return oldValNum <= newValNum
	case ">=":
		return oldValNum >= newValNum
	case "<":
		return oldValNum < newValNum
	case ">":
		return oldValNum > newValNum
	}

	return false
}

func (*Comparison) dateComparator(params Comparison) bool {
	if entity.IsQuoted(params.valLeft) {
		params.valLeft = removeQuotes(params.valLeft)
	}

	if entity.IsQuoted(params.valRight) {
		params.valRight = removeQuotes(params.valRight)
	}

	oldValDate, err := time.Parse(time.DateTime, params.valLeft)
	if err != nil {
		return false
	}

	newValDate, err := time.Parse(time.DateTime, params.valRight)
	if err != nil {
		return false
	}

	switch params.operator {
	case "=":
		return oldValDate.Equal(newValDate)
	case "!=":
		return !oldValDate.Equal(newValDate)
	case "<=":
		return oldValDate.Before(newValDate) || oldValDate.Equal(newValDate)
	case ">=":
		return oldValDate.After(newValDate) || oldValDate.Equal(newValDate)
	case "<":
		return oldValDate.Before(newValDate)
	case ">":
		return oldValDate.After(newValDate)
	}

	return false
}

func (*Comparison) stringComparator(params Comparison) bool {
	if entity.IsQuoted(params.valLeft) {
		params.valLeft = removeQuotes(params.valLeft)
	}

	if entity.IsQuoted(params.valRight) {
		params.valRight = removeQuotes(params.valRight)
	}

	oldVal := params.valLeft
	newVal := params.valRight
	switch params.operator {
	case "=":
		return oldVal == newVal
	case "!=":
		return oldVal != newVal
	case "<=":
		return oldVal <= newVal
	case ">=":
		return oldVal >= newVal
	case "<":
		return oldVal < newVal
	case ">":
		return oldVal > newVal
	}

	return false
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}
