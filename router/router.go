package router

import (
	"alpha-executor/controller"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
)

type Router struct {
	requestController *controller.RequestController
}

func NewRouter(requestController *controller.RequestController) *Router {
	return &Router{
		requestController: requestController,
	}
}

func (r *Router) Server() {
	fmt.Println("server app launched")
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	router.Post("/interpreter/execute", r.requestController.TestingServer)
	router.Post("/interpreter/validate", r.requestController.ValidationServer)

	port := ":8080"
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(fmt.Sprint("can't listen on ", port))
		return
	}
}

func (r *Router) Cli() {
	fmt.Println("cli app launched")

	var err error
	validation := flag.Lookup("validation").Value.String()
	if validation == "true" {
		err = r.requestController.ValidationCli()
	} else {
		var testData *os.File
		testData, err = os.Open("cli/resources/test.json")
		if err != nil {
			log.Fatal("incorrect config path")
		}

		err = r.requestController.TestingCli(testData)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
}
