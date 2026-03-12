package vars

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
)

func (v *Vars) GetProjectVar(ctx context.Context, id string) (*models.Var, error) {
	// TODO: Auth
	res, err := v.projectVars.ProjectVar(ctx, id)
	return res, err
}

func (v *Vars) ListProjectVars(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error) {
	// TODO: Auth
	res, err := v.projectVars.ListProjectVars(ctx, args)
	return res, err
}

func (v *Vars) CreateProjectVar(ctx context.Context, args *models.CreateProjectVarParams) (*models.Var, error) {
	// TODO: Auth
	res, err := v.projectVars.SaveProjectVar(ctx, args)
	if err != nil {
		return nil, err
	}
	return &models.Var{
		Id:    res.Id,
		Key:   args.Key,
		Value: args.Value,
	}, nil
}

func (v *Vars) UpdateProjectVar(ctx context.Context, args *models.UpdateVarParams) error {
	// TODO: Auth
	err := v.projectVars.UpdateProjectVar(ctx, args)
	return err
}

func (v *Vars) DeleteProjectVar(ctx context.Context, id string) error {
	// TODO: Auth
	err := v.projectVars.DeleteProjectVar(ctx, id)
	return err
}
