package model

import "alpha-executor/entity"

type (
	TestingReceiver struct {
		Query     string           `json:"query"`
		Relations entity.Relations `json:"relations"`
	}

	TestingSender struct {
		Result *entity.Relation `json:"result"`
	}
)
