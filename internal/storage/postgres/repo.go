package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repo struct {
	executor QueryExecutor
}

type QueryExecutor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type RepoFactory struct {
	repo *Repo
}

func newRepoFactory(executor QueryExecutor) *RepoFactory {
	repo := &Repo{executor: executor}
	return &RepoFactory{
		repo: repo,
	}
}

func (rf *RepoFactory) Projects() *Repo {
	return rf.repo
}

func (rf *RepoFactory) Frameworks() *Repo {
	return rf.repo
}

func (rf *RepoFactory) DeployConfigs() *Repo {
	return rf.repo
}

func (rf *RepoFactory) Envs() *Repo {
	return rf.repo
}

func (rf *RepoFactory) EnvVars() *Repo {
	return rf.repo
}

func (rf *RepoFactory) ProjectVars() *Repo {
	return rf.repo
}

func (rf *RepoFactory) ResolvedVars() *Repo {
	return rf.repo
}
