package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type TestingRepository struct {
	relations       entity.Relations
	rangedRelations entity.Relations
	result          entity.Relation
}

func NewTestingRepository(
	relations entity.Relations,
	rangedRelations entity.Relations,
	result entity.Relation,
) *TestingRepository {
	return &TestingRepository{
		relations:       relations,
		rangedRelations: rangedRelations,
		result:          result,
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

func (t *TestingRepository) AddRangedRelation(name string, relation *entity.Relation) {
	t.rangedRelations[name] = relation
}

func (t *TestingRepository) GetRangedRelation(name string) (*entity.Relation, error) {
	result := t.rangedRelations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprintf("ranged relation %s is null", name),
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
	clear(t.relations)
	clear(t.rangedRelations)
	clear(t.result)
}
