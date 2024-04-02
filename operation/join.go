package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
)

type Join struct {
}

func (*Join) Execute(relation1, relation2 entity.Pair[string, *entity.Relation], attributes []string) (*entity.Relation, error) {
	joined := make(entity.Relation)
	times := Product{}
	product := times.Execute(relation1.Right, relation2.Right)

	var relations entity.Relations = map[string]*entity.Relation{
		relation1.Left: relation1.Right,
		relation2.Left: relation2.Right,
	}

	attr := model.Attribute{}
	for row := range *product {
		joinedAttributes := 0
		for _, attribute := range attributes {
			slicedAttribute, err := attr.ReturnExistentAttribute(relations, attribute)
			if err != nil {
				return nil, &entity.CustomError{
					ErrorType: entity.ResponseTypes["CE"],
					Message:   err.Error(),
				}
			}

			values, exist := (*row)[slicedAttribute]
			if !exist {
				continue
			}

			if len(values)%2 == 0 {
				part1 := values[:len(values)/2]
				part2 := values[len(values)/2:]

				if entity.OrderedSlicesEqual(part1, part2) {
					joinedAttributes++
					(*row)[slicedAttribute] = part1
				}
			}
		}

		if joinedAttributes == len(attributes) {
			joined[row] = struct{}{}
		}
	}
	return &joined, nil
}
