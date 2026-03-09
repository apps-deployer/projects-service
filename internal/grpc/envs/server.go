package envs

import (
	"context"

	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type envsServer struct {
	projectsv1.UnimplementedEnvServiceServer
	envs Envs
}

type Envs interface {
	GetEnvByGit(ctx context.Context, args *GetByGitParams) (*Env, error)
	Get(ctx context.Context, id string) (*Env, error)
	List(ctx context.Context, args *ListEnvParams) ([]*Env, error)
	Create(ctx context.Context, args *CreateEnvParams) (*Env, error)
	Update(ctx context.Context, args *UpdateEnvParams) error
	Delete(ctx context.Context, id string) error
}

func Register(grpcServer *grpc.Server, envs Envs) {
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
	env, err := s.envs.GetEnvByGit(ctx, &GetByGitParams{
		RepoUrl:      req.GetRepoUrl(),
		TargetBranch: req.GetTargetBranch(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get env: %v", err)
	}
	if env == nil {
		return nil, status.Error(codes.NotFound, "env not found")
	}
	return env.ToProto(), nil
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
		return nil, status.Errorf(codes.Internal, "failed to get env: %v", err)
	}
	if env == nil {
		return nil, status.Error(codes.NotFound, "env not found")
	}
	return env.ToProto(), nil
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
		&ListEnvParams{
			ProjectId: req.GetProjectId(),
			Limit:     req.GetLimit(),
			Offset:    req.GetOffset(),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list envs: %v", err)
	}
	envsResp := make([]*projectsv1.EnvResponse, len(envs))
	for i, env := range envs {
		envsResp[i] = env.ToProto()
	}
	return projectsv1.ListEnvsResponse_builder{
		Envs: envsResp,
	}.Build(), nil
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
	env, err := s.envs.Create(
		ctx,
		&CreateEnvParams{
			Name:         req.GetName(),
			ProjectId:    req.GetProjectId(),
			TargetBranch: req.GetTargetBranch(),
			DomainName:   req.GetDomainName(),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create env: %v", err)
	}
	return env.ToProto(), nil
}

func (s *envsServer) UpdateEnv(
	ctx context.Context,
	req *projectsv1.UpdateEnvRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	id := req.GetId()
	args := &UpdateEnvParams{Id: &id}
	if req.HasName() {
		name := req.GetName()
		args.Name = &name
	}
	if req.HasTargetBranch() {
		targetBranch := req.GetTargetBranch()
		args.TargetBranch = &targetBranch
	}
	if req.HasDomainName() {
		domainName := req.GetDomainName()
		args.DomainName = &domainName
	}
	if err := s.envs.Update(ctx, args); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update env: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *envsServer) DeleteEnv(
	ctx context.Context,
	req *projectsv1.GetEnvRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "env ID is required")
	}
	if err := s.envs.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete env: %v", err)
	}
	return &emptypb.Empty{}, nil
}

type Env struct {
	Id           string
	Name         string
	ProjectId    string
	TargetBranch string
	DomainName   string
}

type GetByGitParams struct {
	RepoUrl      string
	TargetBranch string
}

type CreateEnvParams struct {
	Name         string
	ProjectId    string
	TargetBranch string
	DomainName   string
}

type UpdateEnvParams struct {
	Id           *string
	Name         *string
	TargetBranch *string
	DomainName   *string
}

type ListEnvParams struct {
	ProjectId string
	Limit     int64
	Offset    int64
}

func (p *Env) ToProto() *projectsv1.EnvResponse {
	return projectsv1.EnvResponse_builder{
		Id:           &p.Id,
		Name:         &p.Name,
		ProjectId:    &p.ProjectId,
		TargetBranch: &p.TargetBranch,
		DomainName:   &p.DomainName,
	}.Build()
}
