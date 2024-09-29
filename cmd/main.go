package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"insurance_hack/docs"
	_ "insurance_hack/docs"
	"insurance_hack/internal/db"
	"insurance_hack/internal/httpapi"
	"insurance_hack/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v2"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

const baseUrl = "/api"

//	@title						Insurance hack api
//	@version					0.1
//	@contact.email				boldyrev.now@mail.ru
//	@contact.name				API Support
//	@license.name				Apache 2.0
//
//	@host						79.174.82.229:80
//	@BasePath					/api
//	@schemes					http
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer authorization with JWT

func main() {
	ctx := context.Background()

	var (
		pgUser = os.Getenv("POSTGRES_USER")
		pgPass = os.Getenv("POSTGRES_PASSWORD")
	)

	configFile, err := os.Open("./cmd/config.yaml")
	if err != nil {
		log.Fatalf("failed to open config: %s", err)
	}

	confString, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	var cfg model.ServiceConfig
	if err = yaml.Unmarshal(confString, &cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	docs.SwaggerInfo.Host = cfg.Host

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=%d",
		pgUser,
		pgPass,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.PoolMaxConns,
	)

	connConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("failed to parse connection string: %s", err)
	}

	connPool, err := pgxpool.NewWithConfig(ctx, connConf)
	if err != nil {
		log.Fatalf("failed to connect to postgres")
	}

	dbService := db.New(connPool)

	auth := httpapi.AuthHandler{
		DB: dbService,
	}

	r := chi.NewRouter()
	r.Use(ContentTypeMiddleware)
	r.Use(cors.AllowAll().Handler)
	r.Route(baseUrl, func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler())
		r.Post("/user/create", auth.CreateUser)
		r.Post("/user/token", auth.GetUserToken)
		r.With(auth.AuthMiddleware).Group(func(r chi.Router) {
			r.Get("/user", auth.GetUser)
		})
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
