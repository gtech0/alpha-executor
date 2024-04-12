package main

import (
	"alpha-executor/controller"
	"alpha-executor/entity"
	"alpha-executor/repository"
	"alpha-executor/router"
	"alpha-executor/service"
	"log"
	"os"
)

func main() {
	file, err := os.Open("resources/test.json")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	testingRepository := repository.NewTestingRepository(make(entity.RowsMap), make(entity.Relations), make(entity.Relations), make(entity.Relation))
	executorService := service.NewExecutorService(testingRepository)
	requestController := controller.NewRequestController(executorService)

	requestRouter := router.NewRouter(requestController)
	requestRouter.Server()
}
