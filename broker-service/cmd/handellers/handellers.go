package handellers

import (
	event "broker/cmd/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ServerModel struct {
	RabbitMQconnection *amqp.Connection
}

func NewServerModel(rabbitmq *amqp.Connection) *ServerModel {
	return &ServerModel{
		RabbitMQconnection: rabbitmq,
	}
}

type RequestPayload struct {
	Action string        `json:"action"`
	Auth   AuthPayLoad   `json:"auth,omitempty"`
	Log    LoggerPayload `json:"log,omitempty"`
	Mail   MailerPayload `json:"mail,omitempty"`
}

type AuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoggerPayload struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}
type MailerPayload struct {
	From    string   `json:"from,omitempty"`
	To      []string `json:"to,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Message string   `json:"message,omitempty"`
}

type RPCloggerPayload struct {
	Name string
	Data string
}

func (app *ServerModel) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Welcome to the Broker API",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *ServerModel) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(r, w, &requestPayload)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	log.Println(requestPayload.Action)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		//app.makeLog(w, requestPayload.Log)
		//app.logEventViaRabbitMQ(&w, requestPayload.Log)
		app.logViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.writeJsonError(w, errors.New("unknon method called"))
	}
}

func (app *ServerModel) makeLog(w http.ResponseWriter, a LoggerPayload) {
	log.Println("Making log")
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://logger-services:8080/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	if response.StatusCode != http.StatusCreated {
		app.writeJsonError(w, errors.New("could not create log"))
		return

	}

	jsonPaylod := jsonResponse{
		Error:   false,
		Message: "Log created",
	}

	app.writeJson(w, http.StatusAccepted, jsonPaylod)
}

func (app *ServerModel) authenticate(w http.ResponseWriter, a AuthPayLoad) {
	fmt.Println("Authenticating")
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://auth-services:8080/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		app.writeJsonError(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusOK {
		app.writeJsonError(w, errors.New("could not authenticate"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	if jsonFromService.Error {
		app.writeJsonError(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payloadToSend jsonResponse

	payloadToSend.Error = false
	payloadToSend.Message = "Authenticated"
	payloadToSend.Data = jsonFromService.Data

	app.writeJson(w, http.StatusOK, payloadToSend)
}

func (app *ServerModel) sendMail(w http.ResponseWriter, a MailerPayload) {
	log.Println("Sending mail")
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://mailer-service:8080/sendMail", bytes.NewBuffer(jsonData))
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.writeJsonError(w, errors.New("could not send mail"))
		return
	}

	jsonPaylod := jsonResponse{
		Error:   false,
		Message: "Mail sent",
	}

	app.writeJson(w, http.StatusAccepted, jsonPaylod)
}

func (app *ServerModel) logEventViaRabbitMQ(w *http.ResponseWriter, l LoggerPayload) {
	error := app.pushToQueue(l)
	if error != nil {
		app.writeJsonError(*w, error)
		return
	}

	jsonPaylod := jsonResponse{
		Error:   false,
		Message: "Log created Using RabbitMQ",
	}

	app.writeJson(*w, http.StatusCreated, jsonPaylod)
}

func (app *ServerModel) pushToQueue(l LoggerPayload) error {
	emitter, err := event.NewEmitter(app.RabbitMQconnection)
	if err != nil {
		return err
	}

	json, _ := json.Marshal(l)

	emitter.Push(string(json), "log.INFO")

	return nil
}

func (app *ServerModel) logViaRPC(w http.ResponseWriter, l LoggerPayload) {
	client, err := rpc.Dial("tcp", "logger-services:5001")
	if err != nil {
		log.Println("Error while dialing RPC server:", err)
		app.writeJsonError(w, err)
	}

	var payload RPCloggerPayload
	payload.Name = l.Name
	payload.Data = l.Data

	var response string

	err = client.Call("RPCServer.LogInfo", payload, &response)
	if err != nil {
		log.Println("Error while making RPC call:", err)
		app.writeJsonError(w, err)
	}

	var jsonPaylod jsonResponse

	jsonPaylod.Error = false
	jsonPaylod.Message = response

	app.writeJson(w, http.StatusCreated, jsonPaylod)
}
