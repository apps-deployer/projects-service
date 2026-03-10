package models

import projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"

type Framework struct {
	Id         string
	Name       string
	RootDir    string
	OutputDir  string
	BaseImage  string
	InstallCmd string
	BuildCmd   string
	RunCmd     string
}

type CreateFrameworkParams struct {
	Name       string
	RootDir    string
	OutputDir  string
	BaseImage  string
	InstallCmd string
	BuildCmd   string
	RunCmd     string
}

type UpdateFrameworkParams struct {
	Id         string
	Name       *string
	RootDir    *string
	OutputDir  *string
	BaseImage  *string
	InstallCmd *string
	BuildCmd   *string
	RunCmd     *string
}

func (p *Framework) ToProto() *projectsv1.FrameworkResponse {
	return projectsv1.FrameworkResponse_builder{
		Id:         &p.Id,
		Name:       &p.Name,
		RootDir:    &p.RootDir,
		OutputDir:  &p.OutputDir,
		BaseImage:  &p.BaseImage,
		InstallCmd: &p.InstallCmd,
		BuildCmd:   &p.BuildCmd,
		RunCmd:     &p.RunCmd,
	}.Build()
}
