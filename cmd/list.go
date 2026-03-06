package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shinokada/svelte-next/internal/list"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List Svelte projects in a directory (read-only audit)",
	Long: `Scans [directory] (default: current dir) and prints a summary table of
every Svelte project found — current version, package manager, git status —
without making any changes.

Use --output json to emit machine-readable JSON.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

var (
	listExclude string
	listOutput  string
	listDebug   bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	f := listCmd.Flags()
	f.StringVar(&listExclude, "exclude", "", "Comma-separated directory names to skip")
	f.StringVar(&listOutput, "output", "table", "Output format: table or json")
	f.BoolVarP(&listDebug, "debug", "d", false, "Debug output")
}

func runList(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	var exclude []string
	if listExclude != "" {
		for _, s := range strings.Split(listExclude, ",") {
			if t := strings.TrimSpace(s); t != "" {
				exclude = append(exclude, t)
			}
		}
	}

	opts := list.Options{
		TargetDir: dir,
		Exclude:   exclude,
		Output:    listOutput,
		Debug:     listDebug,
	}

	if err := list.Run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}
