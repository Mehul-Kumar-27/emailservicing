package handellers

import (
	"log"
	"net/http"
)

type mailMessage struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	From    string   `json:"from"`
}

func (app *MailServer) SendMailRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to send the mail")
	var requestPayload mailMessage

	err := app.readJson(r, w, &requestPayload)
	if err != nil {
		log.Println(err)
		app.writeJsonError(w, err, http.StatusBadRequest)
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Body,
	}
	log.Println(msg)

	err = app.Mail.SendMail(&msg)
	if err != nil {
		app.writeJsonError(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Mail sent successfully",
	}

	app.writeJson(w, http.StatusOK, payload)
}
