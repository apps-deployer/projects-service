package services

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

type Storage interface {
	Repos() RepoFactory
	WithinTx(
		ctx context.Context,
		fn func(RepoFactory) error,
	) error
}

type RepoFactory interface {
	DeployConfigs() DeployConfigRepository
	Frameworks() FrameworkRepository
	Envs() EnvRepository
	Projects() ProjectRepository
	ProjectVars() ProjectVarRepository
	EnvVars() EnvVarRepository
	ResolvedVars() ResolvedVarsRepository
}

type DeployConfigRepository interface {
	DeployConfig(ctx context.Context, projectId string) (*models.DeployConfig, error)
	UpdateDeployConfig(ctx context.Context, args *models.UpdateDeployConfigParams) error
	SaveDeployConfig(ctx context.Context, args *models.SaveDeployConfigParams) (*models.SaveDeployConfigResponse, error)
	DeleteDeployConfig(ctx context.Context, projectId string) error
}

type FrameworkRepository interface {
	Framework(ctx context.Context, id string) (*models.Framework, error)
	ListFrameworks(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error)
	SaveFramework(ctx context.Context, args *models.CreateFrameworkParams) (*models.SaveFrameworkResponse, error)
	UpdateFramework(ctx context.Context, args *models.UpdateFrameworkParams) error
	DeleteFramework(ctx context.Context, id string) error
}

type EnvRepository interface {
	Env(ctx context.Context, id string) (*models.Env, error)
	EnvByGit(ctx context.Context, repoUrl string, branch string) (*models.Env, error)
	ListEnvs(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error)
	SaveEnv(ctx context.Context, args *models.CreateEnvParams) (*models.SaveEnvResponse, error)
	UpdateEnv(ctx context.Context, args *models.UpdateEnvParams) error
	DeleteEnv(ctx context.Context, id string) error
}

type ProjectRepository interface {
	Project(ctx context.Context, id string) (*models.Project, error)
	ListProjects(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error)
	SaveProject(ctx context.Context, args *models.SaveProjectParams) (*models.SaveProjectResponse, error)
	UpdateProject(ctx context.Context, args *models.UpdateProjectParams) error
	DeleteProject(ctx context.Context, id string) error
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
