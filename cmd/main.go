package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

type GetUserResponse struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func SimpleHandler(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	newUser := GetUserResponse{
		Id:   id,
		Name: "Someone",
		Age:  33,
	}

	body, err := json.Marshal(newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	_, err = writer.Write(body)
	if err != nil {
		slog.Error("got write error", "error", err)
	}
}

func OriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		fmt.Println("Origin:", origin)
		next.ServeHTTP(writer, request)
	})
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

const baseUrl = "/api"

func main() {
	r := chi.NewRouter()
	r.Use(OriginMiddleware)
	r.Use(ContentTypeMiddleware)
	r.Use(cors.AllowAll().Handler)
	r.Route(baseUrl, func(r chi.Router) {
		r.Get("/user/{id}", SimpleHandler)
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
