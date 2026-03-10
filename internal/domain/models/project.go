package models

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Project struct {
	Id        string
	Name      string
	RepoUrl   string
	OwnerId   string
	CreatedAt timestamp.Timestamp
}

type ListProjectsParams struct {
	OwnerId string
	Limit   int64
	Offset  int64
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
