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
		RootDirOverwrite:    &c.RootDirOverwrite,
		OutputDirOverwrite:  &c.OutputDirOverwrite,
		BaseImageOverwrite:  &c.BaseImageOverwrite,
		InstallCmdOverwrite: &c.InstallCmdOverwrite,
		BuildCmdOverwrite:   &c.BuildCmdOverwrite,
		RunCmdOverwrite:     &c.RunCmdOverwrite,
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
		config.RootDirOverwrite = &rootDir
	}
	if req.HasOutputDirOverwrite() {
		outputDir := req.GetOutputDirOverwrite()
		config.OutputDirOverwrite = &outputDir
	}
	if req.HasBaseImageOverwrite() {
		baseImage := req.GetBaseImageOverwrite()
		config.BaseImageOverwrite = &baseImage
	}
	if req.HasInstallCmdOverwrite() {
		installCmd := req.GetInstallCmdOverwrite()
		config.InstallCmdOverwrite = &installCmd
	}
	if req.HasBuildCmdOverwrite() {
		buildCmd := req.GetBuildCmdOverwrite()
		config.BuildCmdOverwrite = &buildCmd
	}
	if req.HasRunCmdOverwrite() {
		runCmd := req.GetRunCmdOverwrite()
		config.RunCmdOverwrite = &runCmd
	}
	return config
}
