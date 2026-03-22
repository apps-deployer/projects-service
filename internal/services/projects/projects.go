package projects

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type Storage interface {
	Projects() ProjectRepository
	DeployConfigs() DeployConfigRepository

	WithinTx(
		ctx context.Context,
		fn func(TxStorage) error,
	) error
}

type TxStorage interface {
	Projects() ProjectRepository
	DeployConfigs() DeployConfigRepository
}

type ProjectRepository interface {
	Project(ctx context.Context, id string) (*models.Project, error)
	ListProjects(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error)
	SaveProject(ctx context.Context, args *models.SaveProjectParams) (*models.SaveProjectResponse, error)
	UpdateProject(ctx context.Context, args *models.UpdateProjectParams) error
	DeleteProject(ctx context.Context, id string) error
}

type DeployConfigRepository interface {
	SaveDeployConfig(ctx context.Context, args *models.SaveDeployConfigParams) (*models.SaveDeployConfigResponse, error)
	DeleteDeployConfig(ctx context.Context, projectId string) error
}

type Projects struct {
	log     *slog.Logger
	storage Storage
}

func New(
	log *slog.Logger,
	storage Storage,
) *Projects {
	return &Projects{
		log:     log,
		storage: storage,
	}
}

func (p *Projects) Get(ctx context.Context, id string) (*models.Project, error) {
	// TODO: Auth
	op := "Projects.Get"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("getting project")
	res, err := p.storage.Projects().Project(ctx, id)
	if err != nil {
		log.Error("failed to get project", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (p *Projects) List(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error) {
	// TODO: Auth
	op := "Projects.List"
	log := p.log.With(
		slog.String("op", op),
		slog.String("ownerId", args.OwnerId),
	)
	log.Info("listing projects")
	res, err := p.storage.Projects().ListProjects(ctx, args)
	if err != nil {
		log.Error("failed to list projects", sl.Err(err))
		return nil, err
	}
	return res, nil
}

func (p *Projects) Create(ctx context.Context, args *models.CreateProjectParams) (*models.Project, error) {
	// TODO: Auth
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
	err := p.storage.WithinTx(ctx, func(tx TxStorage) error {
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
	// TODO: Auth
	op := "Projects.Update"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating project")
	err := p.storage.Projects().UpdateProject(ctx, args)
	if err != nil {
		log.Error("failed to update project", sl.Err(err))
		return err
	}
	return nil
}

func (p *Projects) Delete(ctx context.Context, id string) error {
	// TODO: Auth
	op := "Projects.Delete"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	log.Info("deleting project")
	err := p.storage.Projects().DeleteProject(ctx, id)
	if err != nil {
		log.Error("failed to delete project", sl.Err(err))
		return err
	}
	return nil
}
