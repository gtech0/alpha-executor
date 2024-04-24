package main

import (
	"alpha-executor/controller"
	"alpha-executor/entity"
	"alpha-executor/repository"
	"alpha-executor/router"
	"alpha-executor/service"
	"flag"
)

func main() {
	var isCli bool
	flag.BoolVar(&isCli, "cli", false, "launch a command line app")
	flag.String("config-path", "", "config file location")
	flag.Bool("validation", false, "executes validation if true, testing if false")
	flag.Parse()

	testingRepository := repository.NewTestingRepository(
		make(entity.RowsMap),
		make(entity.Relations),
		make(entity.Relations),
		make(entity.Relations),
		make(entity.Relations),
	)
	executorService := service.NewExecutorService(testingRepository)
	requestController := controller.NewRequestController(executorService)

	requestRouter := router.NewRouter(requestController)
	requestRouter.Server()
}
