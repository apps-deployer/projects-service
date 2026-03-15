package envs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type EnvStorage interface {
	Env(ctx context.Context, id string) (*models.Env, error)
	EnvByGit(ctx context.Context, repoUrl string, branch string) (*models.Env, error)
	ListEnvs(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error)
	SaveEnv(ctx context.Context, args *models.CreateEnvParams) (*models.SaveEnvResponse, error)
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
	op := "Envs.GetByGit"
	log := e.log.With(
		slog.String("op", op),
		slog.String("repoUrl", args.RepoUrl),
		slog.String("targetBranch", args.TargetBranch),
	)
	log.Info("getting env by git")
	env, err := e.envs.EnvByGit(ctx, args.RepoUrl, args.TargetBranch)
	if err != nil {
		log.Error("failed to get env by git", sl.Err(err))
		return nil, err
	}
	return env, nil
}

func (e *Envs) Get(ctx context.Context, id string) (*models.Env, error) {
	// TODO: Auth
	op := "Envs.Get"
	log := e.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("getting env")
	env, err := e.envs.Env(ctx, id)
	if err != nil {
		log.Error("failed to get env", sl.Err(err))
		return nil, err
	}
	return env, nil
}

func (e *Envs) List(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error) {
	// TODO: Auth
	op := "Envs.List"
	log := e.log.With(
		slog.String("op", op),
		slog.String("projectId", args.ProjectId),
	)
	log.Info("listing envs")
	envs, err := e.envs.ListEnvs(ctx, args)
	if err != nil {
		log.Error("failed to list envs", sl.Err(err))
		return nil, err
	}
	return envs, nil
}

func (e *Envs) Create(ctx context.Context, args *models.CreateEnvParams) (*models.Env, error) {
	// TODO: Auth
	op := "Envs.Create"
	log := e.log.With(
		slog.String("op", op),
		slog.String("projectId", args.ProjectId),
		slog.String("targetBranch", args.TargetBranch),
	)
	log.Info("creating env")
	res, err := e.envs.SaveEnv(ctx, &models.CreateEnvParams{
		Name:         args.Name,
		ProjectId:    args.ProjectId,
		TargetBranch: args.TargetBranch,
		DomainName:   args.DomainName,
	})
	if err != nil {
		log.Error("failed to create env", sl.Err(err))
		return nil, err
	}
	return &models.Env{
		Id:           res.Id,
		Name:         args.Name,
		ProjectId:    args.ProjectId,
		TargetBranch: args.TargetBranch,
		DomainName:   args.DomainName,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (e *Envs) Update(ctx context.Context, args *models.UpdateEnvParams) error {
	// TODO: Auth
	op := "Envs.Update"
	log := e.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating env")
	err := e.envs.UpdateEnv(ctx, args)
	if err != nil {
		log.Error("failed to update env", sl.Err(err))
		return err
	}
	return nil
}

func (e *Envs) Delete(ctx context.Context, id string) error {
	// TODO: Auth
	op := "Envs.Delete"
	log := e.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting env")
	err := e.envs.DeleteEnv(ctx, id)
	if err != nil {
		log.Error("failed to delete env", sl.Err(err))
		return err
	}
	return nil
}
