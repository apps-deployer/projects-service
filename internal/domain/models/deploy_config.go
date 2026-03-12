package models

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

type SaveDeployConfigParams struct {
	ProjectId   string
	FrameworkId string
}
