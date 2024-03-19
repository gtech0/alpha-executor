package main

import (
	"alpha-executor/service"
	"bufio"
	"github.com/kr/pretty"
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
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	program := service.GenerateAST(bufio.NewReader(file))
	pretty.Print(program)
}
