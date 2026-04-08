package deployconfigs_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/services"
	deployconfigs "github.com/apps-deployer/projects-service/internal/services/deploy_configs"
)

const testUserID = "owner-uuid"
const testProjectID = "proj-uuid"

// --- mocks ---

type mockStorage struct {
	factory *mockRepoFactory
}

func (m *mockStorage) Repos() services.RepoFactory { return m.factory }
func (m *mockStorage) WithinTx(ctx context.Context, fn func(services.RepoFactory) error) error {
	return fn(m.factory)
}
func (m *mockStorage) Stop() {}

type mockRepoFactory struct {
	projects      services.ProjectRepository
	deployConfigs services.DeployConfigRepository
	frameworks    services.FrameworkRepository
}

func (m *mockRepoFactory) Projects() services.ProjectRepository           { return m.projects }
func (m *mockRepoFactory) Frameworks() services.FrameworkRepository       { return m.frameworks }
func (m *mockRepoFactory) DeployConfigs() services.DeployConfigRepository { return m.deployConfigs }
func (m *mockRepoFactory) Envs() services.EnvRepository                   { return nil }
func (m *mockRepoFactory) ProjectVars() services.ProjectVarRepository     { return nil }
func (m *mockRepoFactory) EnvVars() services.EnvVarRepository             { return nil }
func (m *mockRepoFactory) ResolvedVars() services.ResolvedVarsRepository  { return nil }

type mockProjectRepo struct {
	project *models.Project
	err     error
}

func (m *mockProjectRepo) Project(_ context.Context, _ string) (*models.Project, error) {
	return m.project, m.err
}
func (m *mockProjectRepo) ListProjects(_ context.Context, _ *models.ListProjectsParams) ([]*models.Project, error) {
	return nil, nil
}
func (m *mockProjectRepo) SaveProject(_ context.Context, _ *models.SaveProjectParams) (*models.SaveProjectResponse, error) {
	return nil, nil
}
func (m *mockProjectRepo) UpdateProject(_ context.Context, _ *models.UpdateProjectParams) error {
	return nil
}
func (m *mockProjectRepo) DeleteProject(_ context.Context, _ string) error { return nil }

type mockDeployConfigRepo struct {
	config   *models.DeployConfig
	configErr error
	ownerID  string
	ownerErr error
	updateErr error
}

func (m *mockDeployConfigRepo) DeployConfig(_ context.Context, _ string) (*models.DeployConfig, error) {
	return m.config, m.configErr
}
func (m *mockDeployConfigRepo) UpdateDeployConfig(_ context.Context, _ *models.UpdateDeployConfigParams) error {
	return m.updateErr
}
func (m *mockDeployConfigRepo) SaveDeployConfig(_ context.Context, _ *models.SaveDeployConfigParams) (*models.SaveDeployConfigResponse, error) {
	return nil, nil
}
func (m *mockDeployConfigRepo) DeleteDeployConfig(_ context.Context, _ string) error { return nil }
func (m *mockDeployConfigRepo) ProjectOwnerByDeployConfigID(_ context.Context, _ string) (string, error) {
	return m.ownerID, m.ownerErr
}

type mockFrameworkRepo struct {
	framework *models.Framework
	err       error
}

func (m *mockFrameworkRepo) Framework(_ context.Context, _ string) (*models.Framework, error) {
	return m.framework, m.err
}
func (m *mockFrameworkRepo) ListFrameworks(_ context.Context, _ *models.ListFrameworksParams) ([]*models.Framework, error) {
	return nil, nil
}
func (m *mockFrameworkRepo) SaveFramework(_ context.Context, _ *models.CreateFrameworkParams) (*models.SaveFrameworkResponse, error) {
	return nil, nil
}
func (m *mockFrameworkRepo) UpdateFramework(_ context.Context, _ *models.UpdateFrameworkParams) error {
	return nil
}
func (m *mockFrameworkRepo) DeleteFramework(_ context.Context, _ string) error { return nil }

// --- helpers ---

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func authedCtx() context.Context {
	return auth.WithUserID(context.Background(), testUserID)
}

func defaultProject() *mockProjectRepo {
	return &mockProjectRepo{
		project: &models.Project{Id: testProjectID, OwnerId: testUserID},
	}
}

func testConfig() *models.DeployConfig {
	now := time.Now()
	return &models.DeployConfig{
		Id: "dc-1", ProjectId: testProjectID, FrameworkId: "fw-1",
		CreatedAt: now, UpdatedAt: now,
	}
}

func testFramework() *models.Framework {
	return &models.Framework{
		Id: "fw-1", Name: "Node.js", BaseImage: "node:20",
		RunCmd: "node index.js",
	}
}

// --- tests ---

func TestGet_HappyPath(t *testing.T) {
	svc := deployconfigs.New(newLogger(), &mockStorage{factory: &mockRepoFactory{
		projects:      defaultProject(),
		deployConfigs: &mockDeployConfigRepo{config: testConfig()},
	}})

	config, err := svc.Get(authedCtx(), testProjectID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Id != "dc-1" {
		t.Errorf("expected id %q, got %q", "dc-1", config.Id)
	}
}

func TestGet_PermissionDenied(t *testing.T) {
	otherProject := &mockProjectRepo{
		project: &models.Project{Id: testProjectID, OwnerId: "other-user"},
	}
	svc := deployconfigs.New(newLogger(), &mockStorage{factory: &mockRepoFactory{
		projects:      otherProject,
		deployConfigs: &mockDeployConfigRepo{config: testConfig()},
	}})

	_, err := svc.Get(authedCtx(), testProjectID)
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}

func TestResolve_HappyPath(t *testing.T) {
	svc := deployconfigs.New(newLogger(), &mockStorage{factory: &mockRepoFactory{
		projects:      defaultProject(),
		deployConfigs: &mockDeployConfigRepo{config: testConfig()},
		frameworks:    &mockFrameworkRepo{framework: testFramework()},
	}})

	config, err := svc.Resolve(authedCtx(), testProjectID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.BaseImage != "node:20" {
		t.Errorf("expected base image %q, got %q", "node:20", config.BaseImage)
	}
}

func TestUpdate_HappyPath(t *testing.T) {
	svc := deployconfigs.New(newLogger(), &mockStorage{factory: &mockRepoFactory{
		projects:      defaultProject(),
		deployConfigs: &mockDeployConfigRepo{ownerID: testUserID},
	}})

	err := svc.Update(authedCtx(), &models.UpdateDeployConfigParams{Id: "dc-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdate_PermissionDenied(t *testing.T) {
	svc := deployconfigs.New(newLogger(), &mockStorage{factory: &mockRepoFactory{
		projects:      defaultProject(),
		deployConfigs: &mockDeployConfigRepo{ownerID: "other-user"},
	}})

	err := svc.Update(authedCtx(), &models.UpdateDeployConfigParams{Id: "dc-1"})
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}
