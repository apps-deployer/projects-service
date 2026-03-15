package models

import "time"

type Var struct {
	Id        string
	Key       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResolvedVar struct {
	Id    string
	Key   string
	Value string
}

type ListProjectVarsParams struct {
	ProjectId string
	Limit     int64
	Offset    int64
}

type ListEnvVarsParams struct {
	EnvId  string
	Limit  int64
	Offset int64
}

type CreateProjectVarParams struct {
	ProjectId string
	Key       string
	Value     string
}

type CreateEnvVarParams struct {
	EnvId string
	Key   string
	Value string
}

type UpdateVarParams struct {
	Id    string
	Value *string
}

type SaveVarResponse struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
}
