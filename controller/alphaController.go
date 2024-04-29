package controller

import (
	"alpha-executor/service"
	"encoding/json"
	"net/http"
	"os"
)

type AlphaController struct {
	executor *service.AlphaService
}

func NewAlphaController(
	executor *service.AlphaService,
) *AlphaController {
	return &AlphaController{
		executor: executor,
	}
}

func (rc *AlphaController) TestingServer(w http.ResponseWriter, r *http.Request) {
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

func (rc *AlphaController) TestingCli(data *os.File) error {
	return rc.executor.TestingCli(data)
}

func (rc *AlphaController) ValidationServer(w http.ResponseWriter, r *http.Request) {
	if err := rc.executor.ValidationServer(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode("OK"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rc *AlphaController) ValidationCli() error {
	return rc.executor.ValidationCli()
}
