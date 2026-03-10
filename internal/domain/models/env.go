package models

import projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"

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
	Id           string
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
