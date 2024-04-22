package model

import "alpha-executor/entity"

type (
	TestingReceiver struct {
		Query     string           `json:"query"`
		Relations entity.Relations `json:"relations"`
	}

	TestingSender struct {
		Results *entity.Relations `json:"results"`
	}
)
