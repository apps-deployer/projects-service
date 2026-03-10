package frameworks

import (
	"context"

	"github.com/apps-deployer/projects-service/internal/domain/models"
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
	Get(ctx context.Context, id string) (*models.Framework, error)
	List(ctx context.Context, args *models.ListFrameworksParams) ([]*models.Framework, error)
	Create(ctx context.Context, args *models.CreateFrameworkParams) (*models.Framework, error)
	Update(ctx context.Context, args *models.UpdateFrameworkParams) error
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
	return newFrameworkResponse(framework), nil
}

func (s *frameworksServer) ListFrameworks(
	ctx context.Context,
	req *projectsv1.ListFrameworksRequest,
) (*projectsv1.ListFrameworksResponse, error) {
	frameworks, err := s.frameworks.List(
		ctx,
		newListFrameworksParams(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list frameworks: %v", err)
	}
	frameworksResponse := make([]*projectsv1.FrameworkResponse, len(frameworks))
	for i, f := range frameworks {
		frameworksResponse[i] = newFrameworkResponse(f)
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
	framework, err := s.frameworks.Create(ctx, newCreateFrameworkParams(req))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create framework: %v", err)
	}
	return newFrameworkResponse(framework), nil
}

func (s *frameworksServer) UpdateFramework(
	ctx context.Context,
	req *projectsv1.UpdateFrameworkRequest,
) (*emptypb.Empty, error) {
	if !req.HasId() {
		return nil, status.Error(codes.InvalidArgument, "framework ID is required")
	}
	args := newUpdateFrameworkParams(req)
	if err := s.frameworks.Update(ctx, args); err != nil {
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

func newFrameworkResponse(f *models.Framework) *projectsv1.FrameworkResponse {
	return projectsv1.FrameworkResponse_builder{
		Id:         &f.Id,
		Name:       &f.Name,
		RootDir:    &f.RootDir,
		OutputDir:  &f.OutputDir,
		BaseImage:  &f.BaseImage,
		InstallCmd: &f.InstallCmd,
		BuildCmd:   &f.BuildCmd,
		RunCmd:     &f.RunCmd,
	}.Build()
}

func newListFrameworksParams(req *projectsv1.ListFrameworksRequest) *models.ListFrameworksParams {
	return &models.ListFrameworksParams{
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}
}

func newCreateFrameworkParams(req *projectsv1.CreateFrameworkRequest) *models.CreateFrameworkParams {
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

func newUpdateFrameworkParams(req *projectsv1.UpdateFrameworkRequest) *models.UpdateFrameworkParams {
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
