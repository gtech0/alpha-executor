package controller

import (
	"alpha-executor/service"
	"net/http"
)

type RequestController struct {
	service *service.ExecutorService
}

func NewRequestController(service *service.ExecutorService) *RequestController {
	return &RequestController{
		service: service,
	}
}

func (c *RequestController) Testing(w http.ResponseWriter, r *http.Request) {

}
