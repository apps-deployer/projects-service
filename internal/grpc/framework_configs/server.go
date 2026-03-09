package deployconfigstemplates

import (
	"context"

	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type frameworkConfigsServer struct {
	projectsv1.UnimplementedFrameworkConfigServiceServer
	frameworkConfigs FrameworkConfigs
}

type FrameworkConfigs interface {
	Get(ctx context.Context, id string) (*FrameworkConfigDTO, error)
	List(ctx context.Context, limit int64, offset int64) ([]*FrameworkConfigDTO, error)
	Create(ctx context.Context, project *CreateFrameworkConfigDTO) (*FrameworkConfigDTO, error)
	Update(ctx context.Context, project *UpdateFrameworkConfigDTO) error
	Delete(ctx context.Context, id string) error
}

func Register(
	grpcServer *grpc.Server, frameworkConfigs FrameworkConfigs) {
	projectsv1.RegisterFrameworkConfigServiceServer(
		grpcServer,
		&frameworkConfigsServer{frameworkConfigs: frameworkConfigs},
	)
}

func (s *frameworkConfigsServer) GetFrameworkConfig(
	ctx context.Context,
	req *projectsv1.GetFrameworkConfigRequest,
) (*projectsv1.FrameworkConfigResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework config ID is required")
	}
	frameworkConfig, err := s.frameworkConfigs.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get framework config: %v", err)
	}
	if frameworkConfig == nil {
		return nil, status.Error(codes.NotFound, "framework config not found")
	}
	return frameworkConfig.ToProto(), nil
}

func (s *frameworkConfigsServer) ListFrameworkConfigs(
	ctx context.Context,
	req *projectsv1.ListFrameworkConfigsRequest,
) (*projectsv1.ListFrameworkConfigsResponse, error) {
	frameworkConfigs, err := s.frameworkConfigs.List(
		ctx,
		req.GetLimit(),
		req.GetOffset(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list frameworks configs: %v", err)
	}
	frameworkConfigsResponse := make([]*projectsv1.FrameworkConfigResponse, len(frameworkConfigs))
	for i, project := range frameworkConfigs {
		frameworkConfigsResponse[i] = project.ToProto()
	}
	return projectsv1.ListFrameworkConfigsResponse_builder{
		FrameworkConfigs: frameworkConfigsResponse,
	}.Build(), nil
}

func (s *frameworkConfigsServer) CreateFrameworkConfig(
	ctx context.Context,
	req *projectsv1.CreateFrameworkConfigRequest,
) (*projectsv1.FrameworkConfigResponse, error) {
	if !req.HasName() {
		return nil, status.Error(codes.InvalidArgument, "framework name is required")
	}
	if !req.HasBaseImage() {
		return nil, status.Error(codes.InvalidArgument, "framework base image is required")
	}
	if !req.HasRunCmd() {
		return nil, status.Error(codes.InvalidArgument, "framework run command is required")
	}
	var (
		name       = req.GetName()
		rootDir    = req.GetRootDir()
		outputDir  = req.GetOutputDir()
		baseImage  = req.GetBaseImage()
		installCmd = req.GetInstallCmd()
		buildCmd   = req.GetBuildCmd()
		runCmd     = req.GetRunCmd()
	)
	frameworkConfig, err := s.frameworkConfigs.Create(ctx, &CreateFrameworkConfigDTO{
		Name:       &name,
		RootDir:    &rootDir,
		OutputDir:  &outputDir,
		BaseImage:  &baseImage,
		InstallCmd: &installCmd,
		BuildCmd:   &buildCmd,
		RunCmd:     &runCmd,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create framework config: %v", err)
	}
	return frameworkConfig.ToProto(), nil
}

func (s *frameworkConfigsServer) UpdateFrameworkConfig(
	ctx context.Context,
	req *projectsv1.UpdateFrameworkConfigRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework config ID is required")
	}
	id := req.GetId()
	frameworkConfig := &UpdateFrameworkConfigDTO{Id: &id}
	if req.HasName() {
		name := req.GetName()
		frameworkConfig.Name = &name
	}
	if req.HasRootDir() {
		rootDir := req.GetRootDir()
		frameworkConfig.RootDir = &rootDir
	}
	if req.HasOutputDir() {
		outputDir := req.GetOutputDir()
		frameworkConfig.OutputDir = &outputDir
	}
	if req.HasBaseImage() {
		baseImage := req.GetBaseImage()
		frameworkConfig.BaseImage = &baseImage
	}
	if req.HasInstallCmd() {
		installCmd := req.GetInstallCmd()
		frameworkConfig.InstallCmd = &installCmd
	}
	if req.HasBuildCmd() {
		buildCmd := req.GetBuildCmd()
		frameworkConfig.BuildCmd = &buildCmd
	}
	if req.HasRunCmd() {
		runCmd := req.GetRunCmd()
		frameworkConfig.RunCmd = &runCmd
	}
	if err := s.frameworkConfigs.Update(ctx, frameworkConfig); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update framework config: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *frameworkConfigsServer) DeleteFrameworkConfig(
	ctx context.Context,
	req *projectsv1.GetFrameworkConfigRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework config ID is required")
	}
	if err := s.frameworkConfigs.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete framework config: %v", err)
	}
	return &emptypb.Empty{}, nil
}

type FrameworkConfigDTO struct {
	Id         *string
	Name       *string
	RootDir    *string
	OutputDir  *string
	BaseImage  *string
	InstallCmd *string
	BuildCmd   *string
	RunCmd     *string
}

type CreateFrameworkConfigDTO struct {
	Name       *string
	RootDir    *string
	OutputDir  *string
	BaseImage  *string
	InstallCmd *string
	BuildCmd   *string
	RunCmd     *string
}

type UpdateFrameworkConfigDTO struct {
	Id         *string
	Name       *string
	RootDir    *string
	OutputDir  *string
	BaseImage  *string
	InstallCmd *string
	BuildCmd   *string
	RunCmd     *string
}

func (p *FrameworkConfigDTO) ToProto() *projectsv1.FrameworkConfigResponse {
	return projectsv1.FrameworkConfigResponse_builder{
		Id:         p.Id,
		Name:       p.Name,
		RootDir:    p.RootDir,
		OutputDir:  p.OutputDir,
		BaseImage:  p.BaseImage,
		InstallCmd: p.InstallCmd,
		BuildCmd:   p.BuildCmd,
		RunCmd:     p.RunCmd,
	}.Build()
}
