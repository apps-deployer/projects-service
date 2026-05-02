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
	AppPortOverride    int32
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
	AppPortOverride    *int32
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
	AppPort    int32
}

type SaveDeployConfigParams struct {
	ProjectId   string
	FrameworkId string
}

type SaveDeployConfigResponse struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
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
		AppPort:    pickInt32(config.AppPortOverride, framework.AppPort),
	}
}

func pick(override, base string) string {
	if override != "" {
		return override
	}
	return base
}

func pickInt32(override, base int32) int32 {
	if override != 0 {
		return override
	}
	return base
}
