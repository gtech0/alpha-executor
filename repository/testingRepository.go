package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type TestingRepository struct {
	rows                entity.RowsMap
	relations           entity.Relations
	calculatedRelations entity.Relations
	result              entity.Relation
}

func NewTestingRepository(
	rows entity.RowsMap,
	relations entity.Relations,
	calculatedRelations entity.Relations,
	result entity.Relation,
) *TestingRepository {
	return &TestingRepository{
		rows:                rows,
		relations:           relations,
		calculatedRelations: calculatedRelations,
		result:              result,
	}
}

func (t *TestingRepository) AddRow(name string, row *entity.RowMap) {
	t.rows[name] = row
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

func (t *TestingRepository) AddRelation(name string, relation *entity.Relation) {
	t.relations[name] = relation
}

func (t *TestingRepository) AddRelations(relations entity.Relations) {
	for name, relation := range relations {
		t.relations[name] = relation
	}
}

func (t *TestingRepository) GetRelation(name string) (*entity.Relation, error) {
	result := t.relations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("relation %s is null", name),
	}
}

func (t *TestingRepository) GetAllRelations() entity.Relations {
	return t.relations
}

func (t *TestingRepository) AddCalculatedRelation(name string, relation *entity.Relation) {
	t.calculatedRelations[name] = relation
}

func (t *TestingRepository) AddCalculatedRelations(relations entity.Relations) {
	for name, relation := range relations {
		t.calculatedRelations[name] = relation
	}
}

func (t *TestingRepository) GetCalculatedRelation(name string) (*entity.Relation, error) {
	result := t.calculatedRelations[name]
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
	clear(t.relations)
	clear(t.calculatedRelations)
	clear(t.result)
}
