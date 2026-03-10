package models

import projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"

type DeployConfig struct {
	Id                  string
	ProjectId           string
	FrameworkId         string
	RootDirOverwrite    string
	OutputDirOverwrite  string
	BaseImageOverwrite  string
	InstallCmdOverwrite string
	BuildCmdOverwrite   string
	RunCmdOverwrite     string
}

type UpdateDeployConfigParams struct {
	Id                  string
	FrameworkId         *string
	RootDirOverwrite    *string
	OutputDirOverwrite  *string
	BaseImageOverwrite  *string
	InstallCmdOverwrite *string
	BuildCmdOverwrite   *string
	RunCmdOverwrite     *string
}

type GeneratedDeployConfig struct {
	Id         string
	ProjectId  string
	RootDir    string
	OutputDir  string
	BaseImage  string
	InstallCmd string
	BuildCmd   string
	RunCmd     string
}

func (p *DeployConfig) ToProto() *projectsv1.DeployConfigResponse {
	return projectsv1.DeployConfigResponse_builder{
		Id:                  &p.Id,
		ProjectId:           &p.ProjectId,
		FrameworkId:         &p.FrameworkId,
		RootDirOverwrite:    &p.RootDirOverwrite,
		OutputDirOverwrite:  &p.OutputDirOverwrite,
		BaseImageOverwrite:  &p.BaseImageOverwrite,
		InstallCmdOverwrite: &p.InstallCmdOverwrite,
		BuildCmdOverwrite:   &p.BuildCmdOverwrite,
		RunCmdOverwrite:     &p.RunCmdOverwrite,
	}.Build()
}

func (p *GeneratedDeployConfig) ToProto() *projectsv1.GenerateDeployConfigResponse {
	return projectsv1.GenerateDeployConfigResponse_builder{
		Id:         &p.Id,
		ProjectId:  &p.ProjectId,
		RootDir:    &p.RootDir,
		OutputDir:  &p.OutputDir,
		BaseImage:  &p.BaseImage,
		InstallCmd: &p.InstallCmd,
		BuildCmd:   &p.BuildCmd,
		RunCmd:     &p.RunCmd,
	}.Build()
}
