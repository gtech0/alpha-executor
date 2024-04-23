package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type TestingRepository struct {
	rows                entity.RowsMap
	relations           entity.Relations
	calculatedRelations entity.Relations
	heldRelations       entity.Relations
	getRelations        entity.Relations
}

func NewTestingRepository(
	rows entity.RowsMap,
	relations entity.Relations,
	calculatedRelations entity.Relations,
	heldRelations entity.Relations,
	getRelations entity.Relations,
) *TestingRepository {
	return &TestingRepository{
		rows:                rows,
		relations:           relations,
		calculatedRelations: calculatedRelations,
		heldRelations:       heldRelations,
		getRelations:        getRelations,
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

func (t *TestingRepository) AddHeldRelation(name string, relation *entity.Relation) {
	t.heldRelations[name] = relation
}

func (t *TestingRepository) GetHeldRelation(name string) (*entity.Relation, error) {
	result := t.heldRelations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("relation %s is null", name),
	}
}

func (t *TestingRepository) AddGetRelation(name string, relation *entity.Relation) {
	t.getRelations[name] = relation
}

func (t *TestingRepository) GetGetRelations() entity.Relations {
	return t.getRelations
}

func (t *TestingRepository) ReleaseHeldRelation(name string) {
	delete(t.heldRelations, name)
}

func (t *TestingRepository) ClearAll() {
	clear(t.rows)
	clear(t.relations)
	clear(t.calculatedRelations)
	clear(t.heldRelations)
	clear(t.getRelations)
}
