package service

import (
	"alpha-executor/model"
	"alpha-executor/repository"
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type ExecutorService struct {
	testingRepository *repository.TestingRepository
}

func NewExecutorService(testingRepository *repository.TestingRepository) *ExecutorService {
	return &ExecutorService{
		testingRepository: testingRepository,
	}
}

func (e *ExecutorService) Execute(body io.ReadCloser) (model.TestingSender, error) {
	var receiver model.TestingReceiver
	if err := json.NewDecoder(body).Decode(&receiver); err != nil {
		return model.TestingSender{}, err
	}

	e.testingRepository.Clear()
	e.testingRepository.AddRelations(receiver.Relations)

	reader := strings.NewReader(receiver.Query)
	program := GenerateAST(bufio.NewReader(reader))
	//pretty.Print(program)
	interpreter := NewInterpreter(e.testingRepository)
	interpreter.Evaluate(program)

	output, err := e.testingRepository.GetResult()
	if err != nil {
		return model.TestingSender{}, err
	}

	return model.TestingSender{
		Result:     output,
		GetResults: e.testingRepository.GetIntermediateRelations(),
	}, nil
}
