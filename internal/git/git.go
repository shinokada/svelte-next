// Package git provides thin wrappers around git CLI commands used by
// svelte-next. Every mutating function accepts a dryRun bool; when true the
// command is printed but not executed.
package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/shinokada/svelte-next/internal/ui"
)

// IsGitRepo returns true if dir is inside a git repository (i.e. `git rev-parse
// --git-dir` succeeds).
func IsGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	return cmd.Run() == nil
}

// IsClean returns true if the working tree in dir has no uncommitted changes.
func IsClean(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("git: status: %w", err)
	}
	return len(strings.TrimSpace(string(out))) == 0, nil
}

// CurrentBranch returns the name of the current branch in dir.
func CurrentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git: branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))
	if branch == "" {
		return "", fmt.Errorf("git: detached HEAD state (no current branch)")
	}
	return branch, nil
}

// Add stages all changes in dir (`git add -A`).
func Add(dir string, dryRun bool) error {
	if dryRun {
		ui.DryRunf("cd %s && git add -A", dir)
		return nil
	}
	return run(dir, "git", "add", "-A")
}

// HasStagedChanges returns true if there is at least one staged change in dir.
func HasStagedChanges(dir string) (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = dir
	err := cmd.Run()
	if err == nil {
		return false, nil // exit 0 → nothing staged
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return true, nil // exit 1 → staged changes present
	}
	return false, fmt.Errorf("git: diff --cached: %w", err)
}

// Commit creates a commit with message in dir.
func Commit(dir, message string, dryRun bool) error {
	if dryRun {
		ui.DryRunf("cd %s && git commit --message %q", dir, message)
		return nil
	}
	return run(dir, "git", "commit", "--message", message)
}

// Push pushes the current branch to origin in dir.
func Push(dir, branch string, dryRun bool) error {
	if dryRun {
		ui.DryRunf("cd %s && git push origin %s", dir, branch)
		return nil
	}
	return run(dir, "git", "push", "origin", branch)
}

// run executes a command in dir, streaming stdout/stderr to the terminal.
func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}


