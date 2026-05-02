package projects

import "testing"

func TestDisplayProjectName(t *testing.T) {
	got, err := displayProjectName(" Python Hello App ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Python Hello App" {
		t.Fatalf("expected display name with spaces preserved, got %q", got)
	}
}

func TestProjectSlug(t *testing.T) {
	got, err := projectSlug("Python Hello App")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "python-hello-app" {
		t.Fatalf("expected python-hello-app, got %q", got)
	}
}

func TestProjectSlugRejectsEmptySlug(t *testing.T) {
	if _, err := projectSlug("!!!"); err == nil {
		t.Fatalf("expected error")
	}
}
