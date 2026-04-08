package envs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
	"github.com/apps-deployer/projects-service/internal/services"
)

type Envs struct {
	log      *slog.Logger
	envs     services.EnvRepository
	projects services.ProjectRepository
}

func New(log *slog.Logger, storage services.Storage) *Envs {
	return &Envs{
		log:      log,
		envs:     storage.Repos().Envs(),
		projects: storage.Repos().Projects(),
	}
}

// checkProjectOwnership fetches the project and verifies the caller owns it.
func (e *Envs) checkProjectOwnership(ctx context.Context, projectID string) error {
	project, err := e.projects.Project(ctx, projectID)
	if err != nil {
		return err
	}
	return auth.CheckOwnership(ctx, project.OwnerId)
}

func (e *Envs) GetByGit(ctx context.Context, args *models.GetEnvByGitParams) (*models.Env, error) {
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
	if err := e.checkProjectOwnership(ctx, env.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	return env, nil
}

func (e *Envs) Get(ctx context.Context, id string) (*models.Env, error) {
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
	if err := e.checkProjectOwnership(ctx, env.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	return env, nil
}

func (e *Envs) List(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error) {
	op := "Envs.List"
	log := e.log.With(
		slog.String("op", op),
		slog.String("projectId", args.ProjectId),
	)
	log.Info("listing envs")
	if err := e.checkProjectOwnership(ctx, args.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	envs, err := e.envs.ListEnvs(ctx, args)
	if err != nil {
		log.Error("failed to list envs", sl.Err(err))
		return nil, err
	}
	return envs, nil
}

func (e *Envs) Create(ctx context.Context, args *models.CreateEnvParams) (*models.Env, error) {
	op := "Envs.Create"
	log := e.log.With(
		slog.String("op", op),
		slog.String("projectId", args.ProjectId),
		slog.String("targetBranch", args.TargetBranch),
	)
	log.Info("creating env")
	if err := e.checkProjectOwnership(ctx, args.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	res, err := e.envs.SaveEnv(ctx, args)
	if err != nil {
		log.Error("failed to create env", sl.Err(err))
		return nil, err
	}
	return models.NewEnvFromSaveResponse(args, res), nil
}

func (e *Envs) Update(ctx context.Context, args *models.UpdateEnvParams) error {
	op := "Envs.Update"
	log := e.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating env")
	env, err := e.envs.Env(ctx, args.Id)
	if err != nil {
		log.Error("failed to get env for ownership check", sl.Err(err))
		return err
	}
	if err := e.checkProjectOwnership(ctx, env.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = e.envs.UpdateEnv(ctx, args)
	if err != nil {
		log.Error("failed to update env", sl.Err(err))
		return err
	}
	return nil
}

func (e *Envs) Delete(ctx context.Context, id string) error {
	op := "Envs.Delete"
	log := e.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting env")
	env, err := e.envs.Env(ctx, id)
	if err != nil {
		log.Error("failed to get env for ownership check", sl.Err(err))
		return err
	}
	if err := e.checkProjectOwnership(ctx, env.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = e.envs.DeleteEnv(ctx, id)
	if err != nil {
		log.Error("failed to delete env", sl.Err(err))
		return err
	}
	return nil
}
