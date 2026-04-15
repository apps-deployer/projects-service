package models

import (
	"time"
)

type Env struct {
	Id           string
	Name         string
	ProjectId    string
	TargetBranch string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type GetEnvByGitParams struct {
	RepoUrl      string
	TargetBranch string
}

type CreateEnvParams struct {
	Name         string
	ProjectId    string
	TargetBranch string
}

type UpdateEnvParams struct {
	Id           string
	Name         *string
	TargetBranch *string
}

type ListEnvsParams struct {
	ProjectId string
	Limit     int64
	Offset    int64
}

type SaveEnvResponse struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEnvFromSaveResponse(
	args *CreateEnvParams,
	res *SaveEnvResponse,
) *Env {
	return &Env{
		Id:           res.Id,
		Name:         args.Name,
		ProjectId:    args.ProjectId,
		TargetBranch: args.TargetBranch,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
}
