package main

import (
	"fmt"
	"log"
	"mailer-service/cmd/handellers"
	"net/http"
)

const webPort = 8080

func main() {
	app := handellers.MailServer{
		Mail: handellers.CreateMail(),
	}

	log.Printf("Starting the mailserver on port %v \n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("Error to start and listen to the server: %v \n", err)
	}

}
