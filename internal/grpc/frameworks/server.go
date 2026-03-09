package frameworks

import (
	"context"

	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type frameworksServer struct {
	projectsv1.UnimplementedFrameworkServiceServer
	frameworks FrameworksService
}

type FrameworksService interface {
	Get(ctx context.Context, id string) (*Framework, error)
	List(ctx context.Context, limit int64, offset int64) ([]*Framework, error)
	Create(ctx context.Context, project *CreateFrameworkParams) (*Framework, error)
	Update(ctx context.Context, project *UpdateFrameworkParams) error
	Delete(ctx context.Context, id string) error
}

func Register(
	grpcServer *grpc.Server, frameworks FrameworksService) {
	projectsv1.RegisterFrameworkServiceServer(
		grpcServer,
		&frameworksServer{frameworks: frameworks},
	)
}

func (s *frameworksServer) GetFramework(
	ctx context.Context,
	req *projectsv1.GetFrameworkRequest,
) (*projectsv1.FrameworkResponse, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	framework, err := s.frameworks.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get framework: %v", err)
	}
	if framework == nil {
		return nil, status.Error(codes.NotFound, "framework not found")
	}
	return framework.ToProto(), nil
}

func (s *frameworksServer) ListFrameworks(
	ctx context.Context,
	req *projectsv1.ListFrameworksRequest,
) (*projectsv1.ListFrameworksResponse, error) {
	frameworks, err := s.frameworks.List(
		ctx,
		req.GetLimit(),
		req.GetOffset(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list frameworks: %v", err)
	}
	frameworksResponse := make([]*projectsv1.FrameworkResponse, len(frameworks))
	for i, project := range frameworks {
		frameworksResponse[i] = project.ToProto()
	}
	return projectsv1.ListFrameworksResponse_builder{
		Frameworks: frameworksResponse,
	}.Build(), nil
}

func (s *frameworksServer) CreateFramework(
	ctx context.Context,
	req *projectsv1.CreateFrameworkRequest,
) (*projectsv1.FrameworkResponse, error) {
	if !req.HasName() {
		return nil, status.Error(codes.InvalidArgument, "framework name is required")
	}
	if !req.HasBaseImage() {
		return nil, status.Error(codes.InvalidArgument, "framework base image is required")
	}
	if !req.HasRunCmd() {
		return nil, status.Error(codes.InvalidArgument, "framework run command is required")
	}
	framework, err := s.frameworks.Create(ctx, &CreateFrameworkParams{
		Name:       req.GetName(),
		RootDir:    req.GetRootDir(),
		OutputDir:  req.GetOutputDir(),
		BaseImage:  req.GetBaseImage(),
		InstallCmd: req.GetInstallCmd(),
		BuildCmd:   req.GetBuildCmd(),
		RunCmd:     req.GetRunCmd(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create framework: %v", err)
	}
	return framework.ToProto(), nil
}

func (s *frameworksServer) UpdateFramework(
	ctx context.Context,
	req *projectsv1.UpdateFrameworkRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	id := req.GetId()
	framework := &UpdateFrameworkParams{Id: id}
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
	if err := s.frameworks.Update(ctx, framework); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update framework: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *frameworksServer) DeleteFramework(
	ctx context.Context,
	req *projectsv1.GetFrameworkRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	if err := s.frameworks.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete framework: %v", err)
	}
	return &emptypb.Empty{}, nil
}

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

func (p *Framework) ToProto() *projectsv1.FrameworkResponse {
	return projectsv1.FrameworkResponse_builder{
		Id:         &p.Id,
		Name:       &p.Name,
		RootDir:    &p.RootDir,
		OutputDir:  &p.OutputDir,
		BaseImage:  &p.BaseImage,
		InstallCmd: &p.InstallCmd,
		BuildCmd:   &p.BuildCmd,
		RunCmd:     &p.RunCmd,
	}.Build()
}
