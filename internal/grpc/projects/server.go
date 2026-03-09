package projects

import (
	"context"

	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type projectsServer struct {
	projectsv1.UnimplementedProjectServiceServer
	projects Projects
}

type Projects interface {
	Get(ctx context.Context, id string) (*ProjectDTO, error)
	List(ctx context.Context, ownerId string) ([]*ProjectDTO, error)
	Create(ctx context.Context, project *CreateProjectDTO) (*ProjectDTO, error)
	Update(ctx context.Context, project *UpdateProjectDTO) error
	Delete(ctx context.Context, id string) error
}

func Register(grpcServer *grpc.Server, projects Projects) {
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
	return project.ToProto(), nil
}

func (s *projectsServer) ListProjects(
	ctx context.Context,
	req *projectsv1.ListProjectsRequest,
) (*projectsv1.ListProjectsResponse, error) {
	projects, err := s.projects.List(ctx, req.GetOwnerId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list projects: %v", err)
	}
	projectsResp := make([]*projectsv1.ProjectResponse, len(projects))
	for i, project := range projects {
		projectsResp[i] = project.ToProto()
	}
	return projectsv1.ListProjectsResponse_builder{
		Projects: projectsResp,
	}.Build(), nil
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
	var (
		name                   = req.GetName()
		repoUrl                = req.GetRepoUrl()
		ownerId                = req.GetOwnerId()
		deployConfigTemplateId = req.GetDeployConfigTemplateId()
	)
	project, err := s.projects.Create(ctx, &CreateProjectDTO{
		Name:                   &name,
		RepoUrl:                &repoUrl,
		OwnerId:                &ownerId,
		DeployConfigTemplateId: &deployConfigTemplateId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create project: %v", err)
	}
	return project.ToProto(), nil
}

func (s *projectsServer) UpdateProject(
	ctx context.Context,
	req *projectsv1.UpdateProjectRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	id := req.GetId()
	project := &UpdateProjectDTO{Id: &id}
	if req.HasName() {
		name := req.GetName()
		project.Name = &name
	}
	if req.HasRepoUrl() {
		repoUrl := req.GetRepoUrl()
		project.RepoUrl = &repoUrl
	}
	if req.HasOwnerId() {
		ownerId := req.GetOwnerId()
		project.OwnerId = &ownerId
	}
	if err := s.projects.Update(ctx, project); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update project: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *projectsServer) DeleteProject(
	ctx context.Context,
	req *projectsv1.GetProjectRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	if err := s.projects.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete project: %v", err)
	}
	return &emptypb.Empty{}, nil
}

type ProjectDTO struct {
	Id        *string
	Name      *string
	RepoUrl   *string
	OwnerId   *string
	CreatedAt *timestamp.Timestamp
}

type CreateProjectDTO struct {
	Name                   *string
	RepoUrl                *string
	OwnerId                *string
	DeployConfigTemplateId *string
}

type UpdateProjectDTO struct {
	Id      *string
	Name    *string
	RepoUrl *string
	OwnerId *string
}

func (p *ProjectDTO) ToProto() *projectsv1.ProjectResponse {
	return projectsv1.ProjectResponse_builder{
		Id:        p.Id,
		Name:      p.Name,
		RepoUrl:   p.RepoUrl,
		OwnerId:   p.OwnerId,
		CreatedAt: p.CreatedAt,
	}.Build()
}
