package postgres

import (
"context"
"strings"
"testing"
"time"

"github.com/apps-deployer/projects-service/internal/domain/models"
"github.com/jackc/pgx/v5"
"github.com/jackc/pgx/v5/pgconn"
)

// captureExecutor records the last SQL query and arguments passed to it.
type captureExecutor struct {
lastSQL  string
lastArgs []any
rowFn    func(sql string, args []any) pgx.Row
rowsFn   func(sql string, args []any) (pgx.Rows, error)
}

func (c *captureExecutor) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
c.lastSQL = sql
c.lastArgs = args
if c.rowsFn != nil {
return c.rowsFn(sql, args)
}
return &emptyRows{}, nil
}

func (c *captureExecutor) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
c.lastSQL = sql
c.lastArgs = args
return pgconn.CommandTag{}, nil
}

func (c *captureExecutor) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
c.lastSQL = sql
c.lastArgs = args
if c.rowFn != nil {
return c.rowFn(sql, args)
}
return &staticRow{}
}

// staticRow is a pgx.Row that scans fixed test values.
type staticRow struct {
id        string
createdAt time.Time
updatedAt time.Time
}

func (r *staticRow) Scan(dest ...any) error {
if len(dest) < 3 {
return nil
}
if s, ok := dest[0].(*string); ok {
*s = r.id
}
if t, ok := dest[1].(*time.Time); ok {
*t = r.createdAt
}
if t, ok := dest[2].(*time.Time); ok {
*t = r.updatedAt
}
return nil
}

// emptyRows implements pgx.Rows and returns no data.
type emptyRows struct{}

func (e *emptyRows) Close()                                       {}
func (e *emptyRows) Err() error                                   { return nil }
func (e *emptyRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (e *emptyRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (e *emptyRows) Next() bool                                   { return false }
func (e *emptyRows) Scan(dest ...any) error                       { return nil }
func (e *emptyRows) Values() ([]any, error)                       { return nil, nil }
func (e *emptyRows) RawValues() [][]byte                          { return nil }
func (e *emptyRows) Conn() *pgx.Conn                              { return nil }

func newTestRepo(exec QueryExecutor, encryptionKey string) *Repo {
return &Repo{executor: exec, encryptionKey: encryptionKey}
}

func TestSaveProjectVar_QueryUsesPgpEncrypt(t *testing.T) {
exec := &captureExecutor{
rowFn: func(sql string, args []any) pgx.Row {
return &staticRow{id: "test-id", createdAt: time.Now(), updatedAt: time.Now()}
},
}
repo := newTestRepo(exec, "test-secret-key")

_, err := repo.SaveProjectVar(context.Background(), &models.CreateProjectVarParams{
ProjectId: "proj-uuid",
Key:       "MY_VAR",
Value:     "plaintext-value",
})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !strings.Contains(exec.lastSQL, "crypto.pgp_sym_encrypt") {
t.Errorf("expected SQL to contain pgp_sym_encrypt, got:\n%s", exec.lastSQL)
}
foundKey := false
for _, arg := range exec.lastArgs {
if s, ok := arg.(string); ok && s == "test-secret-key" {
foundKey = true
break
}
}
if !foundKey {
t.Errorf("expected encryption key to be passed as query argument, args: %v", exec.lastArgs)
}
}

func TestUpdateProjectVar_QueryUsesPgpEncrypt(t *testing.T) {
exec := &captureExecutor{}
repo := newTestRepo(exec, "test-secret-key")

newVal := "updated-value"
err := repo.UpdateProjectVar(context.Background(), &models.UpdateVarParams{
Id:    "var-uuid",
Value: &newVal,
})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !strings.Contains(exec.lastSQL, "crypto.pgp_sym_encrypt") {
t.Errorf("expected SQL to contain pgp_sym_encrypt, got:\n%s", exec.lastSQL)
}
}

func TestSaveEnvVar_QueryUsesPgpEncrypt(t *testing.T) {
exec := &captureExecutor{
rowFn: func(sql string, args []any) pgx.Row {
return &staticRow{id: "evar-id", createdAt: time.Now(), updatedAt: time.Now()}
},
}
repo := newTestRepo(exec, "test-secret-key")

_, err := repo.SaveEnvVar(context.Background(), &models.CreateEnvVarParams{
EnvId: "env-uuid",
Key:   "SECRET",
Value: "s3cr3t",
})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !strings.Contains(exec.lastSQL, "crypto.pgp_sym_encrypt") {
t.Errorf("expected SQL to contain pgp_sym_encrypt, got:\n%s", exec.lastSQL)
}
}

func TestUpdateEnvVar_QueryUsesPgpEncrypt(t *testing.T) {
exec := &captureExecutor{}
repo := newTestRepo(exec, "test-secret-key")

newVal := "new-secret"
err := repo.UpdateEnvVar(context.Background(), &models.UpdateVarParams{
Id:    "evar-uuid",
Value: &newVal,
})
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !strings.Contains(exec.lastSQL, "crypto.pgp_sym_encrypt") {
t.Errorf("expected SQL to contain pgp_sym_encrypt, got:\n%s", exec.lastSQL)
}
}

func TestResolvedVars_QueryUsesPgpDecrypt(t *testing.T) {
exec := &captureExecutor{} // returns emptyRows → no rows
repo := newTestRepo(exec, "test-secret-key")

_, err := repo.ResolvedVars(context.Background(), "env-uuid")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if !strings.Contains(exec.lastSQL, "crypto.pgp_sym_decrypt") {
t.Errorf("expected SQL to contain pgp_sym_decrypt, got:\n%s", exec.lastSQL)
}
foundKey := false
for _, arg := range exec.lastArgs {
if s, ok := arg.(string); ok && s == "test-secret-key" {
foundKey = true
break
}
}
if !foundKey {
t.Errorf("expected encryption key to be passed as query argument, args: %v", exec.lastArgs)
}
}
