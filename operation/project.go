package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"fmt"
)

type Projection struct {
}

func (*Projection) Execute(relation entity.Pair[string, *entity.Relation], attributes []string) (*entity.Relation, error) {
	projected := make(entity.Relation)
	attr := model.Attribute{}
	for row := range *relation.Right {
		newRow := make(entity.RowMap)
		for _, attribute := range attributes {
			slicedAttribute, err := attr.ExtractAttribute(attribute)
			if err != nil {
				continue
			}

			values, exists := (*row)[slicedAttribute.Attribute]
			if !exists {
				return nil, &entity.CustomError{
					ErrorType: entity.ResponseTypes["CE"],
					Message:   fmt.Sprintf("attribute %s doesn't exist", attribute),
				}
			}

			data := newRow[slicedAttribute.Attribute]
			data = append(data, values...)
			newRow[slicedAttribute.Attribute] = data
		}

		duplicate := false
		for projectedRow := range projected {
			if projectedRow.RowsEqual(&newRow) {
				duplicate = true
				break
			}
		}

		if !duplicate {
			projected[&newRow] = struct{}{}
		}
	}
	return &projected, nil
}
