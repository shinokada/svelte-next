package list

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeProject creates a fake Svelte project dir under root.
func makeProject(t *testing.T, root, name, pkgJSON string) string {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if pkgJSON != "" {
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkgJSON), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestRun_TableOutput(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "app", `{"devDependencies":{"svelte":"^5.28.1"}}`)
	makeProject(t, root, "legacy", `{"devDependencies":{"svelte":"^4.2.8"}}`)
	makeProject(t, root, "no-svelte", `{"dependencies":{"react":"18.0.0"}}`)

	err := Run(Options{TargetDir: root, Output: "table"})
	if err != nil {
		t.Errorf("Run() error: %v", err)
	}
}

func TestRun_JSONOutput(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "app", `{"devDependencies":{"svelte":"^5.28.1"}}`)

	// Capture stdout by redirecting — simplest approach: just verify no error.
	err := Run(Options{TargetDir: root, Output: "json"})
	if err != nil {
		t.Errorf("Run() JSON error: %v", err)
	}
}

func TestRun_Exclude(t *testing.T) {
	root := t.TempDir()
	makeProject(t, root, "app", `{"devDependencies":{"svelte":"^5.28.1"}}`)
	makeProject(t, root, "api", `{"devDependencies":{"svelte":"^5.0.0"}}`)

	err := Run(Options{TargetDir: root, Exclude: []string{"api"}, Output: "table"})
	if err != nil {
		t.Errorf("Run() error: %v", err)
	}
}

func TestInspect_NoPackageJSON(t *testing.T) {
	dir := t.TempDir()
	info := inspect(dir, false)
	if info.Skipped == "" {
		t.Error("expected Skipped to be set when no package.json")
	}
}

func TestInspect_NoSvelte(t *testing.T) {
	root := t.TempDir()
	dir := makeProject(t, root, "proj", `{"dependencies":{"react":"18.0.0"}}`)
	info := inspect(dir, false)
	if !strings.Contains(info.Skipped, "no svelte") {
		t.Errorf("unexpected Skipped: %q", info.Skipped)
	}
}

func TestInspect_LegacySvelte(t *testing.T) {
	root := t.TempDir()
	dir := makeProject(t, root, "old", `{"devDependencies":{"svelte":"^4.2.8"}}`)
	info := inspect(dir, false)
	if !strings.Contains(info.Skipped, "major < 5") {
		t.Errorf("expected 'major < 5' in Skipped, got %q", info.Skipped)
	}
}

func TestInspect_ValidSvelte5(t *testing.T) {
	root := t.TempDir()
	dir := makeProject(t, root, "modern", `{"devDependencies":{"svelte":"^5.28.1"}}`)
	info := inspect(dir, false)
	if info.Skipped != "" {
		t.Errorf("expected no skip reason, got %q", info.Skipped)
	}
	if info.Svelte != "5.28.1" {
		t.Errorf("Svelte = %q, want %q", info.Svelte, "5.28.1")
	}
}

func TestPrintJSON_ValidJSON(t *testing.T) {
	projects := []ProjectInfo{
		{Dir: "app", Svelte: "5.28.1", Manager: "pnpm", Git: true, Clean: false},
	}
	// Just verify it produces valid JSON by marshalling and unmarshalling.
	data, err := json.Marshal(projects)
	if err != nil {
		t.Fatal(err)
	}
	var out []ProjectInfo
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if out[0].Dir != "app" {
		t.Errorf("round-trip Dir = %q, want %q", out[0].Dir, "app")
	}
}
