package models

import "time"

type Env struct {
	Id           string
	Name         string
	ProjectId    string
	TargetBranch string
	DomainName   string
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
	DomainName   string
}

type UpdateEnvParams struct {
	Id           string
	Name         *string
	TargetBranch *string
	DomainName   *string
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
