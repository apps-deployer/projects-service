package frameworks

import (
	"github.com/apps-deployer/projects-service/internal/domain/models"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func frameworkToProto(f *models.Framework) *projectsv1.FrameworkResponse {
	return projectsv1.FrameworkResponse_builder{
		Id:         &f.Id,
		Name:       &f.Name,
		RootDir:    &f.RootDir,
		OutputDir:  &f.OutputDir,
		BaseImage:  &f.BaseImage,
		InstallCmd: &f.InstallCmd,
		BuildCmd:   &f.BuildCmd,
		RunCmd:     &f.RunCmd,
		CreatedAt:  timestamppb.New(f.CreatedAt),
		UpdatedAt:  timestamppb.New(f.UpdatedAt),
	}.Build()
}

func frameworksToProto(frameworks []*models.Framework) *projectsv1.ListFrameworksResponse {
	frameworksResponse := make([]*projectsv1.FrameworkResponse, len(frameworks))
	for i, f := range frameworks {
		frameworksResponse[i] = frameworkToProto(f)
	}
	return projectsv1.ListFrameworksResponse_builder{
		Frameworks: frameworksResponse,
	}.Build()
}

func protoToListFrameworksParams(req *projectsv1.ListFrameworksRequest) *models.ListFrameworksParams {
	return &models.ListFrameworksParams{
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}
}

func protoToCreateFrameworkParams(req *projectsv1.CreateFrameworkRequest) *models.CreateFrameworkParams {
	return &models.CreateFrameworkParams{
		Name:       req.GetName(),
		RootDir:    req.GetRootDir(),
		OutputDir:  req.GetOutputDir(),
		BaseImage:  req.GetBaseImage(),
		InstallCmd: req.GetInstallCmd(),
		BuildCmd:   req.GetBuildCmd(),
		RunCmd:     req.GetRunCmd(),
	}
}

func protoToUpdateFrameworkParams(req *projectsv1.UpdateFrameworkRequest) *models.UpdateFrameworkParams {
	framework := &models.UpdateFrameworkParams{Id: req.GetId()}
	if req.HasName() {
		name := req.GetName()
		framework.Name = &name
	}
	if req.HasRootDir() {
		rootDir := req.GetRootDir()
		framework.RootDir = &rootDir
	}
	if req.HasOutputDir() {
		outputDir := req.GetOutputDir()
		framework.OutputDir = &outputDir
	}
	if req.HasBaseImage() {
		baseImage := req.GetBaseImage()
		framework.BaseImage = &baseImage
	}
	if req.HasInstallCmd() {
		installCmd := req.GetInstallCmd()
		framework.InstallCmd = &installCmd
	}
	if req.HasBuildCmd() {
		buildCmd := req.GetBuildCmd()
		framework.BuildCmd = &buildCmd
	}
	if req.HasRunCmd() {
		runCmd := req.GetRunCmd()
		framework.RunCmd = &runCmd
	}
	return framework
}
