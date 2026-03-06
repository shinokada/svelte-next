package pkgmanager

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		lockFile string
		want     Manager
	}{
		{"bun.lockb", Bun},
		{"bun.lock", Bun},
		{"pnpm-lock.yaml", Pnpm},
		{"yarn.lock", Yarn},
		{"package-lock.json", Npm},
	}

	for _, tc := range tests {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, tc.lockFile), []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
		got := Detect(dir)
		if got != tc.want {
			t.Errorf("Detect() with %q = %q, want %q", tc.lockFile, got, tc.want)
		}
	}
}

func TestDetect_Default(t *testing.T) {
	dir := t.TempDir() // no lock files
	if got := Detect(dir); got != Pnpm {
		t.Errorf("Detect() default = %q, want %q", got, Pnpm)
	}
}

func TestDetect_Priority_BunOverPnpm(t *testing.T) {
	dir := t.TempDir()
	// Both bun.lockb and pnpm-lock.yaml present — bun wins.
	for _, f := range []string{"bun.lockb", "pnpm-lock.yaml"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got := Detect(dir); got != Bun {
		t.Errorf("expected Bun to win priority, got %q", got)
	}
}

func TestBuildArgv_Pnpm(t *testing.T) {
	tests := []struct {
		cmd  string
		args []string
		want []string
	}{
		{"install", []string{"-D", "svelte@latest"}, []string{"pnpm", "add", "-D", "svelte@latest"}},
		{"update", nil, []string{"pnpm", "update"}},
		{"update-latest", nil, []string{"pnpm", "up", "-L"}},
		{"run", []string{"test:integration"}, []string{"pnpm", "test:integration"}},
	}
	for _, tc := range tests {
		got, err := buildArgv(Pnpm, tc.cmd, tc.args...)
		if err != nil {
			t.Errorf("buildArgv(pnpm, %q) error: %v", tc.cmd, err)
			continue
		}
		if !slices.Equal(got, tc.want) {
			t.Errorf("buildArgv(pnpm, %q) = %v, want %v", tc.cmd, got, tc.want)
		}
	}
}

func TestBuildArgv_Npm(t *testing.T) {
	got, err := buildArgv(Npm, "update-latest")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"npx", "--yes", "npm-check-updates", "-u"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgv_Bun(t *testing.T) {
	got, err := buildArgv(Bun, "update-latest")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"bun", "update", "--latest"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgv_Yarn(t *testing.T) {
	got, err := buildArgv(Yarn, "update-latest")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"yarn", "upgrade", "--latest"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgv_UnknownCmd(t *testing.T) {
	_, err := buildArgv(Pnpm, "nonexistent-cmd")
	if err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestBuildArgv_UnknownManager(t *testing.T) {
	_, err := buildArgv("pip", "install")
	if err == nil {
		t.Error("expected error for unknown manager")
	}
}


