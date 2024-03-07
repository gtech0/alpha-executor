package main

import (
	"alpha-executor/service"
	"log"
	"os"
)

//func main() {
//	testingRepository := repository.NewTestingRepository(
//		make(entity.Relations),
//		make(entity.Relations),
//	)
//	executorService := service.NewExecutorService(testingRepository)
//	requestController := controller.NewRequestController(executorService)
//	requestRouter := router.NewRequestRouter(requestController)
//
//	requestRouter.Start()
//}

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	service.GenerateAST(file)
}
