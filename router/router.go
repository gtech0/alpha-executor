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

type RequestRouter struct {
	controller *controller.RequestController
}

func NewRequestRouter(controller *controller.RequestController) *RequestRouter {
	return &RequestRouter{
		controller: controller,
	}
}

func (r *RequestRouter) Start() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	router.Post("/interpreter/testing", r.controller.Testing)

	port := ":8080"
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(fmt.Sprint("can't listen on ", port))
		return
	}
}
