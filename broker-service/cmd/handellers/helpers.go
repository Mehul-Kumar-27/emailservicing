package handellers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *ServerModel) readJson(r *http.Request, w http.ResponseWriter, data interface{}) error {
	maxBytes := 1048576 // 1MB

	body := http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(body)

	err := decoder.Decode(data)

	if err != nil {
		return err
	}

	e := decoder.Decode(&struct{}{})
	if e != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

func (app *ServerModel) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {

	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-type", "application/json")

	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	w.Header().Set("Content-Length", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

func (app *ServerModel) writeJsonError(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusInternalServerError

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse

	payload.Error = true
	payload.Message = err.Error()
	payload.Data = nil

	return app.writeJson(w, statusCode, payload)

}
