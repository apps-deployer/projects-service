package envs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type Storage interface {
	Envs() EnvRepository
}

type EnvRepository interface {
	Env(ctx context.Context, id string) (*models.Env, error)
	EnvByGit(ctx context.Context, repoUrl string, branch string) (*models.Env, error)
	ListEnvs(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error)
	SaveEnv(ctx context.Context, args *models.CreateEnvParams) (*models.SaveEnvResponse, error)
	UpdateEnv(ctx context.Context, args *models.UpdateEnvParams) error
	DeleteEnv(ctx context.Context, id string) error
}

type Envs struct {
	log     *slog.Logger
	storage Storage
}

func New(log *slog.Logger, storage Storage) *Envs {
	return &Envs{
		log:     log,
		storage: storage,
	}
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
	env, err := e.storage.Envs().EnvByGit(ctx, args.RepoUrl, args.TargetBranch)
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
	env, err := e.storage.Envs().Env(ctx, id)
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
	envs, err := e.storage.Envs().ListEnvs(ctx, args)
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
	res, err := e.storage.Envs().SaveEnv(ctx, args)
	if err != nil {
		log.Error("failed to create env", sl.Err(err))
		return nil, err
	}
	return models.NewEnvFromSaveResponse(args, res), nil
}

func (e *Envs) Update(ctx context.Context, args *models.UpdateEnvParams) error {
	// TODO: Auth
	op := "Envs.Update"
	log := e.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating env")
	err := e.storage.Envs().UpdateEnv(ctx, args)
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
	err := e.storage.Envs().DeleteEnv(ctx, id)
	if err != nil {
		log.Error("failed to delete env", sl.Err(err))
		return err
	}
	return nil
}
