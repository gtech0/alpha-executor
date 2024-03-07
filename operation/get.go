package operation

import "alpha-executor/entity"

type Get struct {
}

func (g *Get) Execute(relation *entity.Relation, function func(a, b any)) bool {

	return false
}
