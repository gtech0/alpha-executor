package service

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/operation"
	"alpha-executor/repository"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"io"
	"os"
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

	e.testingRepository.ClearAll()
	e.testingRepository.AddRelations(receiver.Relations)

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

	output := e.testingRepository.GetGetRelations()
	return model.TestingSender{
		Results: &output,
	}, nil
}

func (e *ExecutorService) TestingCli(data *os.File) error {
	var receiver model.TestingReceiver
	err := json.NewDecoder(data).Decode(&receiver)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(receiver); err != nil {
		return &entity.CustomError{
			ErrorType: entity.ResponseTypes["CF"],
			Message:   "Unexpected error",
		}
	}

	body := io.NopCloser(&buffer)
	result, err := e.Execute(body)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("resources/solutions/output.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(result); err != nil {
		return err
	}
	return nil
}

func (e *ExecutorService) ValidationCommon(validationReceiver model.ValidationReceiver, data *model.Config) error {
	for testNum := 0; testNum < data.TestCount; testNum++ {
		relationsFile, err := os.Open(fmt.Sprintf("%s/%d.in", data.Tests, testNum))
		if err != nil {
			return err
		}

		var relations entity.Relations
		err = json.NewDecoder(relationsFile).Decode(&relations)
		if err != nil {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["CF"],
				Message:   fmt.Sprintf("Test %d weren't found", testNum+1),
			}
		}

		resultFile, err := os.Open(fmt.Sprintf("%s/%d.out", data.Tests, testNum))
		if err != nil {
			return err
		}

		var result entity.Relations
		err = json.NewDecoder(resultFile).Decode(&result)
		if err != nil {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["CF"],
				Message:   fmt.Sprintf("Test %d weren't found", testNum+1),
			}
		}

		var testData bytes.Buffer
		if err = json.NewEncoder(&testData).Encode(model.TestingReceiver{
			Query:     validationReceiver.Query,
			Relations: relations,
		}); err != nil {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["CF"],
				Message:   fmt.Sprintf("Corrupted data for test %d", testNum+1),
			}
		}

		processingResult, err := e.Execute(io.NopCloser(&testData))
		if err != nil {
			return err
		}

		file, err := os.OpenFile(fmt.Sprintf("%s/%d.ans", data.Output, testNum), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err = encoder.Encode(processingResult.Results); err != nil {
			return err
		}

		if !(&result).RelationsEqual(processingResult.Results) {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["WA"],
				Message:   fmt.Sprintf("Test %d has failed", testNum+1),
			}
		}
	}

	return nil
}

func (e *ExecutorService) ValidationServer(body io.ReadCloser) error {
	var validationReceiver model.ValidationReceiver
	if err := json.NewDecoder(body).Decode(&validationReceiver); err != nil {
		return err
	}

	data, err := model.GetConfig()
	if err != nil {
		return err
	}

	return e.ValidationCommon(validationReceiver, data)
}

func (e *ExecutorService) ValidationCli() error {
	data, err := model.GetConfig()
	if err != nil {
		return err
	}

	source, err := os.Open(data.Source)
	if err != nil {
		return err
	}

	var validationReceiver model.ValidationReceiver
	if err := json.NewDecoder(source).Decode(&validationReceiver); err != nil {
		return err
	}

	if err := e.ValidationCommon(validationReceiver, data); err != nil {
		return err
	}

	return nil
}
