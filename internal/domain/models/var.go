package models

import projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"

type Var struct {
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
	Value string
}

func (p *Var) ToProto() *projectsv1.VarResponse {
	return projectsv1.VarResponse_builder{
		Id:    &p.Id,
		Key:   &p.Key,
		Value: &p.Value,
	}.Build()
}
