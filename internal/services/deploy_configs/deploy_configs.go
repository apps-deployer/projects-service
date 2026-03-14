package deployconfigs

import (
	"context"
	"log/slog"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/lib/logger/sl"
)

type DeployConfigStorage interface {
	DeployConfig(ctx context.Context, projectId string) (*models.DeployConfig, error)
	UpdateDeployConfig(ctx context.Context, args *models.UpdateDeployConfigParams) error
}

type FrameworkStorage interface {
	Framework(ctx context.Context, id string) (*models.Framework, error)
}

type DeployConfigs struct {
	log           *slog.Logger
	deployConfigs DeployConfigStorage
	frameworks    FrameworkStorage
}

func New(log *slog.Logger, deployConfigs DeployConfigStorage, frameworks FrameworkStorage) *DeployConfigs {
	return &DeployConfigs{log: log, deployConfigs: deployConfigs, frameworks: frameworks}
}

func (c *DeployConfigs) Generate(ctx context.Context, projectId string) (*models.GeneratedDeployConfig, error) {
	// TODO: Auth

	op := "DeployConfigs.Generate"
	log := c.log.With(
		slog.String("op", op),
		slog.String("projectId", projectId),
	)
	log.Info("generating deploy config")

	config, err := c.deployConfigs.DeployConfig(ctx, projectId)
	if err != nil {
		log.Error("failed to get deploy config", sl.Err(err))
		return nil, err
	}
	framework, err := c.frameworks.Framework(ctx, config.FrameworkId)
	if err != nil {
		log.Error("failed to get framework", sl.Err(err))
		return nil, err
	}
	res := &models.GeneratedDeployConfig{
		Id:         config.Id,
		ProjectId:  config.ProjectId,
		RootDir:    pick(config.RootDirOverride, framework.RootDir),
		OutputDir:  pick(config.OutputDirOverride, framework.OutputDir),
		BaseImage:  pick(config.BaseImageOverride, framework.BaseImage),
		InstallCmd: pick(config.InstallCmdOverride, framework.InstallCmd),
		BuildCmd:   pick(config.BuildCmdOverride, framework.BuildCmd),
		RunCmd:     pick(config.RunCmdOverride, framework.RunCmd),
	}
	return res, nil
}

func pick(override, base string) string {
	if override != "" {
		return override
	}
	return base
}

func (c *DeployConfigs) Get(ctx context.Context, projectId string) (*models.DeployConfig, error) {
	// TODO: Auth
	op := "DeployConfigs.Get"
	log := c.log.With(
		slog.String("op", op),
		slog.String("projectId", projectId),
	)
	log.Info("getting deploy config")

	config, err := c.deployConfigs.DeployConfig(ctx, projectId)
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

	err := c.deployConfigs.UpdateDeployConfig(ctx, args)
	if err != nil {
		log.Error("failed to update deploy config", sl.Err(err))
		return err
	}
	return nil
}
