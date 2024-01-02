package handellers

import (
	"fmt"
	"logger-service/cmd/data"
	"net/http"
)

type JSONResponse struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *LoggerService) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONResponse

	_ = app.readJson(r, w, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := event.Create()

	if err != nil {
		fmt.Printf("Error while creating log entry: %s", err)
		app.writeJsonError(w, err)
		return
	}

	response := jsonResponse{
		Message: "Log entry created successfully",
		Error:   false,
	}

	app.writeJson(w, http.StatusCreated, response)
}
