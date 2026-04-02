package vars_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	"github.com/apps-deployer/projects-service/internal/services"
	"github.com/apps-deployer/projects-service/internal/services/vars"
)

// --- mock storage / repo factory ---

type mockStorage struct {
	factory *mockRepoFactory
}

func (m *mockStorage) Repos() services.RepoFactory { return m.factory }
func (m *mockStorage) WithinTx(ctx context.Context, fn func(services.RepoFactory) error) error {
	return fn(m.factory)
}
func (m *mockStorage) Stop() {}

type mockRepoFactory struct {
	pv services.ProjectVarRepository
	ev services.EnvVarRepository
	rv services.ResolvedVarsRepository
}

func (m *mockRepoFactory) Projects() services.ProjectRepository               { return nil }
func (m *mockRepoFactory) Frameworks() services.FrameworkRepository           { return nil }
func (m *mockRepoFactory) DeployConfigs() services.DeployConfigRepository     { return nil }
func (m *mockRepoFactory) Envs() services.EnvRepository                       { return nil }
func (m *mockRepoFactory) ProjectVars() services.ProjectVarRepository         { return m.pv }
func (m *mockRepoFactory) EnvVars() services.EnvVarRepository                 { return m.ev }
func (m *mockRepoFactory) ResolvedVars() services.ResolvedVarsRepository      { return m.rv }

// --- mock project var repository ---

type mockProjectVarRepo struct {
	saveResp *models.SaveVarResponse
	saveErr  error
	listResp []*models.Var
	listErr  error
	updateErr error
	deleteErr error
}

func (m *mockProjectVarRepo) ListProjectVars(_ context.Context, _ *models.ListProjectVarsParams) ([]*models.Var, error) {
	return m.listResp, m.listErr
}
func (m *mockProjectVarRepo) SaveProjectVar(_ context.Context, _ *models.CreateProjectVarParams) (*models.SaveVarResponse, error) {
	return m.saveResp, m.saveErr
}
func (m *mockProjectVarRepo) UpdateProjectVar(_ context.Context, _ *models.UpdateVarParams) error {
	return m.updateErr
}
func (m *mockProjectVarRepo) DeleteProjectVar(_ context.Context, _ string) error {
	return m.deleteErr
}

// --- mock env var repository ---

type mockEnvVarRepo struct {
	saveResp  *models.SaveVarResponse
	saveErr   error
	listResp  []*models.Var
	listErr   error
	updateErr error
	deleteErr error
}

func (m *mockEnvVarRepo) ListEnvVars(_ context.Context, _ *models.ListEnvVarsParams) ([]*models.Var, error) {
	return m.listResp, m.listErr
}
func (m *mockEnvVarRepo) SaveEnvVar(_ context.Context, _ *models.CreateEnvVarParams) (*models.SaveVarResponse, error) {
	return m.saveResp, m.saveErr
}
func (m *mockEnvVarRepo) UpdateEnvVar(_ context.Context, _ *models.UpdateVarParams) error {
	return m.updateErr
}
func (m *mockEnvVarRepo) DeleteEnvVar(_ context.Context, _ string) error {
	return m.deleteErr
}

// --- mock resolved vars repository ---

type mockResolvedVarsRepo struct {
	resp []*models.ResolvedVar
	err  error
}

func (m *mockResolvedVarsRepo) ResolvedVars(_ context.Context, _ string) ([]*models.ResolvedVar, error) {
	return m.resp, m.err
}

// --- helpers ---

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newTestStorage(pv services.ProjectVarRepository, ev services.EnvVarRepository, rv services.ResolvedVarsRepository) services.Storage {
	return &mockStorage{
		factory: &mockRepoFactory{pv: pv, ev: ev, rv: rv},
	}
}

// --- tests ---

func TestCreateProjectVar_HappyPath(t *testing.T) {
	now := time.Now()
	pvRepo := &mockProjectVarRepo{
		saveResp: &models.SaveVarResponse{
			Id:        "var-uuid",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	svc := vars.New(newTestLogger(), newTestStorage(pvRepo, nil, nil))

	result, err := svc.CreateProjectVar(context.Background(), &models.CreateProjectVarParams{
		ProjectId: "proj-uuid",
		Key:       "DATABASE_URL",
		Value:     "postgres://localhost/mydb",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Id != "var-uuid" {
		t.Errorf("expected id %q, got %q", "var-uuid", result.Id)
	}
	if result.Key != "DATABASE_URL" {
		t.Errorf("expected key %q, got %q", "DATABASE_URL", result.Key)
	}
}

func TestCreateProjectVar_RepoError(t *testing.T) {
	wantErr := errors.New("db error")
	pvRepo := &mockProjectVarRepo{saveErr: wantErr}
	svc := vars.New(newTestLogger(), newTestStorage(pvRepo, nil, nil))

	_, err := svc.CreateProjectVar(context.Background(), &models.CreateProjectVarParams{
		ProjectId: "proj-uuid",
		Key:       "KEY",
		Value:     "value",
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("expected error %v, got %v", wantErr, err)
	}
}

func TestListProjectVars_HappyPath(t *testing.T) {
	now := time.Now()
	pvRepo := &mockProjectVarRepo{
		listResp: []*models.Var{
			{Id: "v1", Key: "K1", CreatedAt: now, UpdatedAt: now},
			{Id: "v2", Key: "K2", CreatedAt: now, UpdatedAt: now},
		},
	}
	svc := vars.New(newTestLogger(), newTestStorage(pvRepo, nil, nil))

	result, err := svc.ListProjectVars(context.Background(), &models.ListProjectVarsParams{
		ProjectId: "proj-uuid",
		Limit:     10,
		Offset:    0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 vars, got %d", len(result))
	}
}

func TestCreateEnvVar_HappyPath(t *testing.T) {
	now := time.Now()
	evRepo := &mockEnvVarRepo{
		saveResp: &models.SaveVarResponse{
			Id:        "evar-uuid",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	svc := vars.New(newTestLogger(), newTestStorage(nil, evRepo, nil))

	result, err := svc.CreateEnvVar(context.Background(), &models.CreateEnvVarParams{
		EnvId: "env-uuid",
		Key:   "SECRET_KEY",
		Value: "s3cr3t",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Id != "evar-uuid" {
		t.Errorf("expected id %q, got %q", "evar-uuid", result.Id)
	}
	if result.Key != "SECRET_KEY" {
		t.Errorf("expected key %q, got %q", "SECRET_KEY", result.Key)
	}
}

func TestResolveVars_HappyPath(t *testing.T) {
	rvRepo := &mockResolvedVarsRepo{
		resp: []*models.ResolvedVar{
			{Id: "", Key: "DB_URL", Value: "postgres://localhost/db"},
			{Id: "", Key: "SECRET", Value: "topsecret"},
		},
	}
	svc := vars.New(newTestLogger(), newTestStorage(nil, nil, rvRepo))

	result, err := svc.ResolveVars(context.Background(), "env-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 resolved vars, got %d", len(result))
	}
	if result[0].Value != "postgres://localhost/db" {
		t.Errorf("expected value %q, got %q", "postgres://localhost/db", result[0].Value)
	}
}

func TestUpdateProjectVar_HappyPath(t *testing.T) {
	pvRepo := &mockProjectVarRepo{}
	svc := vars.New(newTestLogger(), newTestStorage(pvRepo, nil, nil))

	newVal := "newvalue"
	err := svc.UpdateProjectVar(context.Background(), &models.UpdateVarParams{
		Id:    "var-uuid",
		Value: &newVal,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteProjectVar_HappyPath(t *testing.T) {
	pvRepo := &mockProjectVarRepo{}
	svc := vars.New(newTestLogger(), newTestStorage(pvRepo, nil, nil))

	err := svc.DeleteProjectVar(context.Background(), "var-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
