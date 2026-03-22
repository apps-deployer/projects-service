package projectsgrpc

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type projectsServer struct {
	projectsv1.UnimplementedProjectServiceServer
	projects ProjectsService
}

type ProjectsService interface {
	Get(ctx context.Context, id string) (*models.Project, error)
	List(ctx context.Context, args *models.ListProjectsParams) ([]*models.Project, error)
	Create(ctx context.Context, args *models.CreateProjectParams) (*models.Project, error)
	Update(ctx context.Context, args *models.UpdateProjectParams) error
	Delete(ctx context.Context, id string) error
}

func Register(grpcServer *grpc.Server, projects ProjectsService) {
	projectsv1.RegisterProjectServiceServer(grpcServer, &projectsServer{projects: projects})
}

func (s *projectsServer) GetProject(
	ctx context.Context,
	req *projectsv1.GetProjectRequest,
) (*projectsv1.ProjectResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	project, err := s.projects.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get project: %v", err)
	}
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return projectToProto(project), nil
}

func (s *projectsServer) ListProjects(
	ctx context.Context,
	req *projectsv1.ListProjectsRequest,
) (*projectsv1.ListProjectsResponse, error) {
	if !req.HasOwnerId() {
		return nil, status.Error(codes.InvalidArgument, "owner ID is required")
	}
	projects, err := s.projects.List(ctx, protoToListProjectsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list projects: %v", err)
	}
	return projectsToProto(projects), nil
}

func (s *projectsServer) CreateProject(
	ctx context.Context,
	req *projectsv1.CreateProjectRequest,
) (*projectsv1.ProjectResponse, error) {
	if !req.HasName() {
		return nil, status.Error(codes.InvalidArgument, "project name is required")
	}
	if !req.HasRepoUrl() {
		return nil, status.Error(codes.InvalidArgument, "project repo URL is required")
	}
	if !req.HasOwnerId() {
		return nil, status.Error(codes.InvalidArgument, "project owner ID is required")
	}
	if !req.HasDeployConfigTemplateId() {
		return nil, status.Error(codes.InvalidArgument, "deploy config template is required")
	}
	project, err := s.projects.Create(ctx, protoToCreateProjectsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create project: %v", err)
	}
	return projectToProto(project), nil
}

func (s *projectsServer) UpdateProject(
	ctx context.Context,
	req *projectsv1.UpdateProjectRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	if err := s.projects.Update(ctx, protoToUpdateProjectParams(req)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update project: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *projectsServer) DeleteProject(
	ctx context.Context,
	req *projectsv1.DeleteProjectRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	if err := s.projects.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete project: %v", err)
	}
	return &emptypb.Empty{}, nil
}
