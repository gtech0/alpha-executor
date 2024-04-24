package controller

import (
	"alpha-executor/service"
	"encoding/json"
	"net/http"
)

type RequestController struct {
	executor *service.ExecutorService
}

func NewRequestController(
	executor *service.ExecutorService,
) *RequestController {
	return &RequestController{
		executor: executor,
	}
}

func (rc *RequestController) TestingServer(w http.ResponseWriter, r *http.Request) {
	result, err := rc.executor.Execute(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rc *RequestController) ValidationServer(w http.ResponseWriter, r *http.Request) {
	if err := rc.executor.ValidationServer(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode("OK"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
