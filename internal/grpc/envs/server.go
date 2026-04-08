package envsgrpc

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	grpcutil "github.com/apps-deployer/projects-service/internal/grpc"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type envsServer struct {
	projectsv1.UnimplementedEnvServiceServer
	envs EnvsService
}

type EnvsService interface {
	GetByGit(ctx context.Context, args *models.GetEnvByGitParams) (*models.Env, error)
	Get(ctx context.Context, id string) (*models.Env, error)
	List(ctx context.Context, args *models.ListEnvsParams) ([]*models.Env, error)
	Create(ctx context.Context, args *models.CreateEnvParams) (*models.Env, error)
	Update(ctx context.Context, args *models.UpdateEnvParams) error
	Delete(ctx context.Context, id string) error
}

func Register(grpcServer *grpc.Server, envs EnvsService) {
	projectsv1.RegisterEnvServiceServer(grpcServer, &envsServer{envs: envs})
}

func (s *envsServer) GetEnvByGit(
	ctx context.Context,
	req *projectsv1.GetEnvByGitRequest,
) (*projectsv1.EnvResponse, error) {
	if !req.HasRepoUrl() {
		return nil, status.Error(codes.InvalidArgument, "repo URL is required")
	}
	if !req.HasTargetBranch() {
		return nil, status.Error(codes.InvalidArgument, "target branch is required")
	}
	env, err := s.envs.GetByGit(ctx, protoToGetEnvByGitParams(req))
	if err != nil {
		return nil, grpcutil.MapError(err, "get env by git")
	}
	if env == nil {
		return nil, status.Error(codes.NotFound, "env not found")
	}
	return envToProto(env), nil
}

func (s *envsServer) GetEnv(
	ctx context.Context,
	req *projectsv1.GetEnvRequest,
) (*projectsv1.EnvResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	env, err := s.envs.Get(ctx, req.GetId())
	if err != nil {
		return nil, grpcutil.MapError(err, "get env")
	}
	if env == nil {
		return nil, status.Error(codes.NotFound, "env not found")
	}
	return envToProto(env), nil
}

func (s *envsServer) ListEnvs(
	ctx context.Context,
	req *projectsv1.ListEnvsRequest,
) (*projectsv1.ListEnvsResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	envs, err := s.envs.List(
		ctx,
		protoToListEnvParams(req),
	)
	if err != nil {
		return nil, grpcutil.MapError(err, "list envs")
	}
	return envsToProto(envs), nil
}

func (s *envsServer) CreateEnv(
	ctx context.Context,
	req *projectsv1.CreateEnvRequest,
) (*projectsv1.EnvResponse, error) {
	if !req.HasName() {
		return nil, status.Error(codes.InvalidArgument, "env name is required")
	}
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	if !req.HasTargetBranch() {
		return nil, status.Error(codes.InvalidArgument, "env target branch is required")
	}
	if !req.HasDomainName() {
		return nil, status.Error(codes.InvalidArgument, "env domain name is required")
	}
	env, err := s.envs.Create(ctx, protoToCreateEnvParams(req))
	if err != nil {
		return nil, grpcutil.MapError(err, "create env")
	}
	return envToProto(env), nil
}

func (s *envsServer) UpdateEnv(
	ctx context.Context,
	req *projectsv1.UpdateEnvRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}

	if err := s.envs.Update(ctx, protoToUpdateEnvParams(req)); err != nil {
		return nil, grpcutil.MapError(err, "update env")
	}
	return &emptypb.Empty{}, nil
}

func (s *envsServer) DeleteEnv(
	ctx context.Context,
	req *projectsv1.DeleteEnvRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	if err := s.envs.Delete(ctx, req.GetId()); err != nil {
		return nil, grpcutil.MapError(err, "delete env")
	}
	return &emptypb.Empty{}, nil
}
