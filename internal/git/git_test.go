package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initRepo initialises a git repo in a temp dir and returns its path.
func initRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not found, skipping integration test")
	}
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	return dir
}

// gitRun runs a git command in dir, failing the test on error.
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
}

// writeFile writes content to path, failing the test on error.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile %s: %v", path, err)
	}
}

func TestIsGitRepo_True(t *testing.T) {
	dir := initRepo(t)
	if !IsGitRepo(dir) {
		t.Error("IsGitRepo() = false inside a git repo")
	}
}

func TestIsGitRepo_False(t *testing.T) {
	dir := t.TempDir()
	if IsGitRepo(dir) {
		t.Error("IsGitRepo() = true outside a git repo")
	}
}

func TestIsClean_EmptyRepo(t *testing.T) {
	dir := initRepo(t)
	clean, err := IsClean(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !clean {
		t.Error("expected clean repo on fresh init")
	}
}

func TestIsClean_WithUnstagedFile(t *testing.T) {
	dir := initRepo(t)
	f := filepath.Join(dir, "file.txt")
	writeFile(t, f, "hello")
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-m", "init")

	writeFile(t, f, "changed")
	clean, err := IsClean(dir)
	if err != nil {
		t.Fatal(err)
	}
	if clean {
		t.Error("expected dirty repo after modifying a tracked file")
	}
}

func TestCurrentBranch(t *testing.T) {
	dir := initRepo(t)
	writeFile(t, filepath.Join(dir, "README.md"), "x")
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-m", "init")

	branch, err := CurrentBranch(dir)
	if err != nil {
		t.Fatal(err)
	}
	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestAdd_DryRun(t *testing.T) {
	dir := initRepo(t)
	if err := Add(dir, true); err != nil {
		t.Errorf("Add(dryRun=true) error: %v", err)
	}
}

func TestCommit_DryRun(t *testing.T) {
	dir := initRepo(t)
	if err := Commit(dir, "test commit", true); err != nil {
		t.Errorf("Commit(dryRun=true) error: %v", err)
	}
}

func TestPush_DryRun(t *testing.T) {
	dir := initRepo(t)
	if err := Push(dir, "main", true); err != nil {
		t.Errorf("Push(dryRun=true) error: %v", err)
	}
}

func TestHasStagedChanges_None(t *testing.T) {
	dir := initRepo(t)
	writeFile(t, filepath.Join(dir, "x.txt"), "x")
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-m", "init")

	staged, err := HasStagedChanges(dir)
	if err != nil {
		t.Fatal(err)
	}
	if staged {
		t.Error("expected no staged changes after clean commit")
	}
}

func TestHasStagedChanges_WithStaged(t *testing.T) {
	dir := initRepo(t)
	f := filepath.Join(dir, "x.txt")
	writeFile(t, f, "original")
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "commit", "-m", "init")

	writeFile(t, f, "modified")
	gitRun(t, dir, "add", "-A")

	staged, err := HasStagedChanges(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !staged {
		t.Error("expected staged changes")
	}
}
