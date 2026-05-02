package projects

import (
	"context"
	"errors"
	"testing"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/services"
)

func TestGitHubRepoOwner(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		want    string
	}{
		{name: "https", repoURL: "https://github.com/octocat/hello-world", want: "octocat"},
		{name: "https git suffix", repoURL: "https://github.com/octocat/hello-world.git", want: "octocat"},
		{name: "ssh scp", repoURL: "git@github.com:octocat/hello-world.git", want: "octocat"},
		{name: "ssh url", repoURL: "ssh://git@github.com/octocat/hello-world.git", want: "octocat"},
		{name: "host path", repoURL: "github.com/octocat/hello-world", want: "octocat"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := githubRepoOwner(tt.repoURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected owner %q, got %q", tt.want, got)
			}
		})
	}
}

func TestValidateRepoOwnership(t *testing.T) {
	ctx := auth.WithUser(context.Background(), "user-1", "octocat")
	if err := validateRepoOwnership(ctx, "https://github.com/octocat/hello-world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := validateRepoOwnership(ctx, "https://github.com/other/hello-world")
	if !errors.Is(err, auth.ErrPermissionDenied) {
		t.Fatalf("expected ErrPermissionDenied, got %v", err)
	}

	err = validateRepoOwnership(ctx, "https://gitlab.com/octocat/hello-world")
	if !errors.Is(err, services.ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestValidateRepoOwnership_ServiceAccount(t *testing.T) {
	ctx := auth.WithUserID(context.Background(), "service:webhook")
	if err := validateRepoOwnership(ctx, "https://gitlab.com/not/github"); err != nil {
		t.Fatalf("expected service account bypass, got %v", err)
	}
}
