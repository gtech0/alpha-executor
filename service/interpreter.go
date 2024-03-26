package service

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/operation"
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

func (i *Interpreter) evaluateProgram(program Program) {
	for _, expression := range program.body {
		switch expression.(type) {
		case GetExpression:
			i.evaluateGet(expression.(GetExpression))
			break
		case HoldExpression:
			i.evaluateHold(expression.(HoldExpression))
			break
		case BinaryExpression:
			if expression.(BinaryExpression).kind != "=" {
				position := expression.(BinaryExpression).left.(IdentifierExpression).position
				log.Fatal(fmt.Sprintf("Incorrect assigment operator at %d:%d", position.Line, position.Column))
			}
			//i.evaluateComparison(expression.(BinaryExpression))
			break
		case RangeExpression:
			i.evaluateRange(expression.(RangeExpression))
			break
		default:
			log.Fatal(fmt.Sprintf("Unknown kind %s", expression.GetKind()))
		}
	}
}

func (i *Interpreter) evaluateGet(expression GetExpression) {
	relation := expression.variable.(IdentifierExpression)
	if relation.kind != model.RELATION.String() {
		log.Fatal(fmt.Sprintf("Expected relation on position %d:%d",
			relation.position.Line, relation.position.Column))
	}

	switch expression.expression.GetKind() {
	case model.EXIST.String():
		break
	case model.FOR_ALL.String():
		break
	case
		model.EQUALS.String(),
		model.NOT_EQUALS.String(),
		model.LESS_THAN_EQUALS.String(),
		model.GREATER_THAN_EQUALS.String(),
		model.LESS_THAN.String(),
		model.GREATER_THAN.String():
		result := i.evaluateComparison(expression.expression.(BinaryExpression))

		attributes := make([]string, 0)
		isRelation := false
		for _, data := range expression.relations {
			switch data.GetKind() {
			case model.RELATION.String():
				isRelation = true
				result, err := i.repository.GetRelation(data.(IdentifierExpression).value)
				if err != nil {
					log.Fatal(err)
				}

				i.repository.AddResult(result)
				break
			case model.ATTRIBUTE.String():
				attr := model.Attribute{}
				attribute, err := attr.ExtractAttribute(data.(IdentifierExpression).value)
				if err != nil {
					log.Fatal(err)
				}

				attributes = append(attributes, attribute.Attribute)
				break
			default:
				log.Fatal(fmt.Sprintf("Unexpected type on position %d:%d",
					relation.position.Line, relation.position.Column))
			}

			if isRelation {
				break
			}
		}

		if isRelation {
			break
		}

		projection := operation.Projection{}
		relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: result}
		result, err := projection.Execute(relationPair, attributes)
		if err != nil {
			log.Fatal(err)
		}

		i.repository.AddResult(result)
	}

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

func (i *Interpreter) evaluateHold(expression HoldExpression) {

}

func (i *Interpreter) evaluateComparison(expression BinaryExpression) *entity.Relation {
	comparison := NewComparison(i.repository)
	result, err := comparison.Compare(expression)
	if err != nil {
		log.Fatal(err)
	}

	return result
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

func (i *Interpreter) evaluateRange(expression RangeExpression) {
	relation, err := i.repository.GetRelation(expression.relation.(IdentifierExpression).value)
	if err != nil {
		log.Fatal(err)
	}

	i.repository.AddRelation(expression.variable.(IdentifierExpression).value, relation)
}

func (i *Interpreter) Evaluate(expression Expression) {
	switch expression.GetKind() {
	case model.PROGRAM.String():
		i.evaluateProgram(expression.(Program))
	default:
		log.Fatal("Program undetected: evaluation failed")
	}
	//pretty.Print(i.repository.Relations)
}
