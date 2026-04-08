package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	deployconfigsgrpc "github.com/apps-deployer/projects-service/internal/grpc/deploy_configs"
	envsgrpc "github.com/apps-deployer/projects-service/internal/grpc/envs"
	frameworksgrpc "github.com/apps-deployer/projects-service/internal/grpc/frameworks"
	projectsgrpc "github.com/apps-deployer/projects-service/internal/grpc/projects"
	varsgrpc "github.com/apps-deployer/projects-service/internal/grpc/vars"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int // Порт, на котором будет работать grpc-сервер
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	deployConfigService deployconfigsgrpc.DeployConfigsService,
	envService envsgrpc.EnvsService,
	frameworkService frameworksgrpc.FrameworksService,
	projectService projectsgrpc.ProjectsService,
	projectVarService varsgrpc.ProjectVarsService,
	envVarService varsgrpc.EnvVarsService,
	varsAggregationService varsgrpc.VarsAggregationService,
	port int,
	jwtSecret string,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		AuthInterceptor(jwtSecret),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	deployconfigsgrpc.Register(grpcServer, deployConfigService)
	envsgrpc.Register(grpcServer, envService)
	frameworksgrpc.Register(grpcServer, frameworkService)
	projectsgrpc.Register(grpcServer, projectService)
	varsgrpc.Register(grpcServer, projectVarService, envVarService, varsAggregationService)

	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	// Используем встроенный в gRPCServer механизм graceful shutdown
	a.grpcServer.GracefulStop()
}
