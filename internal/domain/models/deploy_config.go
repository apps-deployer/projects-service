package models

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
