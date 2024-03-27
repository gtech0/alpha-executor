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

func (i *Interpreter) evaluateProgram(program *Program) error {
	var last any

	for _, expression := range program.body {
		var err error
		last, err = i.evaluateExpression(expression)
		if err != nil {
			return err
		}
	}
	i.repository.AddResult(last.(*entity.Relation))
	return nil
}

func (i *Interpreter) evaluateExpression(expression Expression) (any, error) {
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
	default:
		return nil, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Unknown kind %s", expression.GetKind()),
		}
	}
}

func (i *Interpreter) evaluateComparison(expression *BinaryExpression) (any, error) {
	comparison := NewComparison(i.repository)
	result, err := comparison.Compare(expression)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *Interpreter) evaluateGet(expression *GetExpression) (any, error) {
	relation := expression.variable.(*IdentifierExpression)
	if relation.kind != model.RELATION.String() {
		return nil, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message: fmt.Sprintf("Expected relation on position %d:%d",
				relation.position.Line, relation.position.Column),
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
			result, err := i.repository.GetRelation(data.(*IdentifierExpression).value)
			if err != nil {
				return nil, err
			}

			i.repository.AddResult(result)
			break
		case model.ATTRIBUTE.String():
			attr := model.Attribute{}
			attribute, err := attr.ExtractAttribute(data.(*IdentifierExpression).value)
			if err != nil {
				log.Fatal(err)
			}

			attributes = append(attributes, attribute.Attribute)
			break
		default:
			return nil, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message: fmt.Sprintf("Unexpected type on position %d:%d",
					relation.position.Line, relation.position.Column),
			}
		}

		if isRelation {
			break
		}
	}

	projection := Projection{}
	relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: result.(*entity.Relation)}
	result, err = projection.Execute(relationPair, attributes)
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
	//	i.repository.AddIntermediateRelation(relationPair.Left, relPair1.Right)
	//}
}

//func (i *Interpreter) evaluateConjunction(expression *BinaryExpression) (*entity.Relation, error) {
//
//}

//func (i *Interpreter) evaluateDisjunction(expression *BinaryExpression) (*entity.Relation, error) {
//
//}

func (i *Interpreter) evaluateHold(expression *HoldExpression) {

}

//func (i *Interpreter) evaluateForAll(expression BinaryExpression) {
//	comparison := NewComparison(i.repository)
//	err := comparison.Compare(relation, expression)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func (i *Interpreter) evaluateExists(expression BinaryExpression) {
//	comparison := NewComparison(i.repository)
//	err := comparison.Compare(relation, expression)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func (i *Interpreter) evaluateRange(expression *RangeExpression) (any, error) {
	relation, err := i.repository.GetRelation(expression.relation.(*IdentifierExpression).value)
	if err != nil {
		return nil, err
	}

	i.repository.AddRelation(expression.variable.(*IdentifierExpression).value, relation)
	return relation, nil
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
		}
	}
	//pretty.Print(i.repository.Relations)
	return nil
}
