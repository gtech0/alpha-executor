package operation

import (
	"alpha-executor/entity"
)

type Intersection struct {
}

func (*Intersection) Execute(relation1, relation2 *entity.Relation, position entity.Position) (*entity.Relation, error) {
	err := relation1.EqualArity(relation2, position)
	if err != nil {
		return nil, err
	}

	intersection := make(entity.Relation)
	for row1 := range *relation1 {
		for row2 := range *relation2 {
			if row1.RowsEqual(row2) {
				intersection[row1] = struct{}{}
				break
			}
		}
	}
	return &intersection, nil
}
