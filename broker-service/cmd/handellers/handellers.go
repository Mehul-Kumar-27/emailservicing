package handellers

import (
	"net/http"
)

type ServerModel struct {
}

func NewServerModel() *ServerModel {
	return &ServerModel{}
}

func (app *ServerModel) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Welcome to the Broker API",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}
