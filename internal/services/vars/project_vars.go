package vars

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

func (v *Vars) checkProjectOwnership(ctx context.Context, projectID string) error {
	project, err := v.projects.Project(ctx, projectID)
	if err != nil {
		return err
	}
	return auth.CheckOwnership(ctx, project.OwnerId)
}

func (v *Vars) ListProjectVars(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error) {
	op := "Vars.ListProjectVars"
	log := v.log.With(
		slog.String("op", op),
		slog.String("projectId", args.ProjectId),
	)
	log.Info("listing project vars")
	if err := v.checkProjectOwnership(ctx, args.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	res, err := v.pv.ListProjectVars(ctx, args)
	if err != nil {
		log.Error("failed to list project vars", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (v *Vars) CreateProjectVar(ctx context.Context, args *models.CreateProjectVarParams) (*models.Var, error) {
	op := "Vars.CreateProjectVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("projectId", args.ProjectId),
		slog.String("key", args.Key),
	)
	log.Info("creating project var")
	if err := v.checkProjectOwnership(ctx, args.ProjectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	res, err := v.pv.SaveProjectVar(ctx, args)
	if err != nil {
		log.Error("failed to create project var", sl.Err(err))
		return nil, err
	}
	return models.NewVarFromSaveResponse(args.Key, res), nil
}

func (v *Vars) UpdateProjectVar(ctx context.Context, args *models.UpdateVarParams) error {
	op := "Vars.UpdateProjectVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating project var")
	ownerID, err := v.pv.ProjectOwnerByProjectVarID(ctx, args.Id)
	if err != nil {
		log.Error("failed to get project owner for var", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, ownerID); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = v.pv.UpdateProjectVar(ctx, args)
	if err != nil {
		log.Error("failed to update project var", sl.Err(err))
		return err
	}
	return nil
}

func (v *Vars) DeleteProjectVar(ctx context.Context, id string) error {
	op := "Vars.DeleteProjectVar"
	log := v.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting project var")
	ownerID, err := v.pv.ProjectOwnerByProjectVarID(ctx, id)
	if err != nil {
		log.Error("failed to get project owner for var", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, ownerID); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = v.pv.DeleteProjectVar(ctx, id)
	if err != nil {
		log.Error("failed to delete project var", sl.Err(err))
		return err
	}
	return nil
}
