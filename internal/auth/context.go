package auth

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "user_id"

var (
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrPermissionDenied = errors.New("permission denied")
)

// WithUserID returns a new context with the user ID stored.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts the authenticated user ID from the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

// MustUserID extracts the user ID from the context or returns ErrUnauthenticated.
func MustUserID(ctx context.Context) (string, error) {
	id, ok := UserIDFromContext(ctx)
	if !ok || id == "" {
		return "", ErrUnauthenticated
	}
	return id, nil
}

// CheckOwnership verifies that the authenticated user matches the resource owner.
func CheckOwnership(ctx context.Context, ownerID string) error {
	userID, err := MustUserID(ctx)
	if err != nil {
		return err
	}
	if userID != ownerID {
		return ErrPermissionDenied
	}
	return nil
}
