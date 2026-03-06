// Package update implements the core `svelte-next update` loop.
package update

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/shinokada/svelte-next/internal/git"
	"github.com/shinokada/svelte-next/internal/packagejson"
	"github.com/shinokada/svelte-next/internal/pkgmanager"
	"github.com/shinokada/svelte-next/internal/quote"
	"github.com/shinokada/svelte-next/internal/scanner"
	"github.com/shinokada/svelte-next/internal/ui"
)

// Options configures a Run call.
type Options struct {
	TargetDir  string
	SvelteVer  string // "" means "latest"
	Latest     bool
	SkipPkg    bool
	SkipSvelte bool
	SkipTest   bool
	SkipGit    bool
	Debug      bool
	DryRun     bool
	From       int
	Exclude    []string
}

// quoteTimeout is the per-request deadline for fetching a motivational quote.
const quoteTimeout = 5 * time.Second

// Run iterates over subdirectories of opts.TargetDir, updating Svelte 5+
// projects in each one.
func Run(opts Options) error {
	dirs, err := scanner.Scan(scanner.Options{
		TargetDir: opts.TargetDir,
		Exclude:   opts.Exclude,
	})
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if opts.From < 0 || opts.From > len(dirs) {
		return fmt.Errorf("update: --from %d is out of range (0–%d)", opts.From, len(dirs))
	}
	dirs = dirs[opts.From:]

	// Print startup banner.
	dryTag := ""
	if opts.DryRun {
		dryTag = " [dry-run]"
	}
	ui.PrintBanner(fmt.Sprintf("svelte-next update%s\nTarget: %s  (%d dirs)",
		dryTag, opts.TargetDir, len(dirs)), ui.Blue, "*", 50)

	var hadFailure bool

	// NOTE: Each step runs independently; a failure in one step does not block
	// subsequent steps within the same project. This allows all errors to surface
	// in a single run. The git workflow is skipped when prior steps have failed.
	for _, dir := range dirs {
		name := filepath.Base(dir)
		pkgPath := filepath.Join(dir, "package.json")

		// ── 1. Check package.json + svelte dep ──────────────────────────────
		p, err := packagejson.Read(pkgPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				ui.Warnf("  [%s] skipping: no package.json", name)
				continue
			}
			ui.Errorf("  [%s] invalid package.json: %v", name, err)
			hadFailure = true
			continue
		}
		if !p.HasSvelte() {
			ui.Warnf("  [%s] skipping: no svelte dependency", name)
			continue
		}

		// ── 2. Check Svelte major version ────────────────────────────────────
		major, ok := p.SvelteMajor()
		if !ok {
			ui.Warnf("  [%s] skipping: cannot parse svelte version", name)
			continue
		}
		if major < 5 {
			ui.Warnf("  [%s] skipping: svelte major is %d (< 5)", name, major)
			continue
		}

		// ── 3. Detect package manager ────────────────────────────────────────
		mgr := pkgmanager.Detect(dir)
		ui.PrintBanner(fmt.Sprintf("[%s]  manager: %s  svelte: %s",
			name, mgr, p.SvelteVersion()), ui.Cyan, "-", 50)

		projectFailed := false

		// ── 4. Package update ────────────────────────────────────────────────
		if !opts.SkipPkg {
			cmd := "update"
			if opts.Latest {
				cmd = "update-latest"
			}
			ui.Infof("  running: %s %s", mgr, cmd)
			if err := pkgmanager.Run(dir, mgr, opts.DryRun, cmd); err != nil {
				ui.Errorf("  [%s] package update failed: %v", name, err)
				hadFailure = true
				projectFailed = true
			} else if opts.Latest && opts.SkipSvelte && mgr == pkgmanager.Npm {
				// npm update-latest only rewrites package.json via npm-check-updates;
				// a follow-up install is required to sync the lockfile and node_modules.
				ui.Infof("  running: npm install (lockfile sync after npm-check-updates)")
				if err := pkgmanager.Run(dir, mgr, opts.DryRun, "install"); err != nil {
					ui.Errorf("  [%s] npm install failed: %v", name, err)
					hadFailure = true
					projectFailed = true
				}
			}
		}

		// ── 5. Svelte install ────────────────────────────────────────────────
		if !opts.SkipSvelte {
			svelteTarget := "svelte@latest"
			if opts.SvelteVer != "" {
				svelteTarget = "svelte@" + opts.SvelteVer
			}
			ui.Infof("  installing: %s", svelteTarget)
			installArgs := []string{"install", svelteTarget}
			if p.SvelteIsDevDependency() {
				installArgs = []string{"install", "-D", svelteTarget}
			}
			if err := pkgmanager.Run(dir, mgr, opts.DryRun, installArgs...); err != nil {
				ui.Errorf("  [%s] svelte install failed: %v", name, err)
				hadFailure = true
				projectFailed = true
			}
		}

		// ── 6. Integration / e2e tests ───────────────────────────────────────
		// Both test:integration and test:e2e are run if present; skipping one
		// because the other exists could miss a failing suite.
		if !opts.SkipTest {
			for _, script := range []string{"test:integration", "test:e2e"} {
				if !p.HasScript(script) {
					continue
				}
				ui.Infof("  running: %s", script)
				if err := pkgmanager.Run(dir, mgr, opts.DryRun, "run", script); err != nil {
					ui.Errorf("  [%s] %s failed: %v", name, script, err)
					hadFailure = true
					projectFailed = true
				}
			}
		}

		// ── 7. Git workflow ──────────────────────────────────────────────────
		if !opts.SkipGit && !projectFailed && git.IsGitRepo(dir) {
			if err := git.Add(dir, opts.DryRun); err != nil {
				ui.Errorf("  [%s] git add failed: %v", name, err)
				hadFailure = true
				projectFailed = true
			} else {
				staged, err := git.HasStagedChanges(dir)
				if err != nil {
					ui.Errorf("  [%s] could not check staged changes: %v", name, err)
					hadFailure = true
					projectFailed = true
				} else if staged || opts.DryRun {
					branch, err := git.CurrentBranch(dir)
					if err != nil {
						ui.Errorf("  [%s] could not determine current branch: %v", name, err)
						hadFailure = true
						projectFailed = true
					} else {
						newVer := "latest"
						if opts.SvelteVer != "" {
							newVer = opts.SvelteVer
						}
						msg := fmt.Sprintf("chore: update svelte to %s", newVer)
						if err := git.Commit(dir, msg, opts.DryRun); err != nil {
							ui.Errorf("  [%s] git commit failed: %v", name, err)
							hadFailure = true
							projectFailed = true
						} else if err := git.Push(dir, branch, opts.DryRun); err != nil {
							ui.Errorf("  [%s] git push failed: %v", name, err)
							hadFailure = true
							projectFailed = true
						}
					}
				} else {
					ui.Infof("  [%s] nothing to commit", name)
				}
			}
		}

		if projectFailed {
			ui.Warnf("  [%s] completed with errors", name)
		} else {
			ui.Successf("  [%s] done", name)
		}
	}

	// ── 8. Motivational quote ────────────────────────────────────────────────
	if !opts.DryRun {
		if q, err := quote.Fetch(quote.DefaultAPIs, quoteTimeout); err == nil {
			fmt.Println()
			ui.Infof("💬 %s", q)
		}
	}

	if hadFailure {
		return fmt.Errorf("update: one or more projects failed")
	}
	return nil
}
