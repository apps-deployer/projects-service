package frameworks_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/services"
	"github.com/apps-deployer/projects-service/internal/services/frameworks"
)

type mockStorage struct {
	factory *mockRepoFactory
}

func (m *mockStorage) Repos() services.RepoFactory { return m.factory }
func (m *mockStorage) WithinTx(ctx context.Context, fn func(services.RepoFactory) error) error {
	return fn(m.factory)
}
func (m *mockStorage) Stop() {}

type mockRepoFactory struct {
	frameworks services.FrameworkRepository
}

func (m *mockRepoFactory) Projects() services.ProjectRepository           { return nil }
func (m *mockRepoFactory) Frameworks() services.FrameworkRepository       { return m.frameworks }
func (m *mockRepoFactory) DeployConfigs() services.DeployConfigRepository { return nil }
func (m *mockRepoFactory) Envs() services.EnvRepository                   { return nil }
func (m *mockRepoFactory) ProjectVars() services.ProjectVarRepository     { return nil }
func (m *mockRepoFactory) EnvVars() services.EnvVarRepository             { return nil }
func (m *mockRepoFactory) ResolvedVars() services.ResolvedVarsRepository  { return nil }

type mockFrameworkRepo struct {
	framework *models.Framework
	listResp  []*models.Framework
	saveResp  *models.SaveFrameworkResponse
	err       error
}

func (m *mockFrameworkRepo) Framework(_ context.Context, _ string) (*models.Framework, error) {
	return m.framework, m.err
}
func (m *mockFrameworkRepo) ListFrameworks(_ context.Context, _ *models.ListFrameworksParams) ([]*models.Framework, error) {
	return m.listResp, m.err
}
func (m *mockFrameworkRepo) SaveFramework(_ context.Context, _ *models.CreateFrameworkParams) (*models.SaveFrameworkResponse, error) {
	return m.saveResp, m.err
}
func (m *mockFrameworkRepo) UpdateFramework(_ context.Context, _ *models.UpdateFrameworkParams) error {
	return m.err
}
func (m *mockFrameworkRepo) DeleteFramework(_ context.Context, _ string) error {
	return m.err
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newService(repo services.FrameworkRepository) *frameworks.Frameworks {
	return frameworks.New(newLogger(), &mockStorage{factory: &mockRepoFactory{frameworks: repo}})
}

func testFramework() *models.Framework {
	now := time.Now()
	return &models.Framework{
		Id: "fw-1", Name: "Python", BaseImage: "python:3.12-alpine",
		InstallCmd: "pip install -r requirements.txt",
		RunCmd:     "python main.py",
		AppPort:    8080,
		CreatedAt:  now, UpdatedAt: now,
	}
}

func TestGet_HappyPath(t *testing.T) {
	svc := newService(&mockFrameworkRepo{framework: testFramework()})
	result, err := svc.Get(context.Background(), "fw-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Id != "fw-1" {
		t.Fatalf("expected fw-1, got %s", result.Id)
	}
}

func TestList_HappyPath(t *testing.T) {
	svc := newService(&mockFrameworkRepo{listResp: []*models.Framework{testFramework()}})
	result, err := svc.List(context.Background(), &models.ListFrameworksParams{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 framework, got %d", len(result))
	}
}

func TestCreate_HappyPath(t *testing.T) {
	now := time.Now()
	svc := newService(&mockFrameworkRepo{saveResp: &models.SaveFrameworkResponse{
		Id: "fw-1", CreatedAt: now, UpdatedAt: now,
	}})
	result, err := svc.Create(context.Background(), &models.CreateFrameworkParams{
		Name: "Node", BaseImage: "node:20-alpine", RunCmd: "node index.js", AppPort: 3000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Node" || result.AppPort != 3000 {
		t.Fatalf("unexpected framework: %+v", result)
	}
}

func TestUpdate_HappyPath(t *testing.T) {
	svc := newService(&mockFrameworkRepo{})
	name := "Updated"
	if err := svc.Update(context.Background(), &models.UpdateFrameworkParams{Id: "fw-1", Name: &name}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_HappyPath(t *testing.T) {
	svc := newService(&mockFrameworkRepo{})
	if err := svc.Delete(context.Background(), "fw-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepoErrorIsReturned(t *testing.T) {
	wantErr := errors.New("db error")
	svc := newService(&mockFrameworkRepo{err: wantErr})
	if _, err := svc.Get(context.Background(), "fw-1"); !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
