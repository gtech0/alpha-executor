package service

//
//import (
//	"alpha-executor/entity"
//	"alpha-executor/model"
//	"alpha-executor/repository"
//	"fmt"
//	"log"
//)
//
//type Interpreter struct {
//	repository *repository.TestingRepository
//}
//
//func NewInterpreter(repository *repository.TestingRepository) *Interpreter {
//	return &Interpreter{
//		repository: repository,
//	}
//}
//
//func (i *Interpreter) evaluateProgram(program Program) entity.Relation {
//	last := InterpreterValue{
//		Type:  0,
//		Value: nil,
//	}
//
//	for _, expression := range program.body {
//		switch expression.(type) {
//		case GetExpression:
//			i.evaluateGet(expression.(GetExpression))
//			break
//		case BinaryExpression:
//			i.evaluateBinary(expression.(BinaryExpression))
//			break
//		case RangeExpression:
//			i.evaluateRange(expression.(RangeExpression))
//			break
//		default:
//			log.Fatal(fmt.Sprintf("Unknown kind %s", expression.GetKind()))
//		}
//	}
//}
//
//func (i *Interpreter) evaluateGet(expression GetExpression) bool {
//	return false
//}
//
//func (i *Interpreter) evaluateBinary(expression BinaryExpression) bool {
//	left := expression.left.GetKind()
//	right := expression.right.GetKind()
//	if left == model.ATTRIBUTE.String() && right == model.ATTRIBUTE.String() {
//
//	} else if left == model.ATTRIBUTE.String() && right != model.ATTRIBUTE.String() {
//
//	} else if left != model.ATTRIBUTE.String() && right == model.ATTRIBUTE.String() {
//
//	} else {
//
//	}
//	return false
//}
//
//func (i *Interpreter) evaluateBinaryNumeric(expression BinaryExpression) bool {
//	switch expression.kind {
//	case model.EQUALS.String():
//		break
//	case model.NOT_EQUALS.String():
//		break
//	case model.LESS_THAN.String():
//		break
//	case model.LESS_THAN_EQUALS.String():
//		break
//	case model.GREATER_THAN.String():
//		break
//	case model.GREATER_THAN_EQUALS.String():
//		break
//	default:
//		log.Fatal(fmt.Sprintf("Unknown operator %s", expression.kind))
//	}
//	return false
//}
//
//func (i *Interpreter) evaluateRange(expression RangeExpression) bool {
//
//}
//
//func (i *Interpreter) Evaluate(expression Expression) InterpreterValue {
//	switch expression.GetKind() {
//	case model.PROGRAM.String():
//		i.evaluateProgram(expression.(Program))
//	default:
//		log.Fatal("Program undetected: evaluation failed")
//	}
//}
