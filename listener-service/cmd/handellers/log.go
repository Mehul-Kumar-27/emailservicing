package handellers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func logEvent(payload Payload) error {
	log.Println("Making log")
	jsonData, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "http://logger-services:8080/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {

		return err
	}
	if response.StatusCode != http.StatusCreated {
		return err

	}
	return nil
}
