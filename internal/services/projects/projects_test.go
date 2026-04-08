package projects_test

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
	"github.com/apps-deployer/projects-service/internal/services/projects"
)

const testUserID = "owner-uuid"
const otherUserID = "other-uuid"

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
}

func (m *mockRepoFactory) Projects() services.ProjectRepository           { return m.projects }
func (m *mockRepoFactory) Frameworks() services.FrameworkRepository       { return nil }
func (m *mockRepoFactory) DeployConfigs() services.DeployConfigRepository { return m.deployConfigs }
func (m *mockRepoFactory) Envs() services.EnvRepository                   { return nil }
func (m *mockRepoFactory) ProjectVars() services.ProjectVarRepository     { return nil }
func (m *mockRepoFactory) EnvVars() services.EnvVarRepository             { return nil }
func (m *mockRepoFactory) ResolvedVars() services.ResolvedVarsRepository  { return nil }

type mockProjectRepo struct {
	project    *models.Project
	projectErr error
	listResp   []*models.Project
	listErr    error
	saveResp   *models.SaveProjectResponse
	saveErr    error
	updateErr  error
	deleteErr  error
}

func (m *mockProjectRepo) Project(_ context.Context, _ string) (*models.Project, error) {
	return m.project, m.projectErr
}
func (m *mockProjectRepo) ListProjects(_ context.Context, _ *models.ListProjectsParams) ([]*models.Project, error) {
	return m.listResp, m.listErr
}
func (m *mockProjectRepo) SaveProject(_ context.Context, _ *models.SaveProjectParams) (*models.SaveProjectResponse, error) {
	return m.saveResp, m.saveErr
}
func (m *mockProjectRepo) UpdateProject(_ context.Context, _ *models.UpdateProjectParams) error {
	return m.updateErr
}
func (m *mockProjectRepo) DeleteProject(_ context.Context, _ string) error {
	return m.deleteErr
}

type mockDeployConfigRepo struct {
	saveResp *models.SaveDeployConfigResponse
	saveErr  error
}

func (m *mockDeployConfigRepo) DeployConfig(_ context.Context, _ string) (*models.DeployConfig, error) {
	return nil, nil
}
func (m *mockDeployConfigRepo) UpdateDeployConfig(_ context.Context, _ *models.UpdateDeployConfigParams) error {
	return nil
}
func (m *mockDeployConfigRepo) SaveDeployConfig(_ context.Context, _ *models.SaveDeployConfigParams) (*models.SaveDeployConfigResponse, error) {
	return m.saveResp, m.saveErr
}
func (m *mockDeployConfigRepo) DeleteDeployConfig(_ context.Context, _ string) error { return nil }
func (m *mockDeployConfigRepo) ProjectOwnerByDeployConfigID(_ context.Context, _ string) (string, error) {
	return "", nil
}

// --- helpers ---

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func authedCtx(userID string) context.Context {
	return auth.WithUserID(context.Background(), userID)
}

func testProject() *models.Project {
	return &models.Project{
		Id:      "proj-1",
		Name:    "test",
		RepoUrl: "https://github.com/test/repo",
		OwnerId: testUserID,
	}
}

// --- tests ---

func TestGet_HappyPath(t *testing.T) {
	proj := testProject()
	repo := &mockProjectRepo{project: proj}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	result, err := svc.Get(authedCtx(testUserID), "proj-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Id != "proj-1" {
		t.Errorf("expected id %q, got %q", "proj-1", result.Id)
	}
}

func TestGet_PermissionDenied(t *testing.T) {
	proj := testProject()
	repo := &mockProjectRepo{project: proj}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	_, err := svc.Get(authedCtx(otherUserID), "proj-1")
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}

func TestGet_Unauthenticated(t *testing.T) {
	proj := testProject()
	repo := &mockProjectRepo{project: proj}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	_, err := svc.Get(context.Background(), "proj-1")
	if !errors.Is(err, auth.ErrUnauthenticated) {
		t.Errorf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestList_HappyPath(t *testing.T) {
	repo := &mockProjectRepo{listResp: []*models.Project{testProject()}}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	result, err := svc.List(authedCtx(testUserID), &models.ListProjectsParams{
		OwnerId: testUserID, Limit: 10, Offset: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 project, got %d", len(result))
	}
}

func TestCreate_HappyPath(t *testing.T) {
	now := time.Now()
	repo := &mockProjectRepo{
		saveResp: &models.SaveProjectResponse{Id: "new-proj", CreatedAt: now, UpdatedAt: now},
	}
	dcRepo := &mockDeployConfigRepo{
		saveResp: &models.SaveDeployConfigResponse{Id: "dc-1", CreatedAt: now, UpdatedAt: now},
	}
	svc := projects.New(newLogger(), &mockStorage{
		factory: &mockRepoFactory{projects: repo, deployConfigs: dcRepo},
	})

	result, err := svc.Create(authedCtx(testUserID), &models.CreateProjectParams{
		Name:        "new-project",
		RepoUrl:     "https://github.com/test/new",
		OwnerId:     testUserID,
		FrameworkId: "fw-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Id != "new-proj" {
		t.Errorf("expected id %q, got %q", "new-proj", result.Id)
	}
}

func TestUpdate_PermissionDenied(t *testing.T) {
	proj := testProject()
	repo := &mockProjectRepo{project: proj}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	name := "new-name"
	err := svc.Update(authedCtx(otherUserID), &models.UpdateProjectParams{
		Id:   "proj-1",
		Name: &name,
	})
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}

func TestDelete_HappyPath(t *testing.T) {
	proj := testProject()
	repo := &mockProjectRepo{project: proj}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	err := svc.Delete(authedCtx(testUserID), "proj-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_PermissionDenied(t *testing.T) {
	proj := testProject()
	repo := &mockProjectRepo{project: proj}
	svc := projects.New(newLogger(), &mockStorage{factory: &mockRepoFactory{projects: repo}})

	err := svc.Delete(authedCtx(otherUserID), "proj-1")
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}
