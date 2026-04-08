package auth_test

import (
	"context"
	"testing"

	"github.com/apps-deployer/projects-service/internal/auth"
)

func TestWithUserID_And_UserIDFromContext(t *testing.T) {
	ctx := auth.WithUserID(context.Background(), "user-123")
	id, ok := auth.UserIDFromContext(ctx)
	if !ok {
		t.Fatal("expected user ID in context")
	}
	if id != "user-123" {
		t.Errorf("expected %q, got %q", "user-123", id)
	}
}

func TestUserIDFromContext_Missing(t *testing.T) {
	_, ok := auth.UserIDFromContext(context.Background())
	if ok {
		t.Error("expected no user ID in empty context")
	}
}

func TestMustUserID_Success(t *testing.T) {
	ctx := auth.WithUserID(context.Background(), "user-123")
	id, err := auth.MustUserID(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "user-123" {
		t.Errorf("expected %q, got %q", "user-123", id)
	}
}

func TestMustUserID_Missing(t *testing.T) {
	_, err := auth.MustUserID(context.Background())
	if err != auth.ErrUnauthenticated {
		t.Errorf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestCheckOwnership_Match(t *testing.T) {
	ctx := auth.WithUserID(context.Background(), "user-123")
	if err := auth.CheckOwnership(ctx, "user-123"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCheckOwnership_Mismatch(t *testing.T) {
	ctx := auth.WithUserID(context.Background(), "user-123")
	err := auth.CheckOwnership(ctx, "user-456")
	if err != auth.ErrPermissionDenied {
		t.Errorf("expected ErrPermissionDenied, got %v", err)
	}
}

func TestCheckOwnership_NoUser(t *testing.T) {
	err := auth.CheckOwnership(context.Background(), "user-123")
	if err != auth.ErrUnauthenticated {
		t.Errorf("expected ErrUnauthenticated, got %v", err)
	}
}

func TestCheckOwnership_ServiceAccount(t *testing.T) {
	ctx := auth.WithUserID(context.Background(), "service:deployments")
	// Service accounts bypass ownership checks
	if err := auth.CheckOwnership(ctx, "any-owner"); err != nil {
		t.Errorf("expected nil for service account, got %v", err)
	}
}

func TestIsServiceAccount(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"service:deployments", true},
		{"service:auth", true},
		{"user-123", false},
		{"", false},
		{"service", false},
		{"service:", false},
	}
	for _, tt := range tests {
		if got := auth.IsServiceAccount(tt.id); got != tt.want {
			t.Errorf("IsServiceAccount(%q) = %v, want %v", tt.id, got, tt.want)
		}
	}
}
