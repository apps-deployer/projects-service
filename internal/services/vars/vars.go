package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

type ProjectVarStorage interface {
	ProjectVar(ctx context.Context, id string) (*models.Var, error)
	ListProjectVars(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error)
	SaveProjectVar(ctx context.Context, args *models.CreateProjectVarParams) (*models.SaveVarResponse, error)
	UpdateProjectVar(ctx context.Context, args *models.UpdateVarParams) error
	DeleteProjectVar(ctx context.Context, id string) error
}

type EnvVarStorage interface {
	EnvVar(ctx context.Context, id string) (*models.Var, error)
	ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error)
	SaveEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.SaveVarResponse, error)
	UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error
	DeleteEnvVar(ctx context.Context, id string) error
}

type MergedVarsProvider interface {
	MergedVars(ctx context.Context, envId string) ([]*models.Var, error)
}

type Vars struct {
	log         *slog.Logger
	projectVars ProjectVarStorage
	envVars     EnvVarStorage
	mergedVars  MergedVarsProvider
}

func New(
	log *slog.Logger,
	projectVars ProjectVarStorage,
	envVars EnvVarStorage,
	mergedVars MergedVarsProvider,
) *Vars {
	return &Vars{
		log:         log,
		projectVars: projectVars,
		envVars:     envVars,
		mergedVars:  mergedVars,
	}
}
