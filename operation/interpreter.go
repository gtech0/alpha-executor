package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"log"
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
	//case model.RELATION.String():
	//	return i.repository.GetRelation(expression.(*IdentifierExpression).value)
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

	attributes := make(map[string][]string)
	isRelation := false
	for _, data := range expression.relations {
		switch data.GetKind() {
		case model.RELATION.String():
			isRelation = true
			result, err := i.evaluateExpression(expression.expression)
			if err != nil {
				err.(*entity.CustomError).Position = data.(*IdentifierExpression).position
				return false, err
			}

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

			currentAttributes := attributes[attribute.Attribute]
			attributes[attribute.Attribute] = append(currentAttributes, attribute.Relation)
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

	if !isRelation {
		relations := make([]entity.Pair[string, *entity.Relation], 0)
		for rel, attrs := range attributes {
			projection := Projection{}
			ranged, err := i.repository.GetRelation(rel)
			relationPair := entity.Pair[string, *entity.Relation]{Left: rel, Right: ranged}
			result, err := projection.Execute(relationPair, attrs, relation.position)
			if err != nil {
				return false, err
			}

			relations = append(relations, entity.Pair[string, *entity.Relation]{Left: rel, Right: result})
		}

		join := Join{}
		if len(relations) > 0 {
			rel1 := relations[0]
			relations = relations[1:]
			for len(relations) > 0 {
				rel2 := relations[0]
				commonAttributes := make([]string, 0)
				for row1 := range *rel1.Right {
					for row2 := range *rel2.Right {
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

				var err error
				rel1.Right, err = join.Execute(rel1, rel2, commonAttributes)
				if err != nil {
					log.Fatal(err)
				}
			}
			i.repository.AddResult(rel1.Right)
		}
	}

	return true, nil
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

//func (i *Interpreter) evaluateForAll(expression BinaryExpression) {
//	comparison := NewComparison(i.repository)
//	err := comparison.Compare(relation, expression)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

//func (i *Interpreter) evaluateHold(expression *HoldExpression) (*entity.Relation, error) {
//
//}
