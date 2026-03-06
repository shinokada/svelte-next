// Package list implements the `svelte-next list` subcommand: a read-only audit
// that scans a target directory and prints a summary table of every Svelte
// project found, without making any changes.
package list

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shinokada/svelte-next/internal/git"
	"github.com/shinokada/svelte-next/internal/packagejson"
	"github.com/shinokada/svelte-next/internal/pkgmanager"
	"github.com/shinokada/svelte-next/internal/scanner"
	"github.com/shinokada/svelte-next/internal/ui"
)

// Options configures a Run call.
type Options struct {
	TargetDir string
	Exclude   []string
	Output    string // "table" (default) | "json"
	Debug     bool
}

// ProjectInfo holds the read-only state of a single project directory.
type ProjectInfo struct {
	Dir     string `json:"dir"`
	Svelte  string `json:"svelte"`  // version string or "—"
	Manager string `json:"manager"` // detected package manager or "—"
	Git     bool   `json:"git"`
	Clean   bool   `json:"clean"`   // git working tree is clean
	Skipped string `json:"skipped"` // non-empty reason if the project was skipped
}

// Run scans opts.TargetDir and prints a summary of each Svelte project.
func Run(opts Options) error {
	dirs, err := scanner.Scan(scanner.Options{
		TargetDir: opts.TargetDir,
		Exclude:   opts.Exclude,
	})
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}

	if opts.Debug {
		ui.Infof("list: found %d directories to inspect", len(dirs))
	}

	var projects []ProjectInfo

	for _, dir := range dirs {
		info := inspect(dir, opts.Debug)
		projects = append(projects, info)
	}

	switch opts.Output {
	case "json":
		return printJSON(projects)
	default:
		printTable(projects)
		return nil
	}
}

// inspect gathers read-only information about a single directory.
func inspect(dir string, debug bool) ProjectInfo {
	name := filepath.Base(dir)
	info := ProjectInfo{
		Dir:    name,
		Svelte: "—",
	}

	pkgPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkgPath); err != nil {
		info.Skipped = "no package.json"
		if debug {
			ui.Warnf("  %s: no package.json", name)
		}
		return info
	}

	p, err := packagejson.Read(pkgPath)
	if err != nil {
		info.Skipped = "invalid package.json"
		return info
	}

	if !p.HasSvelte() {
		info.Skipped = "no svelte dependency"
		return info
	}

	info.Svelte = p.SvelteVersion()
	info.Manager = string(pkgmanager.Detect(dir))

	major, ok := p.SvelteMajor()
	if ok && major < 5 {
		info.Skipped = fmt.Sprintf("major < 5 (v%d)", major)
	}

	info.Git = git.IsGitRepo(dir)
	if info.Git {
		clean, err := git.IsClean(dir)
		if err == nil {
			info.Clean = clean
		}
	}

	return info
}

// printTable renders the projects as a human-readable table.
func printTable(projects []ProjectInfo) {
	headers := []string{"Directory", "Svelte", "Manager", "Git", "Clean", "Note"}
	rows := make([][]string, len(projects))
	for i, p := range projects {
		gitStr := boolStr(p.Git)
		cleanStr := "—"
		if p.Git {
			cleanStr = boolStr(p.Clean)
		}
		note := p.Skipped
		rows[i] = []string{p.Dir, p.Svelte, managerStr(p.Manager), gitStr, cleanStr, note}
	}
	ui.Table(headers, rows)
}

// printJSON renders the projects as a JSON array.
func printJSON(projects []ProjectInfo) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(projects)
}

func boolStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func managerStr(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
