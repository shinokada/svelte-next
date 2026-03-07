// Package packagejson provides helpers for reading and inspecting a Node.js
// package.json file. It is deliberately dependency-free (stdlib only) so it
// can be used safely in tests without any network or exec requirements.
package packagejson

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// PackageJSON represents the fields of package.json that svelte-next cares
// about. Unknown fields are silently ignored.
type PackageJSON struct {
	Scripts              map[string]string `json:"scripts"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

// Read parses the package.json at path and returns a PackageJSON struct.
func Read(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p PackageJSON
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// SvelteVersion returns the raw version string for the "svelte" package as
// recorded in any of the four dependency maps, with leading "^" or "~"
// stripped. Returns "" if svelte is not listed.
func (p *PackageJSON) SvelteVersion() string {
	v := lookupDep(p, "svelte")
	return stripPrefix(v)
}

// svelteMajorRe matches the leading major version number from semver-like
// specifiers, including plain versions (5.0.0), range prefixes (^~<>=v),
// comparator ranges (>=5 <6), workspace: aliases, and npm: aliases.
// Intentionally does not match file:, git+, http:, or other non-semver
// protocols to avoid false-positive version detection.
var svelteMajorRe = regexp.MustCompile(`^(?:workspace:)?(?:npm:svelte@)?\s*[~^<>=v]*\s*(\d+)`)

// SvelteMajor parses the major version number from SvelteVersion().
// Returns (major, true) on success, (0, false) if svelte is not present or
// the version string cannot be parsed.
func (p *PackageJSON) SvelteMajor() (int, bool) {
	v := p.SvelteVersion()
	if v == "" {
		return 0, false
	}
	// Accepted: 5.28.1 / ^5 / ~5.0.0 / >=5 / >=5 <6 / workspace:^5 / npm:svelte@5
	// Rejected: file:../path / git+https://... / http://... (non-semver protocols)
	// Unions like "^4 || ^5" return the highest major found across all clauses.
	best := -1
	for _, clause := range strings.Split(v, "||") {
		m := svelteMajorRe.FindStringSubmatch(strings.TrimSpace(clause))
		if len(m) != 2 {
			continue
		}
		major, err := strconv.Atoi(m[1])
		if err == nil && major > best {
			best = major
		}
	}
	if best == -1 {
		return 0, false
	}
	return best, true
}

// HasScript returns true if the named script exists in the "scripts" map.
func (p *PackageJSON) HasScript(name string) bool {
	if p.Scripts == nil {
		return false
	}
	_, ok := p.Scripts[name]
	return ok
}

// ScriptValue returns the command string for the named script, or "".
func (p *PackageJSON) ScriptValue(name string) string {
	if p.Scripts == nil {
		return ""
	}
	return p.Scripts[name]
}

// HasSvelte returns true if the "svelte" package appears in any dependency
// section.
func (p *PackageJSON) HasSvelte() bool {
	return lookupDep(p, "svelte") != ""
}

// SvelteIsDevDependency returns true if "svelte" is declared in
// devDependencies (and not in any other section). This is a convenience
// wrapper around SvelteDependencySection.
func (p *PackageJSON) SvelteIsDevDependency() bool {
	sections := p.SvelteDependencySections()
	return len(sections) == 1 && sections[0] == "devDependencies"
}

// SvelteDependencySection returns the name of the first dependency map that
// contains "svelte". Returns "" if svelte is not listed.
// Deprecated: use SvelteDependencySections when you need all matching buckets.
func (p *PackageJSON) SvelteDependencySection() string {
	sections := p.SvelteDependencySections()
	if len(sections) == 0 {
		return ""
	}
	return sections[0]
}

// SvelteDependencySections returns the names of every dependency map that
// contains "svelte" (e.g. both "peerDependencies" and "devDependencies" for
// library packages). Returns nil if svelte is not listed in any section.
func (p *PackageJSON) SvelteDependencySections() []string {
	type section struct {
		name string
		m    map[string]string
	}
	all := []section{
		{"dependencies", p.Dependencies},
		{"devDependencies", p.DevDependencies},
		{"peerDependencies", p.PeerDependencies},
		{"optionalDependencies", p.OptionalDependencies},
	}
	var found []string
	for _, s := range all {
		if _, ok := s.m["svelte"]; ok {
			found = append(found, s.name)
		}
	}
	return found
}

// lookupDep searches all four dependency maps for key and returns its value,
// or "" if not found. Priority: dependencies > devDependencies >
// peerDependencies > optionalDependencies.
func lookupDep(p *PackageJSON, key string) string {
	for _, m := range []map[string]string{
		p.Dependencies,
		p.DevDependencies,
		p.PeerDependencies,
		p.OptionalDependencies,
	} {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return ""
}

// stripPrefix removes a leading "^" or "~" from a semver range string.
func stripPrefix(v string) string {
	return strings.TrimLeft(v, "^~")
}
