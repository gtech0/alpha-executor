package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"github.com/kr/pretty"
	"log"
	"slices"
)

type Interpreter struct {
	repository *repository.TestingRepository
}

func NewInterpreter(repository *repository.TestingRepository) *Interpreter {
	return &Interpreter{
		repository: repository,
	}
}

func (i *Interpreter) Evaluate(expression Expression) error {
	switch expression.GetKind() {
	case model.PROGRAM.String():
		err := i.evaluateProgram(expression.(*Program))
		if err != nil {
			return err
		}
	default:
		return &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "Program undetected: evaluation failed",
			Position:  entity.Position{},
		}
	}
	//pretty.Print(i.repository.relations)
	return nil
}

func (i *Interpreter) evaluateProgram(program *Program) error {
	//var last *entity.Relation

	for _, expression := range program.body {
		//var err error
		_, err := i.evaluateExpression(expression)
		if err != nil {
			return err
		}
	}

	if len(program.body) == 0 {
		return &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprint("No query to compile"),
			Position:  entity.Position{},
		}
	}

	//i.repository.AddResult(last)
	return nil
}

func (i *Interpreter) evaluateExpression(expression Expression) (bool, error) {
	switch expression.GetKind() {
	case model.GET.String():
		return i.evaluateGet(expression.(*GetExpression))
	case model.EQUALS.String(),
		model.NOT_EQUALS.String(),
		model.LESS_THAN_EQUALS.String(),
		model.GREATER_THAN_EQUALS.String(),
		model.LESS_THAN.String(),
		model.GREATER_THAN.String():
		return i.evaluateComparison(expression.(*BinaryExpression))
	case model.RANGE.String():
		return i.evaluateRange(expression.(*RangeExpression))
	case model.CONJUNCTION.String():
		return i.evaluateConjunction(expression.(*BinaryExpression))
	case model.DISJUNCTION.String():
		return i.evaluateDisjunction(expression.(*BinaryExpression))
	case model.EXISTS.String():
		return i.evaluateExists(expression.(*BinaryExpression))
	case model.FOR_ALL.String():
		return i.evaluateForAll(expression.(*BinaryExpression))
	default:
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Unknown kind %s", expression.GetKind()),
		}
	}
}

func (i *Interpreter) evaluateExists(expression *BinaryExpression) (bool, error) {
	relationName := expression.left.(*IdentifierExpression).value
	left, err := i.repository.GetRelation(relationName)
	if err != nil {
		return false, err
	}

	for row := range *left {
		i.repository.AddRow(relationName, row)

		right, err := i.evaluateExpression(expression.right)
		if err != nil {
			return false, err
		}

		if right {
			return true, nil
		}
	}

	return false, nil
}

func (i *Interpreter) evaluateForAll(expression *BinaryExpression) (bool, error) {
	relationName := expression.left.(*IdentifierExpression).value
	left, err := i.repository.GetRelation(relationName)
	if err != nil {
		return false, err
	}

	for row := range *left {
		i.repository.AddRow(relationName, row)

		right, err := i.evaluateExpression(expression.right)
		if err != nil {
			return false, err
		}

		if !right {
			return false, nil
		}
	}

	return true, nil
}

/*
	func relationIterator {
		iterate through a list of relation names {
			get relation
			remove relation name from list
			iterate through relation {

			}
		}
	}
*/
func (i *Interpreter) evaluateFreeRelation(
	relations *[]string,
	expression Expression,
	resultRelations *entity.Relations,
) (bool, error) {
	relationName := (*relations)[0]
	*relations = (*relations)[1:]
	relation, err := i.repository.GetRelation(relationName)
	if err != nil {
		return false, err
	}

	result := false
	newRelation := make(entity.Relation)
	for row := range *relation {
		rowCopy := *row
		i.repository.AddRow(relationName, &rowCopy)
		if len(*relations) > 0 {
			result, err = i.evaluateFreeRelation(relations, expression, resultRelations)
			if err != nil {
				return false, err
			}
		} else {
			result, err = i.evaluateExpression(expression)
			if err != nil {
				return false, err
			}

			if result {
				newRelation[row] = struct{}{}
			}
		}
	}

	(*resultRelations)[relationName] = &newRelation
	return result, nil
}

func (i *Interpreter) evaluateGet(expression *GetExpression) (bool, error) {
	relation := expression.variable.(*IdentifierExpression)
	if relation.kind != model.FREE_RELATION.String() {
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Expected relation"),
			Position:  relation.position,
		}
	}

	relations := make([]string, 0)
	attributes := make([]string, 0)
	isRelation := false
	for _, data := range expression.relations {
		switch data.GetKind() {
		case model.FREE_RELATION.String():
			isRelation = true
			relations = append(relations, data.(*IdentifierExpression).value)
			break
		case model.ATTRIBUTE.String():
			attr := model.Attribute{}
			assertedData := data.(*IdentifierExpression)
			attribute, err := attr.ExtractAttribute(assertedData.value, assertedData.position)
			if err != nil {
				err.(*entity.CustomError).Position = data.(*IdentifierExpression).position
				return false, err
			}

			if !slices.Contains(relations, attribute.Relation) {
				relations = append(relations, attribute.Relation)
			}

			if !slices.Contains(relations, attribute.Attribute) {
				attributes = append(attributes, attribute.Attribute)
			}
			break
		default:
			return false, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   "Unexpected type",
				Position:  relation.position,
			}
		}

		if isRelation {
			break
		}
	}

	resultRelations := make(entity.Relations)
	relationsCopy := relations
	result, err := i.evaluateFreeRelation(&relationsCopy, expression.expression, &resultRelations)
	if err != nil {
		err.(*entity.CustomError).Position = expression.position
		return false, err
	}

	i.repository.AddCalculatedRelations(resultRelations)
	if _, err = pretty.Print(resultRelations); err != nil {
		return false, err
	}

	if !isRelation && result {
		rel, err := i.joiningRelations(relations)
		if err != nil {
			return false, err
		}

		projection := Projection{}
		relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: rel}
		projected, err := projection.Execute(relationPair, attributes, relation.position)
		if err != nil {
			return false, err
		}

		i.repository.AddResult(projected)
	} else if isRelation && result {
		final, err := i.repository.GetRelation(relations[0])
		if err != nil {
			return false, err
		}

		i.repository.AddResult(final)
	}

	return true, nil
}

func (i *Interpreter) joiningRelations(relations []string) (*entity.Relation, error) {
	join := Join{}
	for _, rel := range relations {
		rel1, err := i.repository.GetCalculatedRelation(rel)
		if err != nil {
			return nil, err
		}

		rel1Pair := entity.Pair[string, *entity.Relation]{Left: rel, Right: rel1}
		relations = relations[1:]
		for _, rel := range relations {
			rel2, err := i.repository.GetCalculatedRelation(rel)
			if err != nil {
				return nil, err
			}

			rel2Pair := entity.Pair[string, *entity.Relation]{Left: rel, Right: rel2}
			relations = relations[1:]
			commonAttributes := make([]string, 0)
			for row1 := range *rel1Pair.Right {
				for row2 := range *rel2Pair.Right {
					exists := make(map[string]struct{})
					for key := range *row1 {
						exists[key] = struct{}{}
					}

					for key := range *row2 {
						if _, ok := exists[key]; ok {
							commonAttributes = append(commonAttributes, key)
						}
					}
					break
				}
				break
			}

			rel1Pair.Right, err = join.Execute(rel1Pair, rel2Pair, commonAttributes)
			if err != nil {
				log.Fatal(err)
			}
		}
		return rel1Pair.Right, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["RT"],
		Message:   "Relation join error",
		Position:  entity.Position{},
	}
}

func (i *Interpreter) evaluateComparison(expression *BinaryExpression) (bool, error) {
	comparison := NewComparison(i.repository)
	return comparison.Compare(expression)
}

func (i *Interpreter) evaluateRange(expression *RangeExpression) (bool, error) {
	relation, err := i.repository.GetRelation(expression.relation.(*IdentifierExpression).value)
	if err != nil {
		return false, err
	}

	i.repository.AddRelation(expression.variable.(*IdentifierExpression).value, relation)
	return true, nil
}

func (i *Interpreter) evaluateConjunction(expression *BinaryExpression) (bool, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return false, err
	}

	return left && right, nil
}

func (i *Interpreter) evaluateDisjunction(expression *BinaryExpression) (bool, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return false, err
	}

	return left || right, nil
}

//func (i *Interpreter) evaluateHold(expression *HoldExpression) (*entity.Relation, error) {
//
//}
