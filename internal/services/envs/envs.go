package envs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

type EnvStorage interface {
	Env(ctx context.Context, id string) (*models.Env, error)
	EnvByGit(ctx context.Context, repoUrl string, branch string) (*models.Env, error)
	ListEnvs(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error)
	SaveEnv(ctx context.Context, args *models.SaveEnvParams) (*models.SaveEnvResponse, error)
	UpdateEnv(ctx context.Context, args *models.UpdateEnvParams) error
	DeleteEnv(ctx context.Context, id string) error
}

type Envs struct {
	log  *slog.Logger
	envs EnvStorage
}

func New(log *slog.Logger, envs EnvStorage) *Envs {
	return &Envs{log: log, envs: envs}
}

func (e *Envs) GetByGit(ctx context.Context, args *models.GetEnvByGitParams) (*models.Env, error) {
	// TODO: Auth
	env, err := e.envs.EnvByGit(ctx, args.RepoUrl, args.TargetBranch)
	return env, err
}

func (e *Envs) Get(ctx context.Context, id string) (*models.Env, error) {
	// TODO: Auth
	env, err := e.envs.Env(ctx, id)
	return env, err
}

func (e *Envs) List(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error) {
	// TODO: Auth
	envs, err := e.envs.ListEnvs(ctx, args)
	return envs, err
}

func (e *Envs) Create(ctx context.Context, args *models.CreateEnvParams) (*models.Env, error) {
	// TODO: Auth
	res, err := e.envs.SaveEnv(ctx, &models.SaveEnvParams{
		Name:         args.Name,
		ProjectId:    args.ProjectId,
		TargetBranch: args.TargetBranch,
		DomainName:   args.DomainName,
	})
	if err != nil {
		return nil, err
	}
	return &models.Env{
		Id:           res.Id,
		Name:         args.Name,
		ProjectId:    args.ProjectId,
		TargetBranch: args.TargetBranch,
		DomainName:   args.DomainName,
	}, nil
}

func (e *Envs) Update(ctx context.Context, args *models.UpdateEnvParams) error {
	// TODO: Auth
	err := e.envs.UpdateEnv(ctx, args)
	return err
}

func (e *Envs) Delete(ctx context.Context, id string) error {
	// TODO: Auth
	err := e.envs.DeleteEnv(ctx, id)
	return err
}
