package model

import "alpha-executor/entity"

type (
	TestingReceiver struct {
		Query     entity.Query     `json:"query"`
		Relations entity.Relations `json:"relations"`
	}

	TestingSender struct {
		Result     *entity.Relation  `json:"result"`
		GetResults *entity.Relations `json:"getResults"`
	}
)
