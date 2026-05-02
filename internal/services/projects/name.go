package projects

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/apps-deployer/projects-service/internal/services"
)

var (
	invalidProjectNameChars = regexp.MustCompile(`[^a-z0-9-]+`)
	projectNameDashes       = regexp.MustCompile(`-+`)
)

func displayProjectName(name string) (string, error) {
	displayName := strings.TrimSpace(name)
	if displayName == "" {
		return "", fmt.Errorf("%w: project name is required", services.ErrInvalidArgument)
	}
	if utf8.RuneCountInString(displayName) > 128 {
		return "", fmt.Errorf("%w: project name must be 128 characters or less", services.ErrInvalidArgument)
	}
	return displayName, nil
}

func projectSlug(name string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	normalized = invalidProjectNameChars.ReplaceAllString(normalized, "-")
	normalized = projectNameDashes.ReplaceAllString(normalized, "-")
	normalized = strings.Trim(normalized, "-")
	if len(normalized) > 50 {
		normalized = strings.Trim(normalized[:50], "-")
	}
	if normalized == "" {
		return "", fmt.Errorf("%w: project name must contain latin letters or digits", services.ErrInvalidArgument)
	}
	return normalized, nil
}
