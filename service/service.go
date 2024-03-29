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
	e.testingRepository.AddRelations(receiver.Relations)

	reader := strings.NewReader(receiver.Query)
	program := operation.GenerateAST(bufio.NewReader(reader))
	pretty.Print(program)
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
		//GetResults: e.testingRepository.GetRangedRelation(),
	}, nil
}
