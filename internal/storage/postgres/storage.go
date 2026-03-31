package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/apps-deployer/projects-service/internal/services"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool          *pgxpool.Pool
	encryptionKey string
}

func New(dbUrl string, encryptionKey string) (*Storage, error) {
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
	return &Storage{pool: pool, encryptionKey: encryptionKey}, nil
}

func (s *Storage) Stop() {
	s.pool.Close()
}

func (s *Storage) Repos() services.RepoFactory {
	return newRepoFactory(s.pool, s.encryptionKey)
}

func (s *Storage) WithinTx(
	ctx context.Context,
	fn func(services.RepoFactory) error,
) error {
	const op = "storage.postgres.Storage.WithinTx"
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	repos := newRepoFactory(tx, s.encryptionKey)
	err = fn(repos)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
