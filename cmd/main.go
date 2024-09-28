package main

import (
	"log/slog"
	"net/http"
	"time"

	_ "insurance_hack/docs"
	"insurance_hack/internal/httpapi"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

const baseUrl = "/api"

//	@title			Insurance hack api
//	@version		0.1
//	@contact.email	boldyrev.now@mail.ru
//	@contact.name	API Support
//	@license.name	Apache 2.0
//
//	@host			79.174.82.229:80
//	@BasePath		/api
//	@schemes		http

func main() {
	r := chi.NewRouter()
	r.Use(ContentTypeMiddleware)
	r.Use(cors.AllowAll().Handler)
	r.Route(baseUrl, func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler())
		r.Get("/user/{id}", httpapi.SimpleHandler)
	})

	server := http.Server{
		Handler:     r,
		IdleTimeout: time.Second * 30,
	}

	slog.Info("start server on port 80")
	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to stop server", "err", err)
	}
	slog.Info("stop server")
}
