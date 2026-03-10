package vars

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type varsServer struct {
	projectsv1.UnimplementedVarServiceServer
	projectVars    ProjectVarsService
	envVars        EnvVarsService
	varsAggregator VarsAggregationService
}

type ProjectVarsService interface {
	Get(ctx context.Context, id string) (*models.Var, error)
	List(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error)
	Create(ctx context.Context, args *models.CreateProjectVarParams) (*models.Var, error)
	Update(ctx context.Context, args *models.UpdateVarParams) error
	Delete(ctx context.Context, id string) error
}

type EnvVarsService interface {
	Get(ctx context.Context, id string) (*models.Var, error)
	List(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error)
	Create(ctx context.Context, args *models.CreateEnvVarParams) (*models.Var, error)
	Update(ctx context.Context, args *models.UpdateVarParams) error
	Delete(ctx context.Context, id string) error
}

type VarsAggregationService interface {
	ListAllVars(ctx context.Context, envId string) ([]*models.Var, error)
}

func Register(
	grpcServer *grpc.Server,
	projectVars ProjectVarsService,
	envVars EnvVarsService,
	varsAggregator VarsAggregationService,
) {
	projectsv1.RegisterVarServiceServer(
		grpcServer,
		&varsServer{
			projectVars:    projectVars,
			envVars:        envVars,
			varsAggregator: varsAggregator,
		},
	)
}

func (s *varsServer) ListAllVars(
	ctx context.Context,
	req *projectsv1.ListAllVarsRequest,
) (*projectsv1.ListVarsResponse, error) {
	if !req.HasEnvId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	vars, err := s.varsAggregator.ListAllVars(ctx, req.GetEnvId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list vars: %v", err)
	}
	return varsToProto(vars), nil
}

func (s *varsServer) GetProjectVar(
	ctx context.Context,
	req *projectsv1.GetVarRequest,
) (*projectsv1.VarResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	v, err := s.projectVars.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get var: %v", err)
	}
	if v == nil {
		return nil, status.Error(codes.NotFound, "var not found")
	}
	return varToProto(v), nil
}

func (s *varsServer) ListProjectVars(
	ctx context.Context,
	req *projectsv1.ListProjectVarsRequest,
) (*projectsv1.ListVarsResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	vars, err := s.projectVars.List(ctx, protoToListProjectVarsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list project vars: %v", err)
	}
	return varsToProto(vars), nil
}

func (s *varsServer) CreateProjectVar(
	ctx context.Context,
	req *projectsv1.CreateProjectVarRequest,
) (*projectsv1.VarResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	if !req.HasKey() {
		return nil, status.Error(codes.InvalidArgument, "key is required")
	}
	if !req.HasValue() {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}
	v, err := s.projectVars.Create(ctx, protoToCreateProjectVarParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create env: %v", err)
	}
	return varToProto(v), nil
}

func (s *varsServer) UpdateProjectVar(
	ctx context.Context,
	req *projectsv1.UpdateVarRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	if !req.HasValue() {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}
	err := s.projectVars.Update(ctx, protoToUpdateVarParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update var: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *varsServer) DeleteProjectVar(
	ctx context.Context,
	req *projectsv1.GetVarRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	if err := s.projectVars.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete var: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *varsServer) GetEnvVar(
	ctx context.Context,
	req *projectsv1.GetVarRequest,
) (*projectsv1.VarResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	v, err := s.envVars.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get var: %v", err)
	}
	if v == nil {
		return nil, status.Error(codes.NotFound, "var not found")
	}
	return varToProto(v), nil
}

func (s *varsServer) ListEnvVars(
	ctx context.Context,
	req *projectsv1.ListEnvVarsRequest,
) (*projectsv1.ListVarsResponse, error) {
	if !req.HasEnvId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	vars, err := s.envVars.List(ctx, protoToListEnvVarsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list env vars: %v", err)
	}
	return varsToProto(vars), nil
}

func (s *varsServer) CreateEnvVar(
	ctx context.Context,
	req *projectsv1.CreateEnvVarRequest,
) (*projectsv1.VarResponse, error) {
	if !req.HasEnvId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	if !req.HasKey() {
		return nil, status.Error(codes.InvalidArgument, "key is required")
	}
	if !req.HasValue() {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}
	v, err := s.envVars.Create(ctx, protoToCreateEnvVarParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create env: %v", err)
	}
	return varToProto(v), nil
}

func (s *varsServer) UpdateEnvVar(
	ctx context.Context,
	req *projectsv1.UpdateVarRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	if !req.HasValue() {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}
	err := s.envVars.Update(ctx, protoToUpdateVarParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update var: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *varsServer) DeleteEnvVar(
	ctx context.Context,
	req *projectsv1.GetVarRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	if err := s.envVars.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete var: %v", err)
	}
	return &emptypb.Empty{}, nil
}
