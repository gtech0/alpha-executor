package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
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

func (i *Interpreter) evaluateGet(expression *GetExpression) (bool, error) {
	relation := expression.variable.(*IdentifierExpression)
	if relation.kind != model.RELATION.String() {
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Expected relation"),
			Position:  relation.position,
		}
	}

	result, err := i.evaluateExpression(expression.expression)
	if err != nil {
		err.(*entity.CustomError).Position = expression.position
		return false, err
	}

	relations := make([]string, 0)
	attributes := make([]string, 0)
	isRelation := false
	for _, data := range expression.relations {
		switch data.GetKind() {
		case model.RELATION.String():
			isRelation = true
			if result {
				final, err := i.repository.GetRelation(data.(*IdentifierExpression).value)
				if err != nil {
					return false, err
				}

				i.repository.AddResult(final)
			}
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

	if !isRelation && result {
		rel, err := i.joiningRelations(relations)
		if err != nil {
			return false, err
		}

		projection := Projection{}
		relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: rel}
		result, err := projection.Execute(relationPair, attributes, relation.position)
		if err != nil {
			return false, err
		}

		i.repository.AddResult(result)
	}

	return true, nil
}

func (i *Interpreter) joiningRelations(relations []string) (*entity.Relation, error) {
	join := Join{}
	for _, rel := range relations {
		rel1, err := i.repository.GetRelation(rel)
		if err != nil {
			return nil, err
		}

		rel1Pair := entity.Pair[string, *entity.Relation]{Left: rel, Right: rel1}
		relations = relations[1:]
		for _, rel := range relations {
			rel2, err := i.repository.GetRelation(rel)
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
		return rel1, nil
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

func (i *Interpreter) evaluateExists(expression *BinaryExpression) (bool, error) {
	left, err := i.repository.GetRelation(expression.left.(*IdentifierExpression).value)
	if err != nil {
		return false, err
	}

	for row := range *left {
		relationName := expression.left.(*IdentifierExpression).value
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
	left, err := i.repository.GetRelation(expression.left.(*IdentifierExpression).value)
	if err != nil {
		return false, err
	}

	for row := range *left {
		relationName := expression.left.(*IdentifierExpression).value
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

//func (i *Interpreter) evaluateHold(expression *HoldExpression) (*entity.Relation, error) {
//
//}
