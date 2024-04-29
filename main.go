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

	alphaRepository := repository.NewAlphaRepository(
		make(entity.RowsMap),
		make(entity.Relations),
		make(entity.Relations),
		make(entity.Relations),
		make(entity.Relations),
	)
	alphaService := service.NewAlphaService(alphaRepository)
	alphaController := controller.NewAlphaController(alphaService)

	requestRouter := router.NewRouter(alphaController)
	if isCli {
		requestRouter.Cli()
	} else {
		requestRouter.Server()
	}
}
