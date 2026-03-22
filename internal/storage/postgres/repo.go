package postgres

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/services"
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

func (rf *RepoFactory) Projects() services.ProjectRepository {
	return rf.repo
}

func (rf *RepoFactory) Frameworks() services.FrameworkRepository {
	return rf.repo
}

func (rf *RepoFactory) DeployConfigs() services.DeployConfigRepository {
	return rf.repo
}

func (rf *RepoFactory) Envs() services.EnvRepository {
	return rf.repo
}

func (rf *RepoFactory) EnvVars() services.EnvVarRepository {
	return rf.repo
}

func (rf *RepoFactory) ProjectVars() services.ProjectVarRepository {
	return rf.repo
}

func (rf *RepoFactory) ResolvedVars() services.ResolvedVarsRepository {
	return rf.repo
}
