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

//func main() {
//	testingRepository := repository.NewTestingRepository(
//		make(entity.relations),
//		make(entity.relations),
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

	testingRepository := repository.NewTestingRepository(make(entity.Relations), make(entity.Relations), make(entity.Relation))
	executorService := service.NewExecutorService(testingRepository)
	requestController := controller.NewRequestController(executorService)

	requestRouter := router.NewRouter(requestController)
	requestRouter.Server()

	//var receiver model.TestingReceiver
	//err = json.NewDecoder(file).Decode(&receiver)
	//if err != nil {
	//	panic(err)
	//}

	//reader := strings.NewReader(receiver.Query)
	//program := service.GenerateAST(bufio.NewReader(reader))
	//pretty.Print(program)
	//
	//testingRepository := repository.NewTestingRepository(receiver.relations, make(entity.relations), make(entity.Relation))
	//interpreter := service.NewInterpreter(testingRepository)
	//interpreter.Evaluate(program)
}
