package projects

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type ProjectStorage interface {
	Project(ctx context.Context, id string) (*models.Project, error)
	ListProjects(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error)
	SaveProject(ctx context.Context, args *models.SaveProjectParams) (*models.SaveProjectResponse, error)
	UpdateProject(ctx context.Context, args *models.UpdateProjectParams) error
	DeleteProject(ctx context.Context, id string) error
}

type DeployConfigStorage interface {
	SaveDeployConfig(ctx context.Context, args *models.SaveDeployConfigParams) (id string, err error)
	DeleteDeployConfig(ctx context.Context, projectId string) error
}

type Projects struct {
	log           *slog.Logger
	projects      ProjectStorage
	deployConfigs DeployConfigStorage
}

func New(
	log *slog.Logger,
	projects ProjectStorage,
	deployConfigs DeployConfigStorage,
) *Projects {
	return &Projects{
		log:           log,
		projects:      projects,
		deployConfigs: deployConfigs,
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
	res, err := p.projects.Project(ctx, id)
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
	res, err := p.projects.ListProjects(ctx, args)
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
	res, err := p.projects.SaveProject(ctx, project)
	if err != nil {
		log.Error("failed to save project", sl.Err(err))
		return nil, err
	}
	_, err = p.deployConfigs.SaveDeployConfig(
		ctx, &models.SaveDeployConfigParams{ProjectId: res.Id, FrameworkId: args.FrameworkId})
	if err != nil {
		log.Error("failed to save deploy config", sl.Err(err))
		return nil, err
	}
	return &models.Project{
		Id:        res.Id,
		Name:      args.Name,
		RepoUrl:   args.RepoUrl,
		OwnerId:   args.OwnerId,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (p *Projects) Update(ctx context.Context, args *models.UpdateProjectParams) error {
	// TODO: Auth
	op := "Projects.Update"
	log := p.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating project")
	err := p.projects.UpdateProject(ctx, args)
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
	err := p.projects.DeleteProject(ctx, id)
	if err != nil {
		log.Error("failed to delete project", sl.Err(err))
		return err
	}
	return nil
}
