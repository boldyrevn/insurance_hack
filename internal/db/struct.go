package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) DB {
	return &db{conn: conn}
}
