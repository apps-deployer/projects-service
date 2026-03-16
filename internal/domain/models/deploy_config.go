package models

import (
	"time"
)

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

func NewResolvedDeployConfig(
	config *DeployConfig,
	framework *Framework,
) *ResolvedDeployConfig {
	return &ResolvedDeployConfig{
		Id:         config.Id,
		ProjectId:  config.ProjectId,
		RootDir:    pick(config.RootDirOverride, framework.RootDir),
		OutputDir:  pick(config.OutputDirOverride, framework.OutputDir),
		BaseImage:  pick(config.BaseImageOverride, framework.BaseImage),
		InstallCmd: pick(config.InstallCmdOverride, framework.InstallCmd),
		BuildCmd:   pick(config.BuildCmdOverride, framework.BuildCmd),
		RunCmd:     pick(config.RunCmdOverride, framework.RunCmd),
	}
}

func pick(override, base string) string {
	if override != "" {
		return override
	}
	return base
}
