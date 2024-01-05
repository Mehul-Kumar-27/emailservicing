// main.go
package main

import (
	handellers "broker/cmd/handellers"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	webPort = "8080"
)

func main() {
	logger := log.New(os.Stdout, "API", log.Lshortfile)
	logger.Printf("Starting API server %s", webPort)

	app := handellers.NewServerModel()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // Fix the format specifier
		Handler: app.Routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
