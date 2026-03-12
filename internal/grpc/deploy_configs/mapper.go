package deployconfigs

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
)

func deployConfigToProto(c *models.DeployConfig) *projectsv1.DeployConfigResponse {
	return projectsv1.DeployConfigResponse_builder{
		Id:                  &c.Id,
		ProjectId:           &c.ProjectId,
		FrameworkId:         &c.FrameworkId,
		RootDirOverwrite:    &c.RootDirOverride,
		OutputDirOverwrite:  &c.OutputDirOverride,
		BaseImageOverwrite:  &c.BaseImageOverride,
		InstallCmdOverwrite: &c.InstallCmdOverride,
		BuildCmdOverwrite:   &c.BuildCmdOverride,
		RunCmdOverwrite:     &c.RunCmdOverride,
	}.Build()
}

func generatedDeployConfigToProto(c *models.GeneratedDeployConfig) *projectsv1.GenerateDeployConfigResponse {
	return projectsv1.GenerateDeployConfigResponse_builder{
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
	if req.HasRootDirOverwrite() {
		rootDir := req.GetRootDirOverwrite()
		config.RootDirOverride = &rootDir
	}
	if req.HasOutputDirOverwrite() {
		outputDir := req.GetOutputDirOverwrite()
		config.OutputDirOverride = &outputDir
	}
	if req.HasBaseImageOverwrite() {
		baseImage := req.GetBaseImageOverwrite()
		config.BaseImageOverride = &baseImage
	}
	if req.HasInstallCmdOverwrite() {
		installCmd := req.GetInstallCmdOverwrite()
		config.InstallCmdOverride = &installCmd
	}
	if req.HasBuildCmdOverwrite() {
		buildCmd := req.GetBuildCmdOverwrite()
		config.BuildCmdOverride = &buildCmd
	}
	if req.HasRunCmdOverwrite() {
		runCmd := req.GetRunCmdOverwrite()
		config.RunCmdOverride = &runCmd
	}
	return config
}
