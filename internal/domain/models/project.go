package models

import (
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Project struct {
	Id        string
	Name      string
	RepoUrl   string
	OwnerId   string
	CreatedAt timestamp.Timestamp
}

type CreateProjectParams struct {
	Name                   string
	RepoUrl                string
	OwnerId                string
	DeployConfigTemplateId string
}

type UpdateProjectParams struct {
	Id      string
	Name    *string
	RepoUrl *string
	OwnerId *string
}

func (p *Project) ToProto() *projectsv1.ProjectResponse {
	return projectsv1.ProjectResponse_builder{
		Id:        &p.Id,
		Name:      &p.Name,
		RepoUrl:   &p.RepoUrl,
		OwnerId:   &p.OwnerId,
		CreatedAt: &p.CreatedAt,
	}.Build()
}
