package operation

import (
	"alpha-executor/entity"
)

type Union struct {
}

func (*Union) Execute(relation1, relation2 *entity.Relation, position entity.Position) (*entity.Relation, error) {
	if err := relation1.EqualArity(relation2, position); err != nil {
		return nil, err
	}

	relation := make(entity.Relation)
	for row1 := range *relation1 {
		relation[row1] = struct{}{}
	}

	for row2 := range *relation2 {
		relation[row2] = struct{}{}
	}
	return &relation, nil
}
