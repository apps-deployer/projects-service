package models

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

type ListFrameworksParams struct {
	Limit  int64
	Offset int64
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

type SaveFrameworkParams = CreateFrameworkParams

type SaveFrameworkResponse struct {
	Id string
}
