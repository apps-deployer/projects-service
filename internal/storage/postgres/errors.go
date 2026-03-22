package postgres

import (
	"errors"

	"github.com/apps-deployer/projects-service/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// mapError преобразует ошибки PostgreSQL в ошибки хранилища
func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case "23505": // unique_violation
			return storage.ErrAlreadyExists
		case "23503": // foreign_key_violation
			return storage.ErrConflict
		case "23514": // check_violation
			return storage.ErrConflict
		}
	}

	return err
}
