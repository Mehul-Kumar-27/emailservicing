package main

import (
	"authentication-service/api/handellers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort = "8080"
)

func main() {
	log.Printf("Starting the authentication service at port %s", webPort)

	connection := connectToDB()

	if connection == nil {
		log.Panic("Could not connect to database")
	}

	app := handellers.NewDatabaseModel(connection)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDatase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ping := db.Ping()
	if ping != nil {
		return nil, ping
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	count := 0
	for {
		db, err := openDatase(dsn)
		if err != nil {
			log.Println("Error connecting to database", err)
			log.Println("Retrying in 2 seconds")
			count++
			if count > 10 {
				log.Println("Could not connect to database")
				return nil
			}
			time.Sleep(2 * time.Second)
			continue
		} else {
			log.Println("Connected to database")
			return db
		}

	}
}
