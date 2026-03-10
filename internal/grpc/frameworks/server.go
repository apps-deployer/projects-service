package frameworks

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type frameworksServer struct {
	projectsv1.UnimplementedFrameworkServiceServer
	frameworks FrameworksService
}

type FrameworksService interface {
	Get(ctx context.Context, id string) (*models.Framework, error)
	List(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error)
	Create(ctx context.Context, args *models.CreateFrameworkParams) (*models.Framework, error)
	Update(ctx context.Context, args *models.UpdateFrameworkParams) error
	Delete(ctx context.Context, id string) error
}

func Register(
	grpcServer *grpc.Server, frameworks FrameworksService) {
	projectsv1.RegisterFrameworkServiceServer(
		grpcServer,
		&frameworksServer{frameworks: frameworks},
	)
}

func (s *frameworksServer) GetFramework(
	ctx context.Context,
	req *projectsv1.GetFrameworkRequest,
) (*projectsv1.FrameworkResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	framework, err := s.frameworks.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get framework: %v", err)
	}
	if framework == nil {
		return nil, status.Error(codes.NotFound, "framework not found")
	}
	return frameworkToProto(framework), nil
}

func (s *frameworksServer) ListFrameworks(
	ctx context.Context,
	req *projectsv1.ListFrameworksRequest,
) (*projectsv1.ListFrameworksResponse, error) {
	frameworks, err := s.frameworks.List(
		ctx,
		protoToListFrameworksParams(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list frameworks: %v", err)
	}
	return frameworksToProto(frameworks), nil
}

func (s *frameworksServer) CreateFramework(
	ctx context.Context,
	req *projectsv1.CreateFrameworkRequest,
) (*projectsv1.FrameworkResponse, error) {
	if !req.HasName() {
		return nil, status.Error(codes.InvalidArgument, "framework name is required")
	}
	if !req.HasBaseImage() {
		return nil, status.Error(codes.InvalidArgument, "framework base image is required")
	}
	if !req.HasRunCmd() {
		return nil, status.Error(codes.InvalidArgument, "framework run command is required")
	}
	framework, err := s.frameworks.Create(ctx, protoToCreateFrameworkParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create framework: %v", err)
	}
	return frameworkToProto(framework), nil
}

func (s *frameworksServer) UpdateFramework(
	ctx context.Context,
	req *projectsv1.UpdateFrameworkRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	args := protoToUpdateFrameworkParams(req)
	if err := s.frameworks.Update(ctx, args); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update framework: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *frameworksServer) DeleteFramework(
	ctx context.Context,
	req *projectsv1.GetFrameworkRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	if err := s.frameworks.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete framework: %v", err)
	}
	return &emptypb.Empty{}, nil
}
