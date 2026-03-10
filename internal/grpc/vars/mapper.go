package vars

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
)

func varToProto(v *models.Var) *projectsv1.VarResponse {
	return projectsv1.VarResponse_builder{
		Id:    &v.Id,
		Key:   &v.Key,
		Value: &v.Value,
	}.Build()
}

func varsToProto(vars []*models.Var) *projectsv1.ListVarsResponse {
	varsResp := make([]*projectsv1.VarResponse, len(vars))
	for i, v := range vars {
		varsResp[i] = varToProto(v)
	}
	return projectsv1.ListVarsResponse_builder{
		Vars: varsResp,
	}.Build()
}

func protoToListProjectVarsParams(req *projectsv1.ListProjectVarsRequest) *models.ListProjectVarsParams {
	return &models.ListProjectVarsParams{
		ProjectId: req.GetProjectId(),
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	}
}

func protoToCreateProjectVarParams(req *projectsv1.CreateProjectVarRequest) *models.CreateProjectVarParams {
	return &models.CreateProjectVarParams{
		ProjectId: req.GetProjectId(),
		Key:       req.GetKey(),
		Value:     req.GetValue(),
	}
}

func protoToListEnvVarsParams(req *projectsv1.ListEnvVarsRequest) *models.ListEnvVarsParams {
	return &models.ListEnvVarsParams{
		EnvId:  req.GetEnvId(),
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}
}

func protoToCreateEnvVarParams(req *projectsv1.CreateEnvVarRequest) *models.CreateEnvVarParams {
	return &models.CreateEnvVarParams{
		EnvId: req.GetEnvId(),
		Key:   req.GetKey(),
		Value: req.GetValue(),
	}
}

func protoToUpdateVarParams(req *projectsv1.UpdateVarRequest) *models.UpdateVarParams {
	v := &models.UpdateVarParams{Id: req.GetId()}
	if req.HasValue() {
		value := req.GetValue()
		v.Value = &value
	}
	return v
}
