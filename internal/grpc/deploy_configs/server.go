package deployconfigs

import (
	"context"

	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type deployConfigsServer struct {
	projectsv1.UnimplementedDeployConfigServiceServer
	deployConfigs DeployConfigs
}

type DeployConfigs interface {
	Generate(ctx context.Context, id string) (*GeneratedDeployConfigDTO, error)
	Get(ctx context.Context, id string) (*DeployConfigDTO, error)
	Update(ctx context.Context, project *UpdateDeployConfigDTO) error
}

func Register(grpcServer *grpc.Server, deployConfigs DeployConfigs) {
	projectsv1.RegisterDeployConfigServiceServer(
		grpcServer,
		&deployConfigsServer{deployConfigs: deployConfigs},
	)
}

func (s *deployConfigsServer) GenerateDeployConfig(
	ctx context.Context,
	req *projectsv1.GetDeployConfigRequest,
) (*projectsv1.GenerateDeployConfigResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	config, err := s.deployConfigs.Generate(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate deploy config: %v", err)
	}
	if config == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return config.ToProto(), nil
}

func (s *deployConfigsServer) GetDeployConfig(
	ctx context.Context,
	req *projectsv1.GetDeployConfigRequest,
) (*projectsv1.DeployConfigResponse, error) {
	if !req.HasProjectId() {
		return nil, status.Error(codes.InvalidArgument, "project ID is required")
	}
	config, err := s.deployConfigs.Get(ctx, req.GetProjectId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate deploy config: %v", err)
	}
	if config == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return config.ToProto(), nil
}

func (s *deployConfigsServer) UpdateDeployConfig(
	ctx context.Context,
	req *projectsv1.UpdateDeployConfigRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "deploy config ID is required")
	}
	id := req.GetId()
	config := &UpdateDeployConfigDTO{Id: &id}
	if req.HasFrameworkId() {
		frameworkId := req.GetFrameworkId()
		config.FrameworkId = &frameworkId
	}
	if req.HasRootDirOverwrite() {
		rootDir := req.GetRootDirOverwrite()
		config.RootDirOverwrite = &rootDir
	}
	if req.HasOutputDirOverwrite() {
		outputDir := req.GetOutputDirOverwrite()
		config.OutputDirOverwrite = &outputDir
	}
	if req.HasBaseImageOverwrite() {
		baseImage := req.GetBaseImageOverwrite()
		config.BaseImageOverwrite = &baseImage
	}
	if req.HasInstallCmdOverwrite() {
		installCmd := req.GetInstallCmdOverwrite()
		config.InstallCmdOverwrite = &installCmd
	}
	if req.HasBuildCmdOverwrite() {
		buildCmd := req.GetBuildCmdOverwrite()
		config.BuildCmdOverwrite = &buildCmd
	}
	if req.HasRunCmdOverwrite() {
		runCmd := req.GetRunCmdOverwrite()
		config.RunCmdOverwrite = &runCmd
	}
	if err := s.deployConfigs.Update(ctx, config); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update deploy config: %v", err)
	}
	return &emptypb.Empty{}, nil
}

type DeployConfigDTO struct {
	Id                  *string
	ProjectId           *string
	FrameworkId         *string
	RootDirOverwrite    *string
	OutputDirOverwrite  *string
	BaseImageOverwrite  *string
	InstallCmdOverwrite *string
	BuildCmdOverwrite   *string
	RunCmdOverwrite     *string
}

type UpdateDeployConfigDTO struct {
	Id                  *string
	FrameworkId         *string
	RootDirOverwrite    *string
	OutputDirOverwrite  *string
	BaseImageOverwrite  *string
	InstallCmdOverwrite *string
	BuildCmdOverwrite   *string
	RunCmdOverwrite     *string
}

type GeneratedDeployConfigDTO struct {
	Id         *string
	ProjectId  *string
	RootDir    *string
	OutputDir  *string
	BaseImage  *string
	InstallCmd *string
	BuildCmd   *string
	RunCmd     *string
}

func (p *DeployConfigDTO) ToProto() *projectsv1.DeployConfigResponse {
	return projectsv1.DeployConfigResponse_builder{
		Id:                  p.Id,
		ProjectId:           p.ProjectId,
		FrameworkId:         p.FrameworkId,
		RootDirOverwrite:    p.RootDirOverwrite,
		OutputDirOverwrite:  p.OutputDirOverwrite,
		BaseImageOverwrite:  p.BaseImageOverwrite,
		InstallCmdOverwrite: p.InstallCmdOverwrite,
		BuildCmdOverwrite:   p.BuildCmdOverwrite,
		RunCmdOverwrite:     p.RunCmdOverwrite,
	}.Build()
}

func (p *GeneratedDeployConfigDTO) ToProto() *projectsv1.GenerateDeployConfigResponse {
	return projectsv1.GenerateDeployConfigResponse_builder{
		Id:         p.Id,
		ProjectId:  p.ProjectId,
		RootDir:    p.RootDir,
		OutputDir:  p.OutputDir,
		BaseImage:  p.BaseImage,
		InstallCmd: p.InstallCmd,
		BuildCmd:   p.BuildCmd,
		RunCmd:     p.RunCmd,
	}.Build()
}
