package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

func (v *Vars) ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error) {
	// TODO: Auth
	op := "Vars.ListEnvVars"
	log := v.log.With(
		slog.String("op", op),
		slog.String("envId", args.EnvId),
	)
	log.Info("listing env vars")
	res, err := v.storage.EnvVars().ListEnvVars(ctx, args)
	if err != nil {
		log.Error("failed to list env vars", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (v *Vars) CreateEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.Var, error) {
	// TODO: Auth
	op := "Vars.CreateEnvVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("envId", args.EnvId),
		slog.String("key", args.Key),
	)
	log.Info("creating env var")
	res, err := v.storage.EnvVars().SaveEnvVar(ctx, args)
	if err != nil {
		log.Error("failed to create env var", sl.Err(err))
		return nil, err
	}
	return models.NewVarFromSaveResponse(args.Key, res), nil
}

func (v *Vars) UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error {
	// TODO: Auth
	op := "Vars.UpdateEnvVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating env var")
	err := v.storage.EnvVars().UpdateEnvVar(ctx, args)
	if err != nil {
		log.Error("failed to update env var", sl.Err(err))
		return err
	}
	return nil
}

func (v *Vars) DeleteEnvVar(ctx context.Context, id string) error {
	// TODO: Auth
	op := "Vars.DeleteEnvVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting env var")
	err := v.storage.EnvVars().DeleteEnvVar(ctx, id)
	if err != nil {
		log.Error("failed to delete env var", sl.Err(err))
		return err
	}
	return nil
}
