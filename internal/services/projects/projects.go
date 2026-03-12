package projects

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
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

	res, err := p.projects.Project(ctx, id)
	return res, err
}

func (p *Projects) List(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error) {
	// TODO: Auth

	res, err := p.projects.ListProjects(ctx, args)
	return res, err
}

func (p *Projects) Create(ctx context.Context, args *models.CreateProjectParams) (*models.Project, error) {
	// TODO: Auth

	project := &models.SaveProjectParams{
		Name:    args.Name,
		RepoUrl: args.RepoUrl,
		OwnerId: args.OwnerId,
	}
	res, err := p.projects.SaveProject(ctx, project)
	if err != nil {
		return nil, err
	}
	_, err = p.deployConfigs.SaveDeployConfig(
		ctx, &models.SaveDeployConfigParams{ProjectId: res.Id, FrameworkId: args.FrameworkId})
	if err != nil {
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

	err := p.projects.UpdateProject(ctx, args)
	return err
}

func (p *Projects) Delete(ctx context.Context, id string) error {
	// TODO: Auth

	err := p.projects.DeleteProject(ctx, id)
	return err
}
