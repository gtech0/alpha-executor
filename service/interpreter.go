package service

import (
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
				position := expression.(BinaryExpression).left.(IdentifierExpression).token.Position
				log.Fatal(fmt.Sprintf("Incorrect assigment operator at %d:%d", position.Line, position.Column))
			}
			i.evaluateBinary(expression.(BinaryExpression))
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
}

func (i *Interpreter) evaluateHold(expression HoldExpression) {
}

func (i *Interpreter) evaluateBinary(expression BinaryExpression) {
	//comparison := Comparison{}

	//comparison.compareOperands(expression., expression, make(entity.Relation))
	//left := expression.left.GetKind()
	//right := expression.right.GetKind()
	//if left == model.ATTRIBUTE.String() && right == model.ATTRIBUTE.String() {
	//
	//} else if left == model.ATTRIBUTE.String() && right != model.ATTRIBUTE.String() {
	//
	//} else if left != model.ATTRIBUTE.String() && right == model.ATTRIBUTE.String() {
	//
	//} else {
	//
	//}
}

func (i *Interpreter) evaluateRange(expression RangeExpression) {
	relation, err := i.repository.GetRelation(expression.relation.(IdentifierExpression).token.Value)
	if err != nil {
		log.Fatal(err)
	}

	i.repository.AddRelation(expression.variable.(IdentifierExpression).token.Value, relation)
}

func (i *Interpreter) evaluateBinaryNumeric(expression BinaryExpression) {
	switch expression.kind {
	case model.EQUALS.String():
		break
	case model.NOT_EQUALS.String():
		break
	case model.LESS_THAN.String():
		break
	case model.LESS_THAN_EQUALS.String():
		break
	case model.GREATER_THAN.String():
		break
	case model.GREATER_THAN_EQUALS.String():
		break
	default:
		log.Fatal(fmt.Sprintf("Unknown operator %s", expression.kind))
	}
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
