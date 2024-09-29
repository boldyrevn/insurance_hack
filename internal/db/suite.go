package db

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"insurance_hack/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type TestSuite struct {
	suite.Suite
	db DB
}

func (suite *TestSuite) SetupSuite() {
	ctx := context.Background()

	var (
		pgUser = os.Getenv("POSTGRES_USER")
		pgPass = os.Getenv("POSTGRES_PASSWORD")
	)

	configFile, err := os.Open("../../cmd/config.yaml")
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

	suite.db = New(connPool)
}
