package handellers

import (
	data "authentication-service/api/data"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type DatabaseModel struct {
	DB     *sql.DB
	Models data.Models
}

func NewDatabaseModel(db *sql.DB) *DatabaseModel {
	return &DatabaseModel{
		DB:     db,
		Models: data.New(db),
	}
}

func (dataBaseModel *DatabaseModel) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ReadJson(r, w, &requestPayload)
	if err != nil {
		WriteJsonError(w, err, http.StatusBadRequest)
		return
	}

	user, err := dataBaseModel.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		WriteJsonError(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		WriteJsonError(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Welcome %s", user.Email),
		Data:    user,
	}

	err = dataBaseModel.createLog("login", fmt.Sprintf("user %s logged in", user.Email))
	if err != nil {
		WriteJsonError(w, err, http.StatusInternalServerError)
		return
	}

	WriteJson(w, http.StatusOK, payload)

}

func (app *DatabaseModel) createLog(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.Marshal(entry)

	logServiceUrl := "http://logger-services:8080/log"
	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return errors.New("error creating log")
	}

	return nil

}
