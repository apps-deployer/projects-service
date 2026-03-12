package vars

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (v *Vars) ListAllVars(ctx context.Context, envId string) ([]*models.Var, error) {
	// TODO: Auth
	res, err := v.mergedVars.MergedVars(ctx, envId)
	return res, err
}
