package deployconfigsgrpc

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func deployConfigToProto(c *models.DeployConfig) *projectsv1.DeployConfigResponse {
	return projectsv1.DeployConfigResponse_builder{
		Id:                 &c.Id,
		ProjectId:          &c.ProjectId,
		FrameworkId:        &c.FrameworkId,
		RootDirOverride:    &c.RootDirOverride,
		OutputDirOverride:  &c.OutputDirOverride,
		BaseImageOverride:  &c.BaseImageOverride,
		InstallCmdOverride: &c.InstallCmdOverride,
		BuildCmdOverride:   &c.BuildCmdOverride,
		RunCmdOverride:     &c.RunCmdOverride,
		CreatedAt:          timestamppb.New(c.CreatedAt),
		UpdatedAt:          timestamppb.New(c.UpdatedAt),
	}.Build()
}

func resolvedDeployConfigToProto(c *models.ResolvedDeployConfig) *projectsv1.ResolveDeployConfigResponse {
	return projectsv1.ResolveDeployConfigResponse_builder{
		Id:         &c.Id,
		ProjectId:  &c.ProjectId,
		RootDir:    &c.RootDir,
		OutputDir:  &c.OutputDir,
		BaseImage:  &c.BaseImage,
		InstallCmd: &c.InstallCmd,
		BuildCmd:   &c.BuildCmd,
		RunCmd:     &c.RunCmd,
	}.Build()
}

func protoToUpdateDeployConfigParams(req *projectsv1.UpdateDeployConfigRequest) *models.UpdateDeployConfigParams {
	config := &models.UpdateDeployConfigParams{Id: req.GetId()}
	if req.HasFrameworkId() {
		frameworkId := req.GetFrameworkId()
		config.FrameworkId = &frameworkId
	}
	if req.HasRootDirOverride() {
		rootDir := req.GetRootDirOverride()
		config.RootDirOverride = &rootDir
	}
	if req.HasOutputDirOverride() {
		outputDir := req.GetOutputDirOverride()
		config.OutputDirOverride = &outputDir
	}
	if req.HasBaseImageOverride() {
		baseImage := req.GetBaseImageOverride()
		config.BaseImageOverride = &baseImage
	}
	if req.HasInstallCmdOverride() {
		installCmd := req.GetInstallCmdOverride()
		config.InstallCmdOverride = &installCmd
	}
	if req.HasBuildCmdOverride() {
		buildCmd := req.GetBuildCmdOverride()
		config.BuildCmdOverride = &buildCmd
	}
	if req.HasRunCmdOverride() {
		runCmd := req.GetRunCmdOverride()
		config.RunCmdOverride = &runCmd
	}
	return config
}
