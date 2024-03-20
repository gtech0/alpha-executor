package main

import (
	"alpha-executor/model"
	"alpha-executor/service"
	"bufio"
	"encoding/json"
	"github.com/kr/pretty"
	"log"
	"os"
	"strings"
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
	file, err := os.Open("resources/test.json")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var receiver model.TestingReceiver
	err = json.NewDecoder(file).Decode(&receiver)
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(receiver.Query)
	program := service.GenerateAST(bufio.NewReader(reader))
	pretty.Print(program)
}
