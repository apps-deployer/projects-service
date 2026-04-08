package deployconfigs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
	"github.com/apps-deployer/projects-service/internal/services"
)

type DeployConfigs struct {
	log      *slog.Logger
	storage  services.Storage
	dc       services.DeployConfigRepository
	projects services.ProjectRepository
}

func New(log *slog.Logger, storage services.Storage) *DeployConfigs {
	return &DeployConfigs{
		log:      log,
		storage:  storage,
		dc:       storage.Repos().DeployConfigs(),
		projects: storage.Repos().Projects(),
	}
}

func (c *DeployConfigs) checkProjectOwnership(ctx context.Context, projectID string) error {
	project, err := c.projects.Project(ctx, projectID)
	if err != nil {
		return err
	}
	return auth.CheckOwnership(ctx, project.OwnerId)
}

func (c *DeployConfigs) Resolve(ctx context.Context, projectId string) (*models.ResolvedDeployConfig, error) {
	op := "DeployConfigs.Resolve"
	log := c.log.With(
		slog.String("op", op),
		slog.String("projectId", projectId),
	)
	log.Info("resolving deploy config")
	if err := c.checkProjectOwnership(ctx, projectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	var res *models.ResolvedDeployConfig
	err := c.storage.WithinTx(ctx, func(tx services.RepoFactory) error {
		config, err := tx.DeployConfigs().DeployConfig(ctx, projectId)
		if err != nil {
			return err
		}
		framework, err := tx.Frameworks().Framework(ctx, config.FrameworkId)
		if err != nil {
			return err
		}
		res = models.NewResolvedDeployConfig(config, framework)
		return nil
	})
	if err != nil {
		log.Error("failed to resolve deploy config", sl.Err(err))
	}
	return res, nil
}

func (c *DeployConfigs) Get(ctx context.Context, projectId string) (*models.DeployConfig, error) {
	op := "DeployConfigs.Get"
	log := c.log.With(
		slog.String("op", op),
		slog.String("projectId", projectId),
	)
	log.Info("getting deploy config")
	if err := c.checkProjectOwnership(ctx, projectId); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return nil, err
	}
	config, err := c.dc.DeployConfig(ctx, projectId)
	if err != nil {
		log.Error("failed to get deploy config", sl.Err(err))
		return nil, err
	}
	return config, nil
}

func (c *DeployConfigs) Update(ctx context.Context, args *models.UpdateDeployConfigParams) error {
	op := "DeployConfigs.Update"
	log := c.log.With(
		slog.String("op", op),
		slog.String("id", args.Id),
	)
	log.Info("updating deploy config")
	ownerID, err := c.dc.ProjectOwnerByDeployConfigID(ctx, args.Id)
	if err != nil {
		log.Error("failed to get project owner for deploy config", sl.Err(err))
		return err
	}
	if err := auth.CheckOwnership(ctx, ownerID); err != nil {
		log.Warn("ownership check failed", sl.Err(err))
		return err
	}
	err = c.dc.UpdateDeployConfig(ctx, args)
	if err != nil {
		log.Error("failed to update deploy config", sl.Err(err))
		return err
	}
	return nil
}
