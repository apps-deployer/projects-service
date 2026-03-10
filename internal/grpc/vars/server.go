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
	vars VarsService
}

type VarsService interface {
	GenerateVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error)

	GetProjectVar(ctx context.Context, id string) (*models.Var, error)
	ListProjectVars(ctx context.Context, args *models.ListProjectVarsParams) ([]*models.Var, error)
	CreateProjectVar(ctx context.Context, args *models.CreateProjectVarParams) (*models.Var, error)
	UpdateProjectVar(ctx context.Context, args *models.UpdateVarParams) error
	DeleteProjectVar(ctx context.Context, id string) error

	GetEnvVar(ctx context.Context, id string) (*models.Var, error)
	ListEnvVars(ctx context.Context, args *models.ListEnvVarsParams) ([]*models.Var, error)
	CreateEnvVar(ctx context.Context, args *models.CreateEnvVarParams) (*models.Var, error)
	UpdateEnvVar(ctx context.Context, args *models.UpdateVarParams) error
	DeleteEnvVar(ctx context.Context, id string) error
}

func Register(grpcServer *grpc.Server, vars VarsService) {
	projectsv1.RegisterVarServiceServer(grpcServer, &varsServer{vars: vars})
}

func (s *varsServer) ListVars(
	ctx context.Context,
	req *projectsv1.ListEnvVarsRequest,
) (*projectsv1.ListVarsResponse, error) {
	if !req.HasEnvId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	vars, err := s.vars.GenerateVars(ctx, newListEnvVarsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list vars: %v", err)
	}
	return newListVarsResponse(vars), nil
}

func (s *varsServer) GetProjectVar(
	ctx context.Context,
	req *projectsv1.GetVarRequest,
) (*projectsv1.VarResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "var ID is required")
	}
	v, err := s.vars.GetProjectVar(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get var: %v", err)
	}
	if v == nil {
		return nil, status.Error(codes.NotFound, "var not found")
	}
	return newVarResponse(v), nil
}

func (s *varsServer) ListProjectVars(
	ctx context.Context,
	req *projectsv1.ListProjectVarsRequest,
) (*projectsv1.ListVarsResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	vars, err := s.vars.ListProjectVars(ctx, newListProjectVarsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list project vars: %v", err)
	}
	return newListVarsResponse(vars), nil
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
	v, err := s.vars.CreateProjectVar(ctx, newCreateProjectVarParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create env: %v", err)
	}
	return newVarResponse(v), nil
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
	err := s.vars.UpdateProjectVar(ctx, newUpdateVarParams(req))
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
	if err := s.vars.DeleteProjectVar(ctx, req.GetId()); err != nil {
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
	v, err := s.vars.GetEnvVar(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get var: %v", err)
	}
	if v == nil {
		return nil, status.Error(codes.NotFound, "var not found")
	}
	return newVarResponse(v), nil
}

func (s *varsServer) ListEnvVars(
	ctx context.Context,
	req *projectsv1.ListEnvVarsRequest,
) (*projectsv1.ListVarsResponse, error) {
	if !req.HasEnvId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	vars, err := s.vars.ListEnvVars(ctx, newListEnvVarsParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list env vars: %v", err)
	}
	return newListVarsResponse(vars), nil
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
	v, err := s.vars.CreateEnvVar(ctx, newCreateEnvVarParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create env: %v", err)
	}
	return newVarResponse(v), nil
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
	err := s.vars.UpdateEnvVar(ctx, newUpdateVarParams(req))
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
	if err := s.vars.DeleteEnvVar(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete var: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func newListVarsResponse(vars []*models.Var) *projectsv1.ListVarsResponse {
	varsResp := make([]*projectsv1.VarResponse, len(vars))
	for i, v := range vars {
		varsResp[i] = newVarResponse(v)
	}
	return projectsv1.ListVarsResponse_builder{
		Vars: varsResp,
	}.Build()
}

func newVarResponse(v *models.Var) *projectsv1.VarResponse {
	return projectsv1.VarResponse_builder{
		Id:    &v.Id,
		Key:   &v.Key,
		Value: &v.Value,
	}.Build()
}

func newListProjectVarsParams(req *projectsv1.ListProjectVarsRequest) *models.ListProjectVarsParams {
	return &models.ListProjectVarsParams{
		ProjectId: req.GetProjectId(),
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	}
}

func newCreateProjectVarParams(req *projectsv1.CreateProjectVarRequest) *models.CreateProjectVarParams {
	return &models.CreateProjectVarParams{
		ProjectId: req.GetProjectId(),
		Key:       req.GetKey(),
		Value:     req.GetValue(),
	}
}

func newListEnvVarsParams(req *projectsv1.ListEnvVarsRequest) *models.ListEnvVarsParams {
	return &models.ListEnvVarsParams{
		EnvId:  req.GetEnvId(),
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}
}

func newCreateEnvVarParams(req *projectsv1.CreateEnvVarRequest) *models.CreateEnvVarParams {
	return &models.CreateEnvVarParams{
		EnvId: req.GetEnvId(),
		Key:   req.GetKey(),
		Value: req.GetValue(),
	}
}

func newUpdateVarParams(req *projectsv1.UpdateVarRequest) *models.UpdateVarParams {
	v := &models.UpdateVarParams{Id: req.GetId()}
	if req.HasValue() {
		value := req.GetValue()
		v.Value = &value
	}
	return v
}
