package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(dbUrl string) (*Storage, error) {
	const op = "storage.postgres.New"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) Repositories() *RepoFactory {
	return newRepoFactory(s.pool)
}

func (s *Storage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.pool.BeginTx(ctx, pgx.TxOptions{})
}

func WithinTx[T any](
	ctx context.Context,
	s *Storage,
	fn func(*RepoFactory) (*T, error),
) (*T, error) {
	const op = "storage.postgres.Storage.WithinTx"
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	repos := newRepoFactory(tx)
	res, err := fn(repos)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}
