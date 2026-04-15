package envsgrpc

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func envToProto(e *models.Env) *projectsv1.EnvResponse {
	return projectsv1.EnvResponse_builder{
		Id:           &e.Id,
		Name:         &e.Name,
		ProjectId:    &e.ProjectId,
		TargetBranch: &e.TargetBranch,

		CreatedAt: timestamppb.New(e.CreatedAt),
		UpdatedAt: timestamppb.New(e.UpdatedAt),
	}.Build()
}

func envsToProto(envs []*models.Env) *projectsv1.ListEnvsResponse {
	envsProto := make([]*projectsv1.EnvResponse, len(envs))
	for i, env := range envs {
		envsProto[i] = envToProto(env)
	}
	return projectsv1.ListEnvsResponse_builder{
		Envs: envsProto,
	}.Build()
}

func protoToGetEnvByGitParams(req *projectsv1.GetEnvByGitRequest) *models.GetEnvByGitParams {
	return &models.GetEnvByGitParams{
		RepoUrl:      req.GetRepoUrl(),
		TargetBranch: req.GetTargetBranch(),
	}
}

func protoToListEnvParams(req *projectsv1.ListEnvsRequest) *models.ListEnvsParams {
	return &models.ListEnvsParams{
		ProjectId: req.GetProjectId(),
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	}
}

func protoToCreateEnvParams(req *projectsv1.CreateEnvRequest) *models.CreateEnvParams {
	return &models.CreateEnvParams{
		Name:         req.GetName(),
		ProjectId:    req.GetProjectId(),
		TargetBranch: req.GetTargetBranch(),
	}
}

func protoToUpdateEnvParams(req *projectsv1.UpdateEnvRequest) *models.UpdateEnvParams {
	env := &models.UpdateEnvParams{Id: req.GetId()}
	if req.HasName() {
		name := req.GetName()
		env.Name = &name
	}
	if req.HasTargetBranch() {
		targetBranch := req.GetTargetBranch()
		env.TargetBranch = &targetBranch
	}
	return env
}
