package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

func (v *Vars) ResolveVars(ctx context.Context, envId string) ([]*models.ResolvedVar, error) {
	op := "Vars.ResolveVars"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", envId),
	)
	log.Info("resolving vars")

	if err := v.checkEnvOwnership(ctx, envId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}

	res, err := v.rv.ResolvedVars(ctx, envId)
	if err != nil {
		log.Error("failed to resolve vars", sl.Err(err))
		return nil, err
	}
	return res, nil
}
