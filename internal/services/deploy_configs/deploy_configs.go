package deployconfigs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type Storage interface {
	DeployConfigs() DeployConfigRepository
	Frameworks() FrameworkRepository

	WithinTx(
		ctx context.Context,
		fn func(TxStorage) (*models.ResolvedDeployConfig, error),
	) (*models.ResolvedDeployConfig, error)
}

type TxStorage interface {
	DeployConfigs() DeployConfigRepository
	Frameworks() FrameworkRepository
}

type DeployConfigRepository interface {
	DeployConfig(ctx context.Context, projectId string) (*models.DeployConfig, error)
	UpdateDeployConfig(ctx context.Context, args *models.UpdateDeployConfigParams) error
}

type FrameworkRepository interface {
	Framework(ctx context.Context, id string) (*models.Framework, error)
}

type DeployConfigs struct {
	log     *slog.Logger
	storage Storage
}

func New(log *slog.Logger, storage Storage) *DeployConfigs {
	return &DeployConfigs{
		log:     log,
		storage: storage,
	}
}

func (c *DeployConfigs) Resolve(ctx context.Context, projectId string) (*models.ResolvedDeployConfig, error) {
	// TODO: Auth
	op := "DeployConfigs.Generate"
	log := c.log.With(
		slog.String("op", op),
		slog.String("projectId", projectId),
	)
	log.Info("resolving deploy config")
	res, err := c.storage.WithinTx(ctx, func(tx TxStorage) (*models.ResolvedDeployConfig, error) {
		config, err := tx.DeployConfigs().DeployConfig(ctx, projectId)
		if err != nil {
			return nil, err
		}
		framework, err := tx.Frameworks().Framework(ctx, config.FrameworkId)
		if err != nil {
			return nil, err
		}
		return models.NewResolvedDeployConfig(config, framework), nil
	})
	if err != nil {
		log.Error("failed to resolve deploy config", sl.Err(err))
	}
	return res, nil
}

func (c *DeployConfigs) Get(ctx context.Context, projectId string) (*models.DeployConfig, error) {
	// TODO: Auth
	op := "DeployConfigs.Get"
	log := c.log.With(
		slog.String("op", op),
		slog.String("projectId", projectId),
	)
	log.Info("getting deploy config")

	config, err := c.storage.DeployConfigs().DeployConfig(ctx, projectId)
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

	err := c.storage.DeployConfigs().UpdateDeployConfig(ctx, args)
	if err != nil {
		log.Error("failed to update deploy config", sl.Err(err))
		return err
	}
	return nil
}
