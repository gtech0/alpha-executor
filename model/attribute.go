package model

import (
	"alpha-executor/entity"
	"fmt"
	"strconv"
	"strings"
)

type Attribute struct {
}

type ComplexAttribute struct {
	Relation  string
	Attribute string
}

func (a *Attribute) ReturnExistentAttribute(relations entity.Relations, attribute string) (string, error) {
	for name, relation := range relations {
		for row := range *relation {
			slicedAttribute, err := a.ExtractAttribute(attribute, entity.Position{})
			if err != nil || (slicedAttribute.Relation != name && slicedAttribute.Relation != "") {
				break
			}

			_, exists := (*row)[slicedAttribute.Attribute]
			if exists {
				return slicedAttribute.Attribute, nil
			}
		}
	}

	return "", &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("attribute %s doesn't exist", attribute),
	}
}

func (*Attribute) ExtractAttribute(attribute string, position entity.Position) (ComplexAttribute, error) {
	pointNum := strings.Count(attribute, ".")
	if pointNum == 0 {
		return ComplexAttribute{
			Relation:  "",
			Attribute: attribute,
		}, nil
	}

	if pointNum > 1 {
		return ComplexAttribute{}, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "only one point in attribute allowed",
			Position:  position,
		}
	}

	if entity.IsQuoted(attribute) {
		return ComplexAttribute{}, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "numeric attribute mustn't is quoted",
			Position:  position,
		}
	}

	if _, err := strconv.ParseFloat(attribute, 10); err == nil {
		return ComplexAttribute{}, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "attribute mustn't be numeric",
			Position:  position,
		}
	}

	sliced := strings.Split(attribute, ".")
	return ComplexAttribute{
		Relation:  sliced[0],
		Attribute: sliced[1],
	}, nil
}
