package handellers

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

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

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)

}
