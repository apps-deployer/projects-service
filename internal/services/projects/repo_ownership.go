package projects

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/services"
)

func canonicalGitHubHTTPSRepoURL(repoURL string) (string, error) {
	raw := strings.TrimSpace(repoURL)
	if raw == "" {
		return "", fmt.Errorf("repo URL is empty")
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return "", fmt.Errorf("repo URL must be an HTTPS GitHub clone URL")
	}
	if strings.ToLower(parsed.Hostname()) != "github.com" {
		return "", fmt.Errorf("repo URL host must be github.com")
	}
	if parsed.Port() != "" {
		return "", fmt.Errorf("repo URL must not include a port")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("repo URL must not contain credentials, query, or fragment")
	}
	if strings.HasSuffix(parsed.Path, "/") {
		return "", fmt.Errorf("repo URL must not end with a slash")
	}

	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("repo URL must be in https://github.com/owner/repo.git format")
	}
	if !strings.HasSuffix(parts[1], ".git") || strings.TrimSuffix(parts[1], ".git") == "" {
		return "", fmt.Errorf("repo URL must end with .git")
	}

	return fmt.Sprintf("https://github.com/%s/%s", parts[0], parts[1]), nil
}

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

	canonicalURL, err := canonicalGitHubHTTPSRepoURL(repoURL)
	if err != nil {
		return fmt.Errorf("%w: %v", services.ErrInvalidArgument, err)
	}
	owner, err := githubRepoOwner(canonicalURL)
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
