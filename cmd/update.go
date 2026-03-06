package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shinokada/svelte-next/internal/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [directory]",
	Short: "Update Svelte 5+ projects in a directory",
	Long: `Iterates over non-hidden subdirectories of [directory] (default: current dir),
detects Svelte 5+ projects, and runs package updates, svelte installs,
integration tests, and git commits.

Pass --dry-run to preview every action without executing anything.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdate,
}

// flag vars
var (
	updateLatest     bool
	updateSkipPkg    bool
	updateSkipSvelte bool
	updateSkipTest   bool
	updateSkipGit    bool
	updateDebug      bool
	updateDryRun     bool
	updateFrom       int
	updateSvelteVer  string
	updateExclude    string
)

func init() {
	rootCmd.AddCommand(updateCmd)

	f := updateCmd.Flags()
	f.BoolVarP(&updateLatest, "latest", "L", false, "Update all packages to latest (ignores semver ranges)")
	f.BoolVarP(&updateSkipPkg, "skip-pkg", "p", false, "Skip package manager update")
	f.BoolVarP(&updateSkipSvelte, "skip-svelte", "s", false, "Skip svelte install")
	f.BoolVarP(&updateSkipTest, "skip-test", "t", false, "Skip integration/e2e tests")
	f.BoolVarP(&updateSkipGit, "skip-git", "g", false, "Skip git add/commit/push")
	f.BoolVarP(&updateDebug, "debug", "d", false, "Debug output")
	f.BoolVar(&updateDryRun, "dry-run", false, "Preview actions without executing them")
	f.IntVarP(&updateFrom, "from", "f", 0, "Start processing at subdirectory index N")
	f.StringVarP(&updateSvelteVer, "next", "n", "", "Specific Svelte version to install (e.g. 5.28.1)")
	f.StringVar(&updateExclude, "exclude", "", "Comma-separated directory names to skip (e.g. api,docs)")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	var exclude []string
	if updateExclude != "" {
		for _, s := range strings.Split(updateExclude, ",") {
			if t := strings.TrimSpace(s); t != "" {
				exclude = append(exclude, t)
			}
		}
	}

	opts := update.Options{
		TargetDir:  dir,
		SvelteVer:  updateSvelteVer,
		Latest:     updateLatest,
		SkipPkg:    updateSkipPkg,
		SkipSvelte: updateSkipSvelte,
		SkipTest:   updateSkipTest,
		SkipGit:    updateSkipGit,
		Debug:      updateDebug,
		DryRun:     updateDryRun,
		From:       updateFrom,
		Exclude:    exclude,
	}

	if err := update.Run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}
