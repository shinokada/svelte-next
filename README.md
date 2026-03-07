# svelte-next

> Automate Svelte 5+ version updates across multiple project directories.

`svelte-next` scans a target directory, detects Svelte 5+ projects, and runs package updates, svelte installs, integration/e2e tests, and git commits — all in one command.

Built in Go in a single static binary. It runs natively on macOS, Linux, and Windows.

---

## Features

- **Automatic package manager detection** — bun, pnpm, yarn, or npm, detected from lock files
- **Svelte 5+ targeting** — projects below major version 5 are automatically skipped
- **`list` subcommand** — read-only audit of all projects before making any changes
- **`--dry-run`** — preview every action without executing anything
- **`--exclude`** — skip specific directories by name
- **`--latest` / `-L`** — update all packages to latest, ignoring semver ranges
- **`--from` / `-f`** — resume from a specific subdirectory index
- **Git workflow** — `git add -A`, `git commit`, `git push` after each successful update
- **Integration/e2e tests** — runs `test:integration` or `test:e2e` scripts if present
- **Motivational quote** — fetched from a public API at the end of each run

---

## Installation

### Homebrew (macOS and Linux)

```sh
brew tap shinokada/homebrew-svelte-next
brew install svelte-next
```

To upgrade:

```sh
brew upgrade svelte-next
```

### Debian / Ubuntu (deb)

Download the `.deb` for your architecture from the [releases page](https://github.com/shinokada/svelte-next/releases/latest):

```sh
sudo dpkg -i svelte-next_<version>_linux_amd64.deb   # amd64
sudo dpkg -i svelte-next_<version>_linux_arm64.deb   # arm64
```

### Fedora / RHEL (rpm)

```sh
sudo rpm -i svelte-next_<version>_linux_amd64.rpm    # amd64
sudo rpm -i svelte-next_<version>_linux_arm64.rpm    # arm64
```

To upgrade an existing install, replace `-i` with `-U`.

### Windows

Download the `.zip` for your architecture from the [releases page](https://github.com/shinokada/svelte-next/releases/latest), extract it, and add the directory to your `PATH`. No WSL, Git Bash, or Cygwin needed — the binary runs natively in CMD and PowerShell.

### Manual (tarball)

```sh
VERSION=$(curl -s https://api.github.com/repos/shinokada/svelte-next/releases/latest \
  | grep '"tag_name"' | cut -d'"' -f4 | cut -c2-)
PLATFORM=darwin_arm64  # darwin_amd64 | linux_amd64 | linux_arm64

curl -sLO "https://github.com/shinokada/svelte-next/releases/download/v${VERSION}/svelte-next_${VERSION}_${PLATFORM}.tar.gz"
tar -xzf "svelte-next_${VERSION}_${PLATFORM}.tar.gz"
sudo mv svelte-next /usr/local/bin/
```

### Build from source

```sh
git clone https://github.com/shinokada/svelte-next.git
cd svelte-next
make build          # binary written to ./bin/svelte-next
make install        # installs to $GOPATH/bin
```

---

## Requirements

At runtime the binary only needs:

- A supported package manager on `PATH`: **bun**, **pnpm**, **yarn**, or **npm**
- **git** (only if the git workflow is enabled)

---

## Package Manager Detection

Lock-file priority (first match wins):

| Lock file           | Package manager |
| ------------------- | --------------- |
| `bun.lockb`         | bun             |
| `bun.lock`          | bun (≥ 1.1)     |
| `pnpm-lock.yaml`    | pnpm            |
| `yarn.lock`         | yarn            |
| `package-lock.json` | npm             |
| *(none found)*      | pnpm (default)  |

---

## Usage

### `list` — audit projects without making changes

```sh
# List all Svelte projects in the current directory
svelte-next list .

# List projects, skipping 'api' and 'docs' directories
svelte-next list --exclude "api,docs" .

# Output as JSON (useful for scripting)
svelte-next list --output json .

# List projects in a specific directory
svelte-next list ./projects
```

### `update` — update Svelte projects

```sh
# Update all Svelte 5+ projects in the current directory
svelte-next update .

# Preview every action without executing anything
svelte-next update --dry-run .

# Update projects in a specific directory
svelte-next update ./projects

# Update all packages to latest (ignores semver ranges)
svelte-next update -L .

# Install a specific Svelte version
svelte-next update -n 5.28.1 .

# Skip directories named 'api' and 'docs'
svelte-next update --exclude "api,docs" .

# Start from subdirectory index 3 (useful for resuming)
svelte-next update -f 3 .

# Skip individual workflow steps
svelte-next update -p .    # skip package manager update
svelte-next update -s .    # skip svelte install
svelte-next update -t .    # skip integration/e2e tests
svelte-next update -g .    # skip git add/commit/push

# Combine flags
svelte-next update -pg .
svelte-next update -pst .

# Debug output
svelte-next update -d .
```

### Other

```sh
svelte-next --version
svelte-next --help
svelte-next update --help
svelte-next list --help
```

---

## Flag Reference

### `update`

| Flag            | Short | Default | Description                                    |
| --------------- | ----- | ------- | ---------------------------------------------- |
| `--latest`      | `-L`  | false   | Update all packages to latest (ignores semver) |
| `--skip-pkg`    | `-p`  | false   | Skip package manager update                    |
| `--skip-svelte` | `-s`  | false   | Skip svelte install                            |
| `--skip-test`   | `-t`  | false   | Skip integration/e2e tests                     |
| `--skip-git`    | `-g`  | false   | Skip git add/commit/push                       |
| `--dry-run`     | —     | false   | Preview actions without executing              |
| `--exclude`     | —     | —       | Comma-separated directory names to skip        |
| `--from`        | `-f`  | 0       | Start at subdirectory index N                  |
| `--next`        | `-n`  | —       | Specific Svelte version to install             |
| `--debug`       | `-d`  | false   | Debug output                                   |

### `list`

| Flag        | Default | Description                             |
| ----------- | ------- | --------------------------------------- |
| `--exclude` | —       | Comma-separated directory names to skip |
| `--output`  | `table` | Output format: `table` or `json`        |
| `--debug`   | false   | Debug output                            |

---

## Package Manager Commands

| Action        | pnpm           | npm                        | yarn                    | bun                   |
| ------------- | -------------- | -------------------------- | ----------------------- | --------------------- |
| Install       | `pnpm add`     | `npm install`              | `yarn add`              | `bun add`             |
| Update        | `pnpm update`  | `npm update`               | `yarn upgrade`          | `bun update`          |
| Update latest | `pnpm up -L`   | `npx npm-check-updates -u` | `yarn upgrade --latest` | `bun update --latest` |
| Run script    | `pnpm`         | `npm run`                  | `yarn`                  | `bun`                 |

> **npm note:** `--latest` uses [npm-check-updates](https://github.com/raineorshine/npm-check-updates) since npm has no native equivalent. This rewrites `package.json` before reinstalling.

---

## How It Works

For each subdirectory in the target path, `svelte-next update` does the following:

1. Checks for `package.json` with a `svelte` dependency
2. Skips the directory if Svelte major version is below 5
3. Detects the package manager from lock files
4. Runs the package manager update (unless `-p`)
5. Installs the target Svelte version (unless `-s`)
6. Runs `test:integration` or `test:e2e` if the script exists (unless `-t`)
7. Runs `git add -A`, `git commit`, `git push` if inside a git repo (unless `-g`)

At the end of a successful run, a motivational quote is fetched from a public API.

With `--dry-run`, every step is printed as `[dry-run] <command>` and nothing is executed.

---

## Notes

- Only subdirectories are scanned — files at the top level of the target directory are ignored
- Hidden directories (names starting with `.`) are always skipped
- `--exclude` matches exact directory names (case-sensitive); glob patterns are not supported in v1
- If no lock file is found, pnpm is used as the default package manager
- The git commit message is automatically set to `chore: update svelte to <version>`

---

## License

MIT — see [LICENSE](./LICENSE).
