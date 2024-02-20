package service

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"encoding/json"
	"io"
)

type ExecutorService struct {
	testingRepository *repository.TestingRepository
}

func NewExecutorService(testingRepository *repository.TestingRepository) *ExecutorService {
	return &ExecutorService{
		testingRepository: testingRepository,
	}
}

func (e *ExecutorService) Testing(body io.ReadCloser) (model.TestingSender, error) {
	var receiver model.TestingReceiver
	if err := json.NewDecoder(body).Decode(&receiver); err != nil {
		return model.TestingSender{}, err
	}

	sender, err := e.preExecutionChecks(receiver)
	if err != nil {
		return sender, err
	}

	e.testingRepository.Clear()
	e.testingRepository.AddRelations(receiver.Relations)

	for index := 0; index < len(receiver.Query); index++ {

	}

	output, err := e.testingRepository.GetResult()
	if err != nil {
		return model.TestingSender{}, err
	}

	return model.TestingSender{
		Result:     output,
		GetResults: e.testingRepository.GetIntermediateRelations(),
	}, nil
}

func (*ExecutorService) preExecutionChecks(receiver model.TestingReceiver) (model.TestingSender, error) {
	receiver.Query.DeleteEmpty()
	if len(receiver.Query) == 0 {
		sender := model.TestingSender{
			Result:     new(entity.Relation),
			GetResults: new(entity.Relations),
		}

		return sender, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "Empty query",
		}
	}

	if err := MatchQuery(receiver.Query); err != nil {
		return model.TestingSender{}, err
	}
	return model.TestingSender{}, nil
}
