// Package scanner provides shared directory scanning and --exclude filtering
// used by both the update and list subcommands.
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Options configures a Scan call.
type Options struct {
	// TargetDir is the directory whose immediate subdirectories are scanned.
	TargetDir string
	// Exclude is a list of exact subdirectory names (not paths) to skip.
	// Matching is case-sensitive. Glob patterns are not supported in v1.
	Exclude []string
}

// Scan returns a sorted list of absolute paths of non-hidden subdirectories
// directly inside opts.TargetDir, with any names in opts.Exclude removed.
func Scan(opts Options) ([]string, error) {
	abs, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return nil, fmt.Errorf("scanner: resolving path: %w", err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("scanner: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scanner: %q is not a directory", abs)
	}

	entries, err := os.ReadDir(abs)
	if err != nil {
		return nil, fmt.Errorf("scanner: reading directory: %w", err)
	}

	excludeSet := make(map[string]struct{}, len(opts.Exclude))
	for _, name := range opts.Exclude {
		excludeSet[name] = struct{}{}
	}

	var dirs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		// Skip hidden directories (names starting with ".").
		if len(name) > 0 && name[0] == '.' {
			continue
		}
		// Skip excluded names.
		if _, excluded := excludeSet[name]; excluded {
			continue
		}
		dirs = append(dirs, filepath.Join(abs, name))
	}

	sort.Strings(dirs)
	return dirs, nil
}
