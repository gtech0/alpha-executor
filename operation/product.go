package operation

import (
	"alpha-executor/entity"
)

type Product struct {
}

func (p *Product) Execute(relation1, relation2 *entity.Relation) *entity.Relation {
	relation := make(entity.Relation)
	for row1 := range *relation1 {
		for row2 := range *relation2 {
			relation[p.mergeRows(row1, row2)] = struct{}{}
		}
	}
	return &relation
}

func (*Product) mergeRows(row1, row2 *entity.RowMap) *entity.RowMap {
	row := make(entity.RowMap)
	for key1, values1 := range *row1 {
		row[key1] = append(row[key1], values1...)
	}

	for key2, values2 := range *row2 {
		row[key2] = append(row[key2], values2...)
	}
	return &row
}
