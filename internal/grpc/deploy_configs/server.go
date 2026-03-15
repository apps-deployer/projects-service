package deployconfigs

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type deployConfigsServer struct {
	projectsv1.UnimplementedDeployConfigServiceServer
	deployConfigs DeployConfigsService
}

type DeployConfigsService interface {
	Resolve(ctx context.Context, projectId string) (*models.ResolvedDeployConfig, error)
	Get(ctx context.Context, projectId string) (*models.DeployConfig, error)
	Update(ctx context.Context, args *models.UpdateDeployConfigParams) error
}

func Register(grpcServer *grpc.Server, deployConfigs DeployConfigsService) {
	projectsv1.RegisterDeployConfigServiceServer(
		grpcServer,
		&deployConfigsServer{deployConfigs: deployConfigs},
	)
}

func (s *deployConfigsServer) ResolveDeployConfig(
	ctx context.Context,
	req *projectsv1.GetDeployConfigRequest,
) (*projectsv1.ResolveDeployConfigResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	config, err := s.deployConfigs.Resolve(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate deploy config: %v", err)
	}
	if config == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return resolvedDeployConfigToProto(config), nil
}

func (s *deployConfigsServer) GetDeployConfig(
	ctx context.Context,
	req *projectsv1.GetDeployConfigRequest,
) (*projectsv1.DeployConfigResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	config, err := s.deployConfigs.Get(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate deploy config: %v", err)
	}
	if config == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return deployConfigToProto(config), nil
}

func (s *deployConfigsServer) UpdateDeployConfig(
	ctx context.Context,
	req *projectsv1.UpdateDeployConfigRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "deploy config ID is required")
	}
	if err := s.deployConfigs.Update(ctx, protoToUpdateDeployConfigParams(req)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update deploy config: %v", err)
	}
	return &emptypb.Empty{}, nil
}
