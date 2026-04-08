package envs_test

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
	"github.com/apps-deployer/projects-service/internal/services/envs"
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
	projects services.ProjectRepository
	envs     services.EnvRepository
}

func (m *mockRepoFactory) Projects() services.ProjectRepository           { return m.projects }
func (m *mockRepoFactory) Frameworks() services.FrameworkRepository       { return nil }
func (m *mockRepoFactory) DeployConfigs() services.DeployConfigRepository { return nil }
func (m *mockRepoFactory) Envs() services.EnvRepository                   { return m.envs }
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

type mockEnvRepo struct {
	env      *models.Env
	envErr   error
	listResp []*models.Env
	listErr  error
	saveResp *models.SaveEnvResponse
	saveErr  error
	updateErr error
	deleteErr error
}

func (m *mockEnvRepo) Env(_ context.Context, _ string) (*models.Env, error) {
	return m.env, m.envErr
}
func (m *mockEnvRepo) EnvByGit(_ context.Context, _ string, _ string) (*models.Env, error) {
	return m.env, m.envErr
}
func (m *mockEnvRepo) ListEnvs(_ context.Context, _ *models.ListEnvsParams) ([]*models.Env, error) {
	return m.listResp, m.listErr
}
func (m *mockEnvRepo) SaveEnv(_ context.Context, _ *models.CreateEnvParams) (*models.SaveEnvResponse, error) {
	return m.saveResp, m.saveErr
}
func (m *mockEnvRepo) UpdateEnv(_ context.Context, _ *models.UpdateEnvParams) error {
	return m.updateErr
}
func (m *mockEnvRepo) DeleteEnv(_ context.Context, _ string) error { return m.deleteErr }

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

func testEnv() *models.Env {
	now := time.Now()
	return &models.Env{
		Id: "env-1", Name: "production", ProjectId: testProjectID,
		TargetBranch: "main", DomainName: "example.com",
		CreatedAt: now, UpdatedAt: now,
	}
}

func newService(projRepo *mockProjectRepo, envRepo *mockEnvRepo) *envs.Envs {
	return envs.New(newLogger(), &mockStorage{
		factory: &mockRepoFactory{projects: projRepo, envs: envRepo},
	})
}

// --- tests ---

func TestGet_HappyPath(t *testing.T) {
	svc := newService(defaultProject(), &mockEnvRepo{env: testEnv()})
	env, err := svc.Get(authedCtx(), "env-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Id != "env-1" {
		t.Errorf("expected id %q, got %q", "env-1", env.Id)
	}
}

func TestGet_PermissionDenied(t *testing.T) {
	otherProject := &mockProjectRepo{
		project: &models.Project{Id: testProjectID, OwnerId: "other-user"},
	}
	svc := newService(otherProject, &mockEnvRepo{env: testEnv()})
	_, err := svc.Get(authedCtx(), "env-1")
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}

func TestGetByGit_HappyPath(t *testing.T) {
	svc := newService(defaultProject(), &mockEnvRepo{env: testEnv()})
	env, err := svc.GetByGit(authedCtx(), &models.GetEnvByGitParams{
		RepoUrl: "https://github.com/test/repo", TargetBranch: "main",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Id != "env-1" {
		t.Errorf("expected id %q, got %q", "env-1", env.Id)
	}
}

func TestList_HappyPath(t *testing.T) {
	svc := newService(defaultProject(), &mockEnvRepo{listResp: []*models.Env{testEnv()}})
	result, err := svc.List(authedCtx(), &models.ListEnvsParams{
		ProjectId: testProjectID, Limit: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 env, got %d", len(result))
	}
}

func TestCreate_HappyPath(t *testing.T) {
	now := time.Now()
	svc := newService(defaultProject(), &mockEnvRepo{
		saveResp: &models.SaveEnvResponse{Id: "new-env", CreatedAt: now, UpdatedAt: now},
	})
	env, err := svc.Create(authedCtx(), &models.CreateEnvParams{
		Name: "staging", ProjectId: testProjectID,
		TargetBranch: "develop", DomainName: "staging.example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Id != "new-env" {
		t.Errorf("expected id %q, got %q", "new-env", env.Id)
	}
}

func TestCreate_PermissionDenied(t *testing.T) {
	otherProject := &mockProjectRepo{
		project: &models.Project{Id: testProjectID, OwnerId: "other-user"},
	}
	svc := newService(otherProject, &mockEnvRepo{})
	_, err := svc.Create(authedCtx(), &models.CreateEnvParams{
		Name: "staging", ProjectId: testProjectID,
		TargetBranch: "develop", DomainName: "staging.example.com",
	})
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}

func TestDelete_HappyPath(t *testing.T) {
	svc := newService(defaultProject(), &mockEnvRepo{env: testEnv()})
	err := svc.Delete(authedCtx(), "env-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_PermissionDenied(t *testing.T) {
	otherProject := &mockProjectRepo{
		project: &models.Project{Id: testProjectID, OwnerId: "other-user"},
	}
	svc := newService(otherProject, &mockEnvRepo{env: testEnv()})
	err := svc.Delete(authedCtx(), "env-1")
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}
