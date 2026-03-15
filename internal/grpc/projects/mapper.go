package projects

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func projectToProto(p *models.Project) *projectsv1.ProjectResponse {
	return projectsv1.ProjectResponse_builder{
		Id:        &p.Id,
		Name:      &p.Name,
		RepoUrl:   &p.RepoUrl,
		OwnerId:   &p.OwnerId,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}.Build()
}

func projectsToProto(projects []*models.Project) *projectsv1.ListProjectsResponse {
	projectsResp := make([]*projectsv1.ProjectResponse, len(projects))
	for i, project := range projects {
		projectsResp[i] = projectToProto(project)
	}
	return projectsv1.ListProjectsResponse_builder{
		Projects: projectsResp,
	}.Build()
}

func protoToListProjectsParams(req *projectsv1.ListProjectsRequest) *models.ListProjectsParams {
	return &models.ListProjectsParams{
		OwnerId: req.GetOwnerId(),
		Limit:   req.GetLimit(),
		Offset:  req.GetOffset(),
	}
}

func protoToCreateProjectsParams(req *projectsv1.CreateProjectRequest) *models.CreateProjectParams {
	return &models.CreateProjectParams{
		Name:        req.GetName(),
		RepoUrl:     req.GetRepoUrl(),
		OwnerId:     req.GetOwnerId(),
		FrameworkId: req.GetDeployConfigTemplateId(),
	}
}

func protoToUpdateProjectParams(req *projectsv1.UpdateProjectRequest) *models.UpdateProjectParams {
	project := &models.UpdateProjectParams{Id: req.GetId()}
	if req.HasName() {
		name := req.GetName()
		project.Name = &name
	}
	if req.HasRepoUrl() {
		repoUrl := req.GetRepoUrl()
		project.RepoUrl = &repoUrl
	}
	if req.HasOwnerId() {
		ownerId := req.GetOwnerId()
		project.OwnerId = &ownerId
	}
	return project
}
