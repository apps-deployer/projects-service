package envs

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
)

func envToProto(e *models.Env) *projectsv1.EnvResponse {
	return projectsv1.EnvResponse_builder{
		Id:           &e.Id,
		Name:         &e.Name,
		ProjectId:    &e.ProjectId,
		TargetBranch: &e.TargetBranch,
		DomainName:   &e.DomainName,
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

func protoToListEnvParams(req *projectsv1.ListEnvsRequest) *models.ListEnvParams {
	return &models.ListEnvParams{
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
		DomainName:   req.GetDomainName(),
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
	if req.HasDomainName() {
		domainName := req.GetDomainName()
		env.DomainName = &domainName
	}
	return env
}
