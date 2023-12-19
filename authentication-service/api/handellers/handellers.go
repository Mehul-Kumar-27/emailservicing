package handellers

import (
	data "authentication-service/api/data"
	"database/sql"
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
		Message: fmt.Sprintf("Welcome %s", user.FirstName),
		Data:    user,
	}

	WriteJson(w, http.StatusOK, payload)

}
