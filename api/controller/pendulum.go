package controller

import (
	"net/http"

	pm "github.com/arpit32/conduit/api/model"
	"github.com/arpit32/conduit/api/service"
)

// PendulumController ...
type PendulumController struct {
	BaseController
	PendulumService *service.PendulumService
}

// CreateJob forwards the request to Pendulum Service to create workflow
func (l *PendulumController) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req pm.Trip

	err := l.decodeAndValidate(r, &req)
	if err != nil {
		l.WriteError(r, w, err)
		return
	}

	exec, err := l.PendulumService.CreateJob(r.Context(), req)
	if err != nil {
		l.WriteError(r, w, err)
		return
	}

	l.WriteJSON(r, w, http.StatusOK, exec)
}
