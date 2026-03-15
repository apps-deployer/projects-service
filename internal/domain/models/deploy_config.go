package models

import "time"

type DeployConfig struct {
	Id                 string
	ProjectId          string
	FrameworkId        string
	RootDirOverride    string
	OutputDirOverride  string
	BaseImageOverride  string
	InstallCmdOverride string
	BuildCmdOverride   string
	RunCmdOverride     string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type UpdateDeployConfigParams struct {
	Id                 string
	FrameworkId        *string
	RootDirOverride    *string
	OutputDirOverride  *string
	BaseImageOverride  *string
	InstallCmdOverride *string
	BuildCmdOverride   *string
	RunCmdOverride     *string
}

type ResolvedDeployConfig struct {
	Id         string
	ProjectId  string
	RootDir    string
	OutputDir  string
	BaseImage  string
	InstallCmd string
	BuildCmd   string
	RunCmd     string
}

type SaveDeployConfigParams struct {
	ProjectId   string
	FrameworkId string
}
