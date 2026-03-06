package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// mkdirs creates subdirectories under a temp dir and returns the temp dir path.
func mkdirs(t *testing.T, names ...string) string {
	t.Helper()
	root := t.TempDir()
	for _, name := range names {
		if err := os.MkdirAll(filepath.Join(root, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestScan_Basic(t *testing.T) {
	root := mkdirs(t, "alpha", "beta", "gamma")
	dirs, err := Scan(Options{TargetDir: root})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 3 {
		t.Errorf("expected 3 dirs, got %d: %v", len(dirs), dirs)
	}
}

func TestScan_SkipsHidden(t *testing.T) {
	root := mkdirs(t, "visible", ".hidden", ".git")
	dirs, err := Scan(Options{TargetDir: root})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 {
		t.Errorf("expected 1 visible dir, got %d: %v", len(dirs), dirs)
	}
	if filepath.Base(dirs[0]) != "visible" {
		t.Errorf("unexpected dir name: %s", dirs[0])
	}
}

func TestScan_Exclude(t *testing.T) {
	root := mkdirs(t, "app", "api", "docs", "legacy")
	dirs, err := Scan(Options{TargetDir: root, Exclude: []string{"api", "docs"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 2 {
		t.Errorf("expected 2 dirs after exclude, got %d: %v", len(dirs), dirs)
	}
	for _, d := range dirs {
		name := filepath.Base(d)
		if name == "api" || name == "docs" {
			t.Errorf("excluded dir %q still present", name)
		}
	}
}

func TestScan_Sorted(t *testing.T) {
	root := mkdirs(t, "zebra", "alpha", "mango")
	dirs, err := Scan(Options{TargetDir: root})
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(dirs))
	for i, d := range dirs {
		names[i] = filepath.Base(d)
	}
	want := []string{"alpha", "mango", "zebra"}
	for i, w := range want {
		if names[i] != w {
			t.Errorf("position %d: got %q, want %q", i, names[i], w)
		}
	}
}

func TestScan_SkipsFiles(t *testing.T) {
	root := t.TempDir()
	// Create a file (not a dir) alongside a directory.
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	dirs, err := Scan(Options{TargetDir: root})
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 || filepath.Base(dirs[0]) != "app" {
		t.Errorf("unexpected dirs: %v", dirs)
	}
}

func TestScan_NonexistentDir(t *testing.T) {
	_, err := Scan(Options{TargetDir: "/nonexistent/path/xyz"})
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestScan_FileAsTarget(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Scan(Options{TargetDir: file})
	if err == nil {
		t.Error("expected error when target is a file, not a dir")
	}
}
