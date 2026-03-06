// Package pkgmanager handles package manager detection and command execution
// for svelte-next. Detection is based on lock-file presence; execution
// delegates to the detected package manager binary via os/exec.
package pkgmanager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shinokada/svelte-next/internal/ui"
)

// Manager identifies a Node.js package manager.
type Manager string

const (
	Bun  Manager = "bun"
	Pnpm Manager = "pnpm"
	Yarn Manager = "yarn"
	Npm  Manager = "npm"
)

// Detect returns the package manager for dir based on lock-file priority:
//  1. bun.lockb  or bun.lock  (Bun ≥ 1.1 uses text-format bun.lock)
//  2. pnpm-lock.yaml
//  3. yarn.lock
//  4. package-lock.json
//  5. pnpm (default)
func Detect(dir string) Manager {
	checks := []struct {
		file    string
		manager Manager
	}{
		{"bun.lockb", Bun},
		{"bun.lock", Bun},
		{"pnpm-lock.yaml", Pnpm},
		{"yarn.lock", Yarn},
		{"package-lock.json", Npm},
	}
	for _, c := range checks {
		if fileExists(filepath.Join(dir, c.file)) {
			return c.manager
		}
	}
	return Pnpm
}

// Run executes a package-manager command in dir. If dryRun is true the
// command is printed (prefixed "[dry-run]") but not executed.
//
// Supported cmd values:
//
//	"install"        — add a package (args: e.g. ["-D", "svelte@latest"])
//	"update"         — update packages within semver ranges
//	"update-latest"  — update all packages ignoring semver ranges (-L flag)
//	"run"            — run a script (args: [scriptName])
func Run(dir string, mgr Manager, dryRun bool, cmd string, args ...string) error {
	argv, err := buildArgv(mgr, cmd, args...)
	if err != nil {
		return err
	}

	if dryRun {
		ui.DryRunf("cd %s && %s %v", dir, string(mgr), argv)
		return nil
	}

	c := exec.Command(argv[0], argv[1:]...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// CheckInstalled returns an error if the binary for mgr is not found in PATH.
func CheckInstalled(mgr Manager) error {
	_, err := exec.LookPath(string(mgr))
	if err != nil {
		return fmt.Errorf("package manager %q not found in PATH: %w", mgr, err)
	}
	return nil
}

// buildArgv constructs the argv slice for a given manager + command pair.
func buildArgv(mgr Manager, cmd string, args ...string) ([]string, error) {
	switch mgr {
	case Bun:
		return bunArgv(cmd, args...)
	case Pnpm:
		return pnpmArgv(cmd, args...)
	case Yarn:
		return yarnArgv(cmd, args...)
	case Npm:
		return npmArgv(cmd, args...)
	default:
		return nil, fmt.Errorf("pkgmanager: unknown manager %q", mgr)
	}
}

func bunArgv(cmd string, args ...string) ([]string, error) {
	switch cmd {
	case "install":
		return append([]string{"bun", "add"}, args...), nil
	case "update":
		return append([]string{"bun", "update"}, args...), nil
	case "update-latest":
		return []string{"bun", "update", "--latest"}, nil
	case "run":
		return append([]string{"bun"}, args...), nil
	}
	return nil, fmt.Errorf("pkgmanager: unknown command %q for bun", cmd)
}

func pnpmArgv(cmd string, args ...string) ([]string, error) {
	switch cmd {
	case "install":
		return append([]string{"pnpm", "add"}, args...), nil
	case "update":
		return append([]string{"pnpm", "update"}, args...), nil
	case "update-latest":
		return []string{"pnpm", "up", "-L"}, nil
	case "run":
		return append([]string{"pnpm"}, args...), nil
	}
	return nil, fmt.Errorf("pkgmanager: unknown command %q for pnpm", cmd)
}

func yarnArgv(cmd string, args ...string) ([]string, error) {
	switch cmd {
	case "install":
		return append([]string{"yarn", "add"}, args...), nil
	case "update":
		return append([]string{"yarn", "upgrade"}, args...), nil
	case "update-latest":
		// yarn upgrade --latest (classic v1).
		// Berry (v2+) would need "yarn up '*'" but we use the simple form here;
		// the full version-detection logic can be added in a future release.
		return []string{"yarn", "upgrade", "--latest"}, nil
	case "run":
		return append([]string{"yarn"}, args...), nil
	}
	return nil, fmt.Errorf("pkgmanager: unknown command %q for yarn", cmd)
}

func npmArgv(cmd string, args ...string) ([]string, error) {
	switch cmd {
	case "install":
		return append([]string{"npm", "install"}, args...), nil
	case "update":
		return append([]string{"npm", "update"}, args...), nil
	case "update-latest":
		// npm has no native "update to latest" — use npx npm-check-updates.
		// The caller is expected to follow up with a separate "install" call.
		return []string{"npx", "--yes", "npm-check-updates", "-u"}, nil
	case "run":
		return append([]string{"npm", "run"}, args...), nil
	}
	return nil, fmt.Errorf("pkgmanager: unknown command %q for npm", cmd)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
