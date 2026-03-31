package app

import (
	"log/slog"

	grpcapp "github.com/apps-deployer/projects-service/internal/app/grpc"
	"github.com/apps-deployer/projects-service/internal/services"
	deployconfigs "github.com/apps-deployer/projects-service/internal/services/deploy_configs"
	"github.com/apps-deployer/projects-service/internal/services/envs"
	"github.com/apps-deployer/projects-service/internal/services/frameworks"
	"github.com/apps-deployer/projects-service/internal/services/projects"
	"github.com/apps-deployer/projects-service/internal/services/vars"
	"github.com/apps-deployer/projects-service/internal/storage/postgres"
)

type App struct {
	GrpcServer *grpcapp.App
	storage    services.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	dbUrl string,
	encryptionKey string,
) *App {
	storage, err := postgres.New(dbUrl, encryptionKey)
	if err != nil {
		panic(err)
	}

	deployConfigService := deployconfigs.New(log, storage)
	envService := envs.New(log, storage)
	frameworkService := frameworks.New(log, storage)
	projectService := projects.New(log, storage)
	varService := vars.New(log, storage)

	grpcApp := grpcapp.New(
		log,
		deployConfigService,
		envService,
		frameworkService,
		projectService,
		varService,
		varService,
		varService,
		grpcPort,
	)

	return &App{
		GrpcServer: grpcApp,
		storage:    storage,
	}
}

func (a *App) Stop() {
	a.GrpcServer.Stop()
	a.storage.Stop()
}
