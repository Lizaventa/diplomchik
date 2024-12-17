package repository

import "github.com/jackc/pgx/v4/pgxpool"

type PsqlConnection struct {
	pool *pgxpool.Pool
}

func NewPsqlConnection(pool *pgxpool.Pool) *PsqlConnection {
	return &PsqlConnection{pool: pool}
}
