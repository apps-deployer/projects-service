package grpcutil

import (
	"errors"

	"github.com/apps-deployer/projects-service/internal/auth"
	"github.com/apps-deployer/projects-service/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapError converts domain/storage errors to appropriate gRPC status errors.
func MapError(err error, action string) error {
	switch {
	case errors.Is(err, auth.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, "user not authenticated")
	case errors.Is(err, auth.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, "permission denied")
	case errors.Is(err, storage.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, storage.ErrAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "already exists")
	case errors.Is(err, storage.ErrConflict):
		return status.Errorf(codes.FailedPrecondition, "conflict")
	default:
		return status.Errorf(codes.Internal, "failed to %s: %v", action, err)
	}
}
