// routes.go
package handellers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app * MailServer)Routes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
		AllowCredentials: true,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/sendMail", app.SendMailRoute)

	return mux
}
