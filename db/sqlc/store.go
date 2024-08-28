package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Provides all function to execute db queries and transactions
type Store interface {
	Querier
}

type SQLStore struct {
	*Queries // satisfies the Store interface
	pool     *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		pool:    pool,
		Queries: New(pool),
	}
}
