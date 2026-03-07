package packagejson

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// writeTemp writes content to a temp file and returns its path.
func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRead_DevDependency(t *testing.T) {
	path := writeTemp(t, `{
		"devDependencies": { "svelte": "^5.28.1" }
	}`)
	p, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if !p.HasSvelte() {
		t.Error("HasSvelte() = false, want true")
	}
	if got := p.SvelteVersion(); got != "5.28.1" {
		t.Errorf("SvelteVersion() = %q, want %q", got, "5.28.1")
	}
}

func TestRead_Dependency(t *testing.T) {
	path := writeTemp(t, `{
		"dependencies": { "svelte": "~5.0.0" }
	}`)
	p, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := p.SvelteVersion(); got != "5.0.0" {
		t.Errorf("SvelteVersion() = %q, want %q", got, "5.0.0")
	}
}

func TestSvelteMajor(t *testing.T) {
	tests := []struct {
		json      string
		wantMajor int
		wantOK    bool
	}{
		{`{"devDependencies":{"svelte":"^5.28.1"}}`, 5, true},
		{`{"devDependencies":{"svelte":"^4.2.8"}}`, 4, true},
		{`{"devDependencies":{"svelte":"5.0.0-next.1"}}`, 5, true},
		{`{"devDependencies":{"svelte":">=5"}}`, 5, true},
		{`{"devDependencies":{"svelte":">=5 <6"}}`, 5, true},
		{`{"devDependencies":{"svelte":">= 5.0.0"}}`, 5, true},
		{`{"devDependencies":{"svelte":"workspace:^5.0.0"}}`, 5, true},
		{`{"devDependencies":{"svelte":"npm:svelte@5.0.0"}}`, 5, true},
		{`{"devDependencies":{"svelte":"^4 || ^5"}}`, 5, true},
		{`{"devDependencies":{"svelte":"^3 || ^4"}}`, 4, true},
		{`{"devDependencies":{"svelte":"^5 || ^4"}}`, 5, true},
		{`{"devDependencies":{"svelte":"file:../svelte5-local"}}`, 0, false},
		{`{"devDependencies":{"svelte":"git+https://github.com/sveltejs/svelte"}}`, 0, false},
		{`{"devDependencies":{"svelte":"https://example.com/svelte.tgz"}}`, 0, false},
		{`{"dependencies":{"react":"18.0.0"}}`, 0, false},
		{`{}`, 0, false},
	}
	for _, tc := range tests {
		path := writeTemp(t, tc.json)
		p, err := Read(path)
		if err != nil {
			t.Fatal(err)
		}
		gotMajor, gotOK := p.SvelteMajor()
		if gotMajor != tc.wantMajor || gotOK != tc.wantOK {
			t.Errorf("SvelteMajor() = (%d, %v), want (%d, %v) for %s",
				gotMajor, gotOK, tc.wantMajor, tc.wantOK, tc.json)
		}
	}
}

func TestHasScript(t *testing.T) {
	path := writeTemp(t, `{
		"scripts": {
			"test:integration": "playwright test",
			"build": "vite build"
		}
	}`)
	p, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if !p.HasScript("test:integration") {
		t.Error("expected HasScript(test:integration) = true")
	}
	if p.HasScript("test:e2e") {
		t.Error("expected HasScript(test:e2e) = false")
	}
}

func TestRead_MissingFile(t *testing.T) {
	_, err := Read("/nonexistent/package.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("expected fs.ErrNotExist, got %v", err)
	}
}

func TestSvelteIsDevDependency(t *testing.T) {
	tests := []struct {
		name string
		json string
		want bool
	}{
		{"in devDependencies", `{"devDependencies":{"svelte":"^5.0.0"}}`, true},
		{"in dependencies", `{"dependencies":{"svelte":"^5.0.0"}}`, false},
		{"in peerDependencies", `{"peerDependencies":{"svelte":"^5.0.0"}}`, false},
		{"in dev and peer", `{"devDependencies":{"svelte":"^5.0.0"},"peerDependencies":{"svelte":"^5.0.0"}}`, false},
		{"absent", `{}`, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTemp(t, tc.json)
			p, err := Read(path)
			if err != nil {
				t.Fatal(err)
			}
			if got := p.SvelteIsDevDependency(); got != tc.want {
				t.Errorf("SvelteIsDevDependency() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSvelteDependencySections(t *testing.T) {
	tests := []struct {
		name string
		json string
		want []string
	}{
		{
			"peer and dev (library pattern)",
			`{"peerDependencies":{"svelte":"^5.0.0"},"devDependencies":{"svelte":"^5.0.0"}}`,
			[]string{"devDependencies", "peerDependencies"},
		},
		{
			"only devDependencies",
			`{"devDependencies":{"svelte":"^5.0.0"}}`,
			[]string{"devDependencies"},
		},
		{
			"absent",
			`{}`,
			nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTemp(t, tc.json)
			p, err := Read(path)
			if err != nil {
				t.Fatal(err)
			}
			got := p.SvelteDependencySections()
			if len(got) != len(tc.want) {
				t.Errorf("SvelteDependencySections() = %v, want %v", got, tc.want)
				return
			}
			gotSet := make(map[string]bool, len(got))
			for _, s := range got {
				gotSet[s] = true
			}
			for _, w := range tc.want {
				if !gotSet[w] {
					t.Errorf("SvelteDependencySections() = %v, missing %q", got, w)
				}
			}
		})
	}
}

func TestRead_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{ invalid json `)
	_, err := Read(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestStripPrefix(t *testing.T) {
	tests := []struct{ in, want string }{
		{"^5.0.0", "5.0.0"},
		{"~4.2.8", "4.2.8"},
		{"5.0.0", "5.0.0"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := stripPrefix(tc.in); got != tc.want {
			t.Errorf("stripPrefix(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
