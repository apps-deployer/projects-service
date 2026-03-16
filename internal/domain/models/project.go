package models

import (
	"time"
)

type Project struct {
	Id        string
	Name      string
	RepoUrl   string
	OwnerId   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListProjectsParams struct {
	OwnerId string
	Limit   int64
	Offset  int64
}

type CreateProjectParams struct {
	Name        string
	RepoUrl     string
	OwnerId     string
	FrameworkId string
}

type UpdateProjectParams struct {
	Id      string
	Name    *string
	RepoUrl *string
	OwnerId *string
}

type SaveProjectParams struct {
	Name    string
	RepoUrl string
	OwnerId string
}

type SaveProjectResponse struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProjectFromSaveResponse(
	args *SaveProjectParams,
	res *SaveProjectResponse,
) *Project {
	return &Project{
		Id:        res.Id,
		Name:      args.Name,
		RepoUrl:   args.RepoUrl,
		OwnerId:   args.OwnerId,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}
}
