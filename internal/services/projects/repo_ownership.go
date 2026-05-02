package projects

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/services"
)

func validateRepoOwnership(ctx context.Context, repoURL string) error {
	userID, err := auth.MustUserID(ctx)
	if err != nil {
		return err
	}
	if auth.IsServiceAccount(userID) {
		return nil
	}

	githubLogin, ok := auth.GitHubLoginFromContext(ctx)
	if !ok || strings.TrimSpace(githubLogin) == "" {
		return auth.ErrUnauthenticated
	}

	owner, err := githubRepoOwner(repoURL)
	if err != nil {
		return fmt.Errorf("%w: %v", services.ErrInvalidArgument, err)
	}
	if !strings.EqualFold(owner, githubLogin) {
		return auth.ErrPermissionDenied
	}
	return nil
}

func githubRepoOwner(repoURL string) (string, error) {
	raw := strings.TrimSpace(repoURL)
	if raw == "" {
		return "", fmt.Errorf("repo URL is empty")
	}

	if strings.HasPrefix(strings.ToLower(raw), "git@github.com:") {
		path := strings.TrimPrefix(raw, "git@github.com:")
		return ownerFromGitHubPath(path)
	}

	parsed, err := url.Parse(raw)
	if err == nil && parsed.Host != "" {
		host := strings.ToLower(parsed.Hostname())
		if host != "github.com" {
			return "", fmt.Errorf("repo URL host must be github.com")
		}
		return ownerFromGitHubPath(parsed.Path)
	}

	if strings.HasPrefix(strings.ToLower(raw), "github.com/") {
		return ownerFromGitHubPath(strings.TrimPrefix(raw, "github.com/"))
	}

	return "", fmt.Errorf("unsupported GitHub repo URL")
}

func ownerFromGitHubPath(path string) (string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || strings.TrimSuffix(parts[1], ".git") == "" {
		return "", fmt.Errorf("repo URL must include owner and repository name")
	}
	return parts[0], nil
}
