package auth

import (
	"context"
	"errors"
)

type contextKey string

const (
	userIDKey      contextKey = "user_id"
	githubLoginKey contextKey = "github_login"
)

var (
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrPermissionDenied = errors.New("permission denied")
)

// WithUserID returns a new context with the user ID stored.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithGitHubLogin returns a new context with the GitHub login stored.
func WithGitHubLogin(ctx context.Context, login string) context.Context {
	return context.WithValue(ctx, githubLoginKey, login)
}

// WithUser returns a new context with authenticated user fields stored.
func WithUser(ctx context.Context, userID string, githubLogin string) context.Context {
	ctx = WithUserID(ctx, userID)
	if githubLogin != "" {
		ctx = WithGitHubLogin(ctx, githubLogin)
	}
	return ctx
}

// UserIDFromContext extracts the authenticated user ID from the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

// GitHubLoginFromContext extracts the authenticated GitHub login from the context.
func GitHubLoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(githubLoginKey).(string)
	return login, ok
}

// MustUserID extracts the user ID from the context or returns ErrUnauthenticated.
func MustUserID(ctx context.Context) (string, error) {
	id, ok := UserIDFromContext(ctx)
	if !ok || id == "" {
		return "", ErrUnauthenticated
	}
	return id, nil
}

// IsServiceAccount returns true if the user ID belongs to a service account.
func IsServiceAccount(userID string) bool {
	return len(userID) > 8 && userID[:8] == "service:"
}

// CheckOwnership verifies that the authenticated user matches the resource owner.
// Service accounts bypass ownership checks.
func CheckOwnership(ctx context.Context, ownerID string) error {
	userID, err := MustUserID(ctx)
	if err != nil {
		return err
	}
	if IsServiceAccount(userID) {
		return nil
	}
	if userID != ownerID {
		return ErrPermissionDenied
	}
	return nil
}
