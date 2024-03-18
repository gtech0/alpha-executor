package service

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

//func evaluateProgram(program Program) InterpreterValue {
//	last := InterpreterValue{
//		Type:  0,
//		Value: nil,
//	}
//
//	for i, i := range program.body {
//
//	}
//}
//
//func (i *Interpreter) Evaluate(expression Expression) InterpreterValue {
//
//	switch expression.GetKind() {
//	case PROGRAM:
//		evaluateProgram(expression.(Program))
//	default:
//		log.Fatal("Type undetected: evaluation failed")
//	}
//}
