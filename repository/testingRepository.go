package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type TestingRepository struct {
	Relations             entity.Relations
	intermediateRelations entity.Relations
	result                entity.Relation
}

func NewTestingRepository(
	relations entity.Relations,
	intermediateRelations entity.Relations,
	result entity.Relation,
) *TestingRepository {
	return &TestingRepository{
		Relations:             relations,
		intermediateRelations: intermediateRelations,
		result:                result,
	}
}

func (t *TestingRepository) AddRelation(name string, relation *entity.Relation) {
	t.Relations[name] = relation
}

func (t *TestingRepository) AddIntermediateRelation(name string, relation *entity.Relation) {
	t.intermediateRelations[name] = relation
}

func (t *TestingRepository) AddRelations(relations entity.Relations) {
	for name, relation := range relations {
		t.Relations[name] = relation
	}
}

func (t *TestingRepository) AddResult(rel *entity.Relation) {
	t.result = *rel
}

func (t *TestingRepository) GetRelation(name string) (*entity.Relation, error) {
	result := t.Relations[name]
	if result != nil {
		return result, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprint("relation ", name, " is null"),
	}
}

func (t *TestingRepository) GetIntermediateRelations() *entity.Relations {
	return &t.intermediateRelations
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
	clear(t.Relations)
	clear(t.intermediateRelations)
	clear(t.result)
}
