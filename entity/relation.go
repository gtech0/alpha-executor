package entity

import (
	"encoding/json"
	"fmt"
)

type Relation map[*RowMap]struct{}
type Relations map[string]*Relation

func (r *Relation) MarshalJSON() ([]byte, error) {
	set := r.keysToSlice()
	index := 0
	for rowMap := range *r {
		set[index] = rowMap
		index++
	}

	marshal, err := json.Marshal(set)
	return marshal, err
}

func (r *Relation) keysToSlice() []*RowMap {
	set := make([]*RowMap, len(*r))
	index := 0
	for rowMap := range *r {
		set[index] = rowMap
		index++
	}
	return set
}

func (r *Relation) UnmarshalJSON(data []byte) error {
	var relation []*RowMap
	err := json.Unmarshal(data, &relation)
	if err != nil {
		return err
	}

	relationConverted := make(Relation)
	for _, value := range relation {
		relationConverted[value] = struct{}{}
	}

	*r = relationConverted

	return err
}

func (r *Relation) EqualArity(r2 *Relation, position Position) error {
	for row1 := range *r {
		for row2 := range *r2 {
			if !row1.keysEqual(row2) {
				return &CustomError{
					ErrorType: ResponseTypes["RT"],
					Message:   fmt.Sprintf("Incorrect arity at %d:%d", position.Line, position.Column),
					Position:  position,
				}
			}
		}
	}

	return nil
}

func (r *Relation) RelationEqual(r2 *Relation) bool {
	for row1 := range *r {
		equal := false
		for row2 := range *r2 {
			if row1.RowsEqual(row2) {
				equal = true
			}
		}

		if !equal {
			return false
		}
	}
	return true
}

func (r *Relations) RelationsEqual(r2 *Relations) bool {
	for _, relation1 := range *r {
		equal := false
		for _, relation2 := range *r2 {
			if relation1.RelationEqual(relation2) {
				equal = true
			}
		}

		if !equal {
			return false
		}
	}
	return true
}
