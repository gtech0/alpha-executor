package repository

import (
	"alpha-executor/entity"
	"fmt"
)

type TestingRepository struct {
	relations             entity.Relations
	intermediateRelations entity.Relations
}

func NewTestingRepository(relations entity.Relations, intermediateRelations entity.Relations) TestingRepository {
	return TestingRepository{
		relations:             relations,
		intermediateRelations: intermediateRelations,
	}
}

func (t *TestingRepository) AddRelation(name string, relation *entity.Relation) {
	t.relations[name] = relation
}

func (t *TestingRepository) AddIntermediateRelation(name string, relation *entity.Relation) {
	t.intermediateRelations[name] = relation
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
	return result, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   fmt.Sprint("relation ", name, " is null"),
	}
}

func (t *TestingRepository) GetIntermediateRelations() *entity.Relations {
	return &t.intermediateRelations
}

func (t *TestingRepository) GetResult() (*entity.Relation, error) {
	result := t.relations[""]

	if result != nil {
		return result, nil
	}
	return result, &entity.CustomError{
		ErrorType: entity.ResponseTypes["CE"],
		Message:   "result is null",
	}
}

func (t *TestingRepository) Clear() {
	clear(t.relations)
	clear(t.intermediateRelations)
}
