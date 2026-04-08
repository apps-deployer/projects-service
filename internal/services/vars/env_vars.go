package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

func (v *Vars) checkEnvOwnership(ctx context.Context, envID string) error {
	env, err := v.envs.Env(ctx, envID)
	if err != nil {
		return err
	}
	return v.checkProjectOwnership(ctx, env.ProjectId)
}

func (v *Vars) ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error) {
	op := "Vars.ListEnvVars"
	log := v.log.With(
		slog.String("op", op),
		slog.String("envId", args.EnvId),
	)
	log.Info("listing env vars")
	if err := v.checkEnvOwnership(ctx, args.EnvId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	res, err := v.ev.ListEnvVars(ctx, args)
	if err != nil {
		log.Error("failed to list env vars", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (v *Vars) CreateEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.Var, error) {
	op := "Vars.CreateEnvVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("envId", args.EnvId),
		slog.String("key", args.Key),
	)
	log.Info("creating env var")
	if err := v.checkEnvOwnership(ctx, args.EnvId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	res, err := v.ev.SaveEnvVar(ctx, args)
	if err != nil {
		log.Error("failed to create env var", sl.Err(err))
		return nil, err
	}
	return models.NewVarFromSaveResponse(args.Key, res), nil
}

func (v *Vars) UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error {
	op := "Vars.UpdateEnvVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating env var")
	ownerID, err := v.ev.ProjectOwnerByEnvVarID(ctx, args.Id)
	if err != nil {
		log.Error("failed to get project owner for env var", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, ownerID); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = v.ev.UpdateEnvVar(ctx, args)
	if err != nil {
		log.Error("failed to update env var", sl.Err(err))
		return err
	}
	return nil
}

func (v *Vars) DeleteEnvVar(ctx context.Context, id string) error {
	op := "Vars.DeleteEnvVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting env var")
	ownerID, err := v.ev.ProjectOwnerByEnvVarID(ctx, id)
	if err != nil {
		log.Error("failed to get project owner for env var", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, ownerID); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = v.ev.DeleteEnvVar(ctx, id)
	if err != nil {
		log.Error("failed to delete env var", sl.Err(err))
		return err
	}
	return nil
}
