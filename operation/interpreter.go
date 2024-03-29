package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
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
	var last *entity.Relation

	for _, expression := range program.body {
		var err error
		last, err = i.evaluateExpression(expression)
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

	i.repository.AddResult(last)
	return nil
}

func (i *Interpreter) evaluateExpression(expression Expression) (*entity.Relation, error) {
	switch expression.GetKind() {
	case model.RELATION.String():
		return i.repository.GetRelation(expression.(*IdentifierExpression).value)
	case model.RANGED_RELATION.String():
		return i.repository.GetRangedRelation(expression.(*IdentifierExpression).value)
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
		return nil, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Unknown kind %s", expression.GetKind()),
		}
	}
}

func (i *Interpreter) evaluateGet(expression *GetExpression) (*entity.Relation, error) {
	relation := expression.variable.(*IdentifierExpression)
	if relation.kind != model.RELATION.String() {
		return nil, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Expected relation"),
			Position:  relation.position,
		}
	}

	result, err := i.evaluateExpression(expression.expression)
	if err != nil {
		return nil, err
	}

	attributes := make([]string, 0)
	isRelation := false
	for _, data := range expression.relations {
		switch data.GetKind() {
		case model.RELATION.String():
			isRelation = true
			result, err := i.evaluateExpression(data.(*IdentifierExpression))
			if err != nil {
				err.(*entity.CustomError).Position = data.(*IdentifierExpression).position
				return nil, err
			}

			i.repository.AddResult(result)
			break
		case model.ATTRIBUTE.String():
			attr := model.Attribute{}
			assertedData := data.(*IdentifierExpression)
			attribute, err := attr.ExtractAttribute(assertedData.value, assertedData.position)
			if err != nil {
				err.(*entity.CustomError).Position = data.(*IdentifierExpression).position
				return nil, err
			}

			attributes = append(attributes, attribute.Attribute)
			break
		default:
			return nil, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   "Unexpected type",
				Position:  relation.position,
			}
		}

		if isRelation {
			break
		}
	}

	projection := Projection{}
	relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: result}
	result, err = projection.Execute(relationPair, attributes, relation.position)
	if err != nil {
		return nil, err
	}

	return result, nil
	//join := operation.Join{}
	//if len(attributes) > 0 {
	//	rel1, err := i.repository.GetRelation(attributes[0].Relation)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	relPair1 := entity.Pair[string, *entity.Relation]{Left: attributes[0].Relation, Right: rel1}
	//	attributes = attributes[1:]
	//
	//	for len(attributes) > 0 {
	//		rel2, err := i.repository.GetRelation(attributes[0].Relation)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		relPair2 := entity.Pair[string, *entity.Relation]{Left: attributes[0].Relation, Right: rel2}
	//		attributes = attributes[1:]
	//
	//		commonAttributes := make([]string, 0)
	//		for row1 := range *relPair1.Right {
	//			for row2 := range *relPair2.Right {
	//				exists := make(map[string]struct{})
	//				for key := range *row1 {
	//					exists[key] = struct{}{}
	//				}
	//
	//				for key := range *row2 {
	//					if _, ok := exists[key]; ok {
	//						commonAttributes = append(commonAttributes, key)
	//					}
	//				}
	//				break
	//			}
	//			break
	//		}
	//
	//		relPair1.Right, err = join.Execute(relPair1, relPair2, commonAttributes)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	}
	//	i.repository.AddRangedRelation(relationPair.Left, relPair1.Right)
	//}
}

func (i *Interpreter) evaluateComparison(expression *BinaryExpression) (*entity.Relation, error) {
	comparison := NewComparison(i.repository)
	result, err := comparison.Compare(expression)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *Interpreter) evaluateRange(expression *RangeExpression) (*entity.Relation, error) {
	relation, err := i.repository.GetRelation(expression.relation.(*IdentifierExpression).value)
	if err != nil {
		return nil, err
	}

	i.repository.AddRangedRelation(expression.variable.(*IdentifierExpression).value, relation)
	return relation, nil
}

func (i *Interpreter) evaluateConjunction(expression *BinaryExpression) (*entity.Relation, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return nil, err
	}

	intersect := Intersection{}
	return intersect.Execute(left, right, expression.position)
}

func (i *Interpreter) evaluateDisjunction(expression *BinaryExpression) (*entity.Relation, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return nil, err
	}

	union := Union{}
	return union.Execute(left, right, expression.position)
}

func (i *Interpreter) evaluateExists(expression *BinaryExpression) (*entity.Relation, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return nil, err
	}

	relation := make(entity.Relation)
	for row := range *left {
		singleRowRelation := make(entity.Relation)
		singleRowRelation[row] = struct{}{}
		i.repository.AddRelation(expression.left.(*IdentifierExpression).value, &singleRowRelation)

		right, err := i.evaluateExpression(expression.right)
		if err != nil {
			return nil, err
		}

		for newRow := range *right {
			relation[newRow] = struct{}{}
		}
	}

	return &relation, nil
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
