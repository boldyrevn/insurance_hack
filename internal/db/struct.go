package db

import "github.com/jackc/pgx"

type db struct {
	conn *pgx.ConnPool
}

func New(conn *pgx.ConnPool) DB {
	return &db{conn: conn}
}
