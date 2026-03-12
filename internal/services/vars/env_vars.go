package vars

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (v *Vars) GetEnvVar(ctx context.Context, id string) (*models.Var, error) {
	// TODO: Auth
	res, err := v.envVars.EnvVar(ctx, id)
	return res, err
}

func (v *Vars) ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error) {
	// TODO: Auth
	res, err := v.envVars.ListEnvVars(ctx, args)
	return res, err
}

func (v *Vars) CreateEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.Var, error) {
	// TODO: Auth
	res, err := v.envVars.SaveEnvVar(ctx, args)
	if err != nil {
		return nil, err
	}
	return &models.Var{
		Id:    res.Id,
		Key:   args.Key,
		Value: args.Value,
	}, nil
}

func (v *Vars) UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error {
	// TODO: Auth
	err := v.envVars.UpdateEnvVar(ctx, args)
	return err
}

func (v *Vars) DeleteEnvVar(ctx context.Context, id string) error {
	// TODO: Auth
	err := v.envVars.DeleteEnvVar(ctx, id)
	return err
}
