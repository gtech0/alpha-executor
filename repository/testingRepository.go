package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type TestingRepository struct {
	rows          entity.RowsMap
	freeRelations entity.Relations
	result        entity.Relation
}

func NewTestingRepository(
	rows entity.RowsMap,
	freeRelations entity.Relations,
	result entity.Relation,
) *TestingRepository {
	return &TestingRepository{
		rows:          rows,
		freeRelations: freeRelations,
		result:        result,
	}
}

func (t *TestingRepository) AddRow(name string, row *entity.RowMap) {
	t.rows[name] = row
}

func (t *TestingRepository) AddRows(rows entity.RowsMap) {
	for name, row := range rows {
		t.rows[name] = row
	}
}

func (t *TestingRepository) GetRow(name string) (*entity.RowMap, error) {
	result := t.rows[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("row %s is null", name),
	}
}

func (t *TestingRepository) AddFreeRelation(name string, relation *entity.Relation) {
	t.freeRelations[name] = relation
}

func (t *TestingRepository) AddFreeRelations(relations entity.Relations) {
	for name, relation := range relations {
		t.freeRelations[name] = relation
	}
}

func (t *TestingRepository) GetFreeRelation(name string) (*entity.Relation, error) {
	result := t.freeRelations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("relation %s is null", name),
	}
}

func (t *TestingRepository) AddResult(rel *entity.Relation) {
	t.result = *rel
}

func (t *TestingRepository) GetResult() (*entity.Relation, error) {
	result := t.result
	if result != nil {
		return &result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   "result is null",
	}
}

func (t *TestingRepository) Clear() {
	clear(t.rows)
	clear(t.freeRelations)
	clear(t.result)
}
