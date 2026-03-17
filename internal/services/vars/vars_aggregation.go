package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

func (v *Vars) ResolveVars(ctx context.Context, envId string) ([]*models.Var, error) {
	// TODO: Auth
	op := "Vars.ListAllVars"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", envId),
	)
	log.Info("listing all vars")

	res, err := v.storage.ResolvedVars().ResolvedVars(ctx, envId)
	if err != nil {
		log.Error("failed to list all vars", sl.Err(err))
		return nil, err
	}
	return res, nil
}
