package router

import (
	"alpha-executor/controller"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
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

	port := ":8080"
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(fmt.Sprint("can't listen on ", port))
		return
	}
}

//func (r *Router) Cli() {
//	fmt.Println("cli app launched")
//
//	solution := flag.Lookup("solution").Value.String()
//	data, err := os.Open(solution)
//	if err != nil {
//		log.Fatal("incorrect solution path")
//	}
//
//	validation := flag.Lookup("validation").Value.String()
//	if validation == "true" {
//		err = r.requestController.ValidationCli(data)
//	} else {
//		err = r.requestController.TestingCli(data)
//	}
//
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//}
