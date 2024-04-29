package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type AlphaRepository struct {
	rows                entity.RowsMap
	relations           entity.Relations
	calculatedRelations entity.Relations
	heldRelations       entity.Relations
	getRelations        entity.Relations
}

func NewAlphaRepository(
	rows entity.RowsMap,
	relations entity.Relations,
	calculatedRelations entity.Relations,
	heldRelations entity.Relations,
	getRelations entity.Relations,
) *AlphaRepository {
	return &AlphaRepository{
		rows:                rows,
		relations:           relations,
		calculatedRelations: calculatedRelations,
		heldRelations:       heldRelations,
		getRelations:        getRelations,
	}
}

func (t *AlphaRepository) AddRow(name string, row *entity.RowMap) {
	t.rows[name] = row
}

func (t *AlphaRepository) GetRow(name string) (*entity.RowMap, error) {
	result := t.rows[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("row %s is null", name),
	}
}

func (t *AlphaRepository) AddRelation(name string, relation *entity.Relation) {
	t.relations[name] = relation
}

func (t *AlphaRepository) AddRelations(relations entity.Relations) {
	for name, relation := range relations {
		t.relations[name] = relation
	}
}

func (t *AlphaRepository) GetRelation(name string) (*entity.Relation, error) {
	result := t.relations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("relation %s is null", name),
	}
}

func (t *AlphaRepository) GetAllRelations() entity.Relations {
	return t.relations
}

func (t *AlphaRepository) AddCalculatedRelations(relations entity.Relations) {
	for name, relation := range relations {
		t.calculatedRelations[name] = relation
	}
}

func (t *AlphaRepository) GetCalculatedRelation(name string) (*entity.Relation, error) {
	result := t.calculatedRelations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("relation %s is null", name),
	}
}

func (t *AlphaRepository) AddHeldRelation(name string, relation *entity.Relation) {
	t.heldRelations[name] = relation
}

func (t *AlphaRepository) GetHeldRelation(name string) (*entity.Relation, error) {
	result := t.heldRelations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("relation %s is null", name),
	}
}

func (t *AlphaRepository) AddGetRelation(name string, relation *entity.Relation) {
	t.getRelations[name] = relation
}

func (t *AlphaRepository) GetGetRelations() entity.Relations {
	return t.getRelations
}

func (t *AlphaRepository) ReleaseHeldRelation(name string) {
	delete(t.heldRelations, name)
}

func (t *AlphaRepository) ClearAll() {
	clear(t.rows)
	clear(t.relations)
	clear(t.calculatedRelations)
	clear(t.heldRelations)
	clear(t.getRelations)
}
