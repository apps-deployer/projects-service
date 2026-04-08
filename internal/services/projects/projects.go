package projects

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
	"github.com/apps-deployer/projects-service/internal/services"
)

type Projects struct {
	log      *slog.Logger
	storage  services.Storage
	projects services.ProjectRepository
}

func New(
	log *slog.Logger,
	storage services.Storage,
) *Projects {
	return &Projects{
		log:      log,
		storage:  storage,
		projects: storage.Repos().Projects(),
	}
}

func (p *Projects) Get(ctx context.Context, id string) (*models.Project, error) {
	op := "Projects.Get"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("getting project")
	res, err := p.projects.Project(ctx, id)
	if err != nil {
		log.Error("failed to get project", sl.Err(err))
		return nil, err
	}
	if err := auth.CheckOwnership(ctx, res.OwnerId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (p *Projects) List(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error) {
	op := "Projects.List"
	log := p.log.With(
		slog.String("op", op),
		slog.String("ownerId", args.OwnerId),
	)
	log.Info("listing projects")
	res, err := p.projects.ListProjects(ctx, args)
	if err != nil {
		log.Error("failed to list projects", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (p *Projects) Create(ctx context.Context, args *models.CreateProjectParams) (*models.Project, error) {
	op := "Projects.Create"
	log := p.log.With(
		slog.String("op", op),
		slog.String("name", args.Name),
		slog.String("repoUrl", args.RepoUrl),
	)
	log.Info("creating project")
	project := &models.SaveProjectParams{
		Name:    args.Name,
		RepoUrl: args.RepoUrl,
		OwnerId: args.OwnerId,
	}
	var response *models.SaveProjectResponse
	err := p.storage.WithinTx(ctx, func(tx services.RepoFactory) error {
		res, err := tx.Projects().SaveProject(ctx, project)
		if err != nil {
			return err
		}
		_, err = tx.DeployConfigs().SaveDeployConfig(
			ctx, &models.SaveDeployConfigParams{
				ProjectId:   res.Id,
				FrameworkId: args.FrameworkId,
			},
		)
		if err != nil {
			return err
		}
		response = res
		return nil
	})
	if err != nil {
		log.Error("failed to create project", sl.Err(err))
		return nil, err
	}
	return models.NewProjectFromSaveResponse(project, response), nil
}

func (p *Projects) Update(ctx context.Context, args *models.UpdateProjectParams) error {
	op := "Projects.Update"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating project")
	project, err := p.projects.Project(ctx, args.Id)
	if err != nil {
		log.Error("failed to get project for ownership check", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, project.OwnerId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = p.projects.UpdateProject(ctx, args)
	if err != nil {
		log.Error("failed to update project", sl.Err(err))
		return err
	}
	return nil
}

func (p *Projects) Delete(ctx context.Context, id string) error {
	op := "Projects.Delete"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting project")
	project, err := p.projects.Project(ctx, id)
	if err != nil {
		log.Error("failed to get project for ownership check", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, project.OwnerId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = p.projects.DeleteProject(ctx, id)
	if err != nil {
		log.Error("failed to delete project", sl.Err(err))
		return err
	}
	return nil
}
