package varsgrpc_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/apps-deployer/projects-service/internal/domain/models"
	varsgrpc "github.com/apps-deployer/projects-service/internal/grpc/vars"
	projectsv1 "github.com/apps-deployer/protos/gen/go/projects/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// --- mock services ---

type mockProjectVarsSvc struct {
	createResp *models.Var
	createErr  error
	listResp   []*models.Var
	listErr    error
	updateErr  error
	deleteErr  error
}

func (m *mockProjectVarsSvc) ListProjectVars(_ context.Context, _ *models.ListProjectVarsParams) ([]*models.Var, error) {
	return m.listResp, m.listErr
}
func (m *mockProjectVarsSvc) CreateProjectVar(_ context.Context, _ *models.CreateProjectVarParams) (*models.Var, error) {
	return m.createResp, m.createErr
}
func (m *mockProjectVarsSvc) UpdateProjectVar(_ context.Context, _ *models.UpdateVarParams) error {
	return m.updateErr
}
func (m *mockProjectVarsSvc) DeleteProjectVar(_ context.Context, _ string) error {
	return m.deleteErr
}

type mockEnvVarsSvc struct {
	createResp *models.Var
	createErr  error
	listResp   []*models.Var
	listErr    error
	updateErr  error
	deleteErr  error
}

func (m *mockEnvVarsSvc) ListEnvVars(_ context.Context, _ *models.ListEnvVarsParams) ([]*models.Var, error) {
	return m.listResp, m.listErr
}
func (m *mockEnvVarsSvc) CreateEnvVar(_ context.Context, _ *models.CreateEnvVarParams) (*models.Var, error) {
	return m.createResp, m.createErr
}
func (m *mockEnvVarsSvc) UpdateEnvVar(_ context.Context, _ *models.UpdateVarParams) error {
	return m.updateErr
}
func (m *mockEnvVarsSvc) DeleteEnvVar(_ context.Context, _ string) error {
	return m.deleteErr
}

type mockVarsAggSvc struct {
	resp []*models.ResolvedVar
	err  error
}

func (m *mockVarsAggSvc) ResolveVars(_ context.Context, _ string) ([]*models.ResolvedVar, error) {
	return m.resp, m.err
}

// --- test server setup ---

func newTestServer(
	t *testing.T,
	pvSvc varsgrpc.ProjectVarsService,
	evSvc varsgrpc.EnvVarsService,
	aggSvc varsgrpc.VarsAggregationService,
) *grpc.ClientConn {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	varsgrpc.Register(srv, pvSvc, evSvc, aggSvc)

	go func() {
		if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Logf("server error: %v", err)
		}
	}()
	t.Cleanup(srv.GracefulStop)

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

// --- tests: project vars ---

func TestGRPC_CreateProjectVar_HappyPath(t *testing.T) {
	now := time.Now()
	pvSvc := &mockProjectVarsSvc{
		createResp: &models.Var{Id: "var-1", Key: "MY_KEY", CreatedAt: now, UpdatedAt: now},
	}
	conn := newTestServer(t, pvSvc, &mockEnvVarsSvc{}, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	projectId := "proj-1"
	key := "MY_KEY"
	value := "my_value"
	resp, err := client.CreateProjectVar(context.Background(), projectsv1.CreateProjectVarRequest_builder{
		ProjectId: &projectId,
		Key:       &key,
		Value:     &value,
	}.Build())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GetId() != "var-1" {
		t.Errorf("expected id %q, got %q", "var-1", resp.GetId())
	}
	if resp.GetKey() != "MY_KEY" {
		t.Errorf("expected key %q, got %q", "MY_KEY", resp.GetKey())
	}
}

func TestGRPC_CreateProjectVar_MissingProjectId(t *testing.T) {
	conn := newTestServer(t, &mockProjectVarsSvc{}, &mockEnvVarsSvc{}, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	key := "MY_KEY"
	value := "my_value"
	_, err := client.CreateProjectVar(context.Background(), projectsv1.CreateProjectVarRequest_builder{
		Key:   &key,
		Value: &value,
	}.Build())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected %v, got %v", codes.InvalidArgument, st.Code())
	}
}

func TestGRPC_CreateProjectVar_MissingKey(t *testing.T) {
	conn := newTestServer(t, &mockProjectVarsSvc{}, &mockEnvVarsSvc{}, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	projectId := "proj-1"
	value := "my_value"
	_, err := client.CreateProjectVar(context.Background(), projectsv1.CreateProjectVarRequest_builder{
		ProjectId: &projectId,
		Value:     &value,
	}.Build())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected %v, got %v", codes.InvalidArgument, st.Code())
	}
}

func TestGRPC_ListProjectVars_HappyPath(t *testing.T) {
	now := time.Now()
	pvSvc := &mockProjectVarsSvc{
		listResp: []*models.Var{
			{Id: "v1", Key: "K1", CreatedAt: now, UpdatedAt: now},
			{Id: "v2", Key: "K2", CreatedAt: now, UpdatedAt: now},
		},
	}
	conn := newTestServer(t, pvSvc, &mockEnvVarsSvc{}, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	projectId := "proj-1"
	resp, err := client.ListProjectVars(context.Background(), projectsv1.ListProjectVarsRequest_builder{
		ProjectId: &projectId,
	}.Build())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.GetVars()) != 2 {
		t.Errorf("expected 2 vars, got %d", len(resp.GetVars()))
	}
}

func TestGRPC_ListProjectVars_MissingProjectId(t *testing.T) {
	conn := newTestServer(t, &mockProjectVarsSvc{}, &mockEnvVarsSvc{}, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	_, err := client.ListProjectVars(context.Background(), projectsv1.ListProjectVarsRequest_builder{}.Build())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected %v, got %v", codes.InvalidArgument, st.Code())
	}
}

// --- tests: env vars ---

func TestGRPC_CreateEnvVar_HappyPath(t *testing.T) {
	now := time.Now()
	evSvc := &mockEnvVarsSvc{
		createResp: &models.Var{Id: "evar-1", Key: "SECRET", CreatedAt: now, UpdatedAt: now},
	}
	conn := newTestServer(t, &mockProjectVarsSvc{}, evSvc, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	envId := "env-1"
	key := "SECRET"
	value := "s3cr3t"
	resp, err := client.CreateEnvVar(context.Background(), projectsv1.CreateEnvVarRequest_builder{
		EnvId: &envId,
		Key:   &key,
		Value: &value,
	}.Build())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GetId() != "evar-1" {
		t.Errorf("expected id %q, got %q", "evar-1", resp.GetId())
	}
}

// --- tests: resolve vars ---

func TestGRPC_ResolveVars_HappyPath(t *testing.T) {
	aggSvc := &mockVarsAggSvc{
		resp: []*models.ResolvedVar{
			{Key: "DB_URL", Value: "postgres://localhost/db"},
			{Key: "SECRET", Value: "topsecret"},
		},
	}
	conn := newTestServer(t, &mockProjectVarsSvc{}, &mockEnvVarsSvc{}, aggSvc)
	client := projectsv1.NewVarServiceClient(conn)

	envId := "env-1"
	resp, err := client.ResolveVars(context.Background(), projectsv1.ResolveVarsRequest_builder{
		EnvId: &envId,
	}.Build())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.GetVars()) != 2 {
		t.Errorf("expected 2 resolved vars, got %d", len(resp.GetVars()))
	}
}

func TestGRPC_ResolveVars_MissingEnvId(t *testing.T) {
	conn := newTestServer(t, &mockProjectVarsSvc{}, &mockEnvVarsSvc{}, &mockVarsAggSvc{})
	client := projectsv1.NewVarServiceClient(conn)

	_, err := client.ResolveVars(context.Background(), projectsv1.ResolveVarsRequest_builder{}.Build())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected %v, got %v", codes.InvalidArgument, st.Code())
	}
}
