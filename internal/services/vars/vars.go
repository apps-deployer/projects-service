package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

type Storage interface {
	Repos() RepoFactory
}

type RepoFactory interface {
	ProjectVars() ProjectVarRepository
	EnvVars() EnvVarRepository
	ResolvedVars() ResolvedVarsRepository
}

type ProjectVarRepository interface {
	ListProjectVars(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error)
	SaveProjectVar(ctx context.Context, args *models.CreateProjectVarParams) (*models.SaveVarResponse, error)
	UpdateProjectVar(ctx context.Context, args *models.UpdateVarParams) error
	DeleteProjectVar(ctx context.Context, id string) error
}

type EnvVarRepository interface {
	ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error)
	SaveEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.SaveVarResponse, error)
	UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error
	DeleteEnvVar(ctx context.Context, id string) error
}

type ResolvedVarsRepository interface {
	ResolvedVars(ctx context.Context, envId string) ([]*models.ResolvedVar, error)
}

type Vars struct {
	log *slog.Logger
	pv  ProjectVarRepository
	ev  EnvVarRepository
	rv  ResolvedVarsRepository
}

func New(log *slog.Logger, storage Storage) *Vars {
	repos := storage.Repos()
	return &Vars{
		log: log,
		pv:  repos.ProjectVars(),
		ev:  repos.EnvVars(),
		rv:  repos.ResolvedVars(),
	}
}
