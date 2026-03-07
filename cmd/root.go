package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "svelte-next",
	Short: "Update Svelte 5+ projects across multiple directories",
	Long: `svelte-next scans a target directory, detects Svelte 5+ projects,
and runs package updates, svelte installs, tests, and git commits.

Use --dry-run on any subcommand to preview actions without executing them.`,
}

// Execute is called by main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("svelte-next {{.Version}}\n")
}
