package service

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"github.com/kr/pretty"
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
			//i.evaluateBinary(expression.(BinaryExpression))
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

	relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: &entity.Relation{}}
	attributes := make([]model.ComplexAttribute, 0)
	for _, data := range expression.relations {
		isRelation := false
		switch data.GetKind() {
		case model.RELATION.String():
			isRelation = true
			result, err := i.repository.GetRelation(data.(IdentifierExpression).value)
			if err != nil {
				log.Fatal(err)
			}

			relationPair.Right = result
			break
		case model.ATTRIBUTE.String():
			attr := model.Attribute{}
			attribute, err := attr.ExtractAttribute(data.(IdentifierExpression).value)
			if err != nil {
				log.Fatal(err)
			}

			attributes = append(attributes, attribute)
			//result, err := i.repository.GetRelation(attribute.Relation)
			//row := &entity.RowMap{}
			//for rowMap := range *result {
			//	values, exists := (*rowMap)[attribute.Attribute]
			//	if !exists {
			//		log.Fatal(fmt.Sprintf("Attribute %s of relation %s doesn't exist",
			//			attribute.Attribute, attribute.Relation))
			//	}
			//	(*row)[attribute.Attribute] = values
			//}
			//(*relationPair.Right)[row] = struct{}{}
			//break
		default:
			log.Fatal(fmt.Sprintf("Unexpected type on position %d:%d",
				relation.position.Line, relation.position.Column))
		}

		if isRelation {
			break
		}
	}

	//join := operation.Join{}
	//if len(attributes) >= 2 {
	//	rel1, err := i.repository.GetRelation(attributes[0].Relation)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	relPair1 := entity.Pair[string, *entity.Relation]{Left: attributes[0].Relation, Right: rel1}
	//
	//	for len(attributes) > 0 {
	//		rel2, err := i.repository.GetRelation(attributes[0].Relation)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		relPair2 := entity.Pair[string, *entity.Relation]{Left: attributes[0].Relation, Right: rel2}
	//		result, err := join.Execute(relPair1, relPair2)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	}
	//}
}

func (i *Interpreter) evaluateHold(expression HoldExpression) {

}

func (i *Interpreter) evaluateBinary(relation entity.Pair[string, *entity.Relation], expression BinaryExpression) {
	comparison := NewComparison(i.repository)
	err := comparison.compareOperands(relation, expression)
	if err != nil {
		log.Fatal(err)
	}
}

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
	pretty.Print(i.repository.Relations)
}
