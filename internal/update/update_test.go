package update

import (
	"os"
	"path/filepath"
	"testing"
)

// makeProject creates a fake project directory with the given package.json content.
func makeProject(t *testing.T, root, name, content string) {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRun_DryRun_NoExec(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "app", `{"devDependencies":{"svelte":"^5.28.1"}}`)

	// dry-run should complete without error and without modifying anything.
	err := Run(Options{
		TargetDir: root,
		DryRun:    true,
		SkipGit:   true, // no git repo in temp dir
	})
	if err != nil {
		t.Errorf("Run(DryRun=true) unexpected error: %v", err)
	}
}

func TestRun_SkipsLegacySvelte(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "old", `{"devDependencies":{"svelte":"^4.2.8"}}`)

	// No exec calls should happen; legacy project should just be skipped.
	err := Run(Options{
		TargetDir: root,
		DryRun:    true,
		SkipGit:   true,
	})
	// Should not error — skipping is not a failure.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_SkipsNoSvelte(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "react-app", `{"dependencies":{"react":"18.0.0"}}`)

	err := Run(Options{
		TargetDir: root,
		DryRun:    true,
		SkipGit:   true,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_ExcludeDir(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "app", `{"devDependencies":{"svelte":"^5.28.1"}}`)
	makeProject(t, root, "api", `{"devDependencies":{"svelte":"^5.0.0"}}`)

	// Exclude "api" — only "app" should be processed.
	err := Run(Options{
		TargetDir: root,
		DryRun:    true,
		SkipGit:   true,
		Exclude:   []string{"api"},
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_FromIndex(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "aaa", `{"devDependencies":{"svelte":"^5.0.0"}}`)
	makeProject(t, root, "bbb", `{"devDependencies":{"svelte":"^5.1.0"}}`)

	// From=1 should skip "aaa" and only process "bbb".
	err := Run(Options{
		TargetDir: root,
		DryRun:    true,
		SkipGit:   true,
		From:      1,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_FromOutOfRange(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "app", `{"devDependencies":{"svelte":"^5.28.1"}}`)

	err := Run(Options{
		TargetDir: root,
		DryRun:    true,
		From:      99,
	})
	if err == nil {
		t.Error("expected error for out-of-range --from")
	}
}
