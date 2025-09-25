package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc_repo "qrcodegen/sqlc/generated"
)

type Repository interface {
	sqlc_repo.Querier
	WithTX(tx pgx.Tx) Repository
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type repository struct {
	*sqlc_repo.Queries
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{
		Queries: sqlc_repo.New(pool),
		pool:    pool,
	}
}

func (r *repository) WithTX(tx pgx.Tx) Repository {
	return &repository{
		Queries: sqlc_repo.New(tx),
		pool:    r.pool,
	}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}
