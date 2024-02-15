package main

import (
	"alpha-executor/controller"
	"alpha-executor/entity"
	"alpha-executor/repository"
	"alpha-executor/router"
	"alpha-executor/service"
)

func main() {
	testingRepository := repository.NewTestingRepository(
		make(entity.Relations),
		make(entity.Relations),
	)
	executorService := service.NewExecutorService(testingRepository)
	requestController := controller.NewRequestController(executorService)
	requestRouter := router.NewRequestRouter(requestController)

	requestRouter.Start()
}
