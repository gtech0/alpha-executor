package service

import (
	"alpha-executor/model"
	"alpha-executor/operation"
	"alpha-executor/repository"
	"bufio"
	"encoding/json"
	"github.com/kr/pretty"
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
	e.testingRepository.AddFreeRelations(receiver.Relations)

	reader := strings.NewReader(receiver.Query)
	program := operation.GenerateAST(bufio.NewReader(reader))
	if _, err := pretty.Print(program); err != nil {
		return model.TestingSender{}, err
	}

	interpreter := operation.NewInterpreter(e.testingRepository)
	err := interpreter.Evaluate(&program)
	if err != nil {
		return model.TestingSender{}, err
	}

	output, err := e.testingRepository.GetResult()
	if err != nil {
		return model.TestingSender{}, err
	}

	return model.TestingSender{
		Result: output,
	}, nil
}
