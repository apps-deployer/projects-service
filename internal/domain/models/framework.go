package models

import "time"

type Framework struct {
	Id         string
	Name       string
	RootDir    string
	OutputDir  string
	BaseImage  string
	InstallCmd string
	BuildCmd   string
	RunCmd     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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

type SaveFrameworkResponse struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
}
