# Svelte-Next

## Automate Svelte Version Updates

This script automates updating Svelte versions in project directories. If you have multiple Svelte projects in subdirectories, this script will update them all.

## Features:

- Automatic package manager detection (supports pnpm, npm, yarn, and bun)
- Updates Svelte to the specified version (defaults to "latest")
- Option to update all packages to their latest versions, ignoring semver ranges (`-L` / `--latest`)
- Option to run package updates, integration/e2e tests, and git commands (`git add`, `git commit`, `git push`)
- Option to run in a debug mode
- Option to start at a certain index of project subdirectories
- Displays colored messages for informative progress

## Windows Support

`svelte-next` is a Bash script and does not run natively on Windows. Windows users have three options, in order of recommendation:

### Option 1: WSL (Recommended)

WSL (Windows Subsystem for Linux) runs a full Linux environment inside Windows and is the most compatible option. Install it once from PowerShell (run as Administrator):

```powershell
wsl --install
```

Restart your PC, then open the Ubuntu terminal and install `svelte-next` normally using the Homebrew or manual instructions below. Your Windows project files are accessible inside WSL at `/mnt/c/Users/YourName/...`.

### Option 2: Git Bash

If you already have [Git for Windows](https://git-scm.com/download/win) installed, Git Bash is available immediately — no extra setup needed. Open Git Bash and use the manual (curl or git clone) install methods below.

> **Note:** Git Bash is missing `tput` and some GNU coreutils, so colored output will not display correctly. The script will still run, but expect some visual noise.

### Option 3: Cygwin

[Cygwin](https://www.cygwin.com/) provides a fuller Unix environment than Git Bash but is heavier to install. During Cygwin setup, make sure to select `bash`, `curl`, `git`, and `jq` packages.

---

## Requirements:

- One of the following package managers:
  - pnpm (default if no lock file found)
  - npm
  - yarn
  - bun
- git
- jq

## Package Manager Detection

The script automatically detects which package manager to use based on the lock files present:
- `bun.lockb` → uses Bun
- `pnpm-lock.yaml` → uses pnpm
- `yarn.lock` → uses Yarn
- `package-lock.json` → uses npm
- If no lock file is found → defaults to pnpm

## Installation

### Homebrew (macOS and Linux)

```sh
brew tap shinokada/svelte-next
brew install svelte-next
```

To upgrade to a newer version:

```sh
brew upgrade svelte-next
```

### Debian / Ubuntu (deb)

Download the `.deb` for your architecture from the [releases page](https://github.com/shinokada/svelte-next/releases/latest), then install:

```sh
# amd64
sudo dpkg -i svelte-next_<version>_linux_amd64.deb

# arm64
sudo dpkg -i svelte-next_<version>_linux_arm64.deb
```

To upgrade, download the new `.deb` and run `dpkg -i` again.

### Fedora / RHEL (rpm)

Download the `.rpm` for your architecture from the [releases page](https://github.com/shinokada/svelte-next/releases/latest), then install:

```sh
# amd64
sudo rpm -i svelte-next_<version>_linux_amd64.rpm

# arm64
sudo rpm -i svelte-next_<version>_linux_arm64.rpm
```

To upgrade an existing install, replace `-i` with `-U`.

### awesome package manager

Install [awesome](https://github.com/shinokada/awesome):

```sh
curl -s https://raw.githubusercontent.com/shinokada/awesome/main/install | bash -s install
```

Add the following to your terminal config file, such as `.zshrc` or `.bashrc`:

```sh
export PATH=$HOME/.local/share/bin:$PATH
```

Then source the config file or open a new terminal tab:

```sh
# for example
. ~/.zshrc
```

Install `svelte-next`:

```sh
awesome install shinokada/svelte-next
```

### Manual (tarball)

Download the tarball for your platform from the [releases page](https://github.com/shinokada/svelte-next/releases), then install:

```sh
# Fetch the latest version number automatically
VERSION=$(curl -s https://api.github.com/repos/shinokada/svelte-next/releases/latest | grep '"tag_name"' | cut -d'"' -f4 | cut -c2-)
PLATFORM=darwin_arm64  # darwin_amd64 | linux_amd64 | linux_arm64

curl -sLO "https://github.com/shinokada/svelte-next/releases/download/v${VERSION}/svelte-next_${VERSION}_${PLATFORM}.tar.gz"
tar -xzf "svelte-next_${VERSION}_${PLATFORM}.tar.gz"
cd "svelte-next-${VERSION}"

# Install binary and support files
sudo cp svelte-next /usr/local/bin/svelte-next
sudo chmod +x /usr/local/bin/svelte-next
sudo mkdir -p /usr/local/share/svelte-next
sudo cp -r lib src /usr/local/share/svelte-next/
```

### Manual (git clone)

```sh
git clone https://github.com/shinokada/svelte-next.git
cd svelte-next
sudo cp svelte-next /usr/local/bin/svelte-next
sudo chmod +x /usr/local/bin/svelte-next
sudo mkdir -p /usr/local/share/svelte-next
sudo cp -r lib src /usr/local/share/svelte-next/
```

> **Note:** The `git clone` method keeps `lib/` and `src/` next to the script in the cloned directory, so the manual `cp` to `/usr/local/share/svelte-next` is only needed if you move the binary to `/usr/local/bin`. If you run the script directly from the cloned directory (e.g. `./svelte-next`), no extra steps are needed.

## Usage

```sh
# Install the latest and run package manager update, 
# test:integration and git add, commit, and push 
# if it is a git repo in subdirectories of the CURRENT directory.
svelte-next update .

# run the script in the ./Runes directory
svelte-next update ./Runes

# Use -v param to install a certain Svelte next version.
svelte-next update -v 5.x.x .

# Use -L or --latest to update all packages to their latest versions (ignores semver ranges):
svelte-next update -L .

# Use -p flag to NOT run package updates:
svelte-next update -p .

# Use -s flag to NOT run updating svelte:
svelte-next update -s .

# Use -g flag to NOT run git add, commit, and push:
svelte-next update -g .

# Use -t flag to NOT run integration/e2e tests:
svelte-next update -t .

# Use -f <number> for starting index of subdirectory:
svelte-next update -f 3 .

# Use -d to run in debug mode:
svelte-next update -d .

# Combine the flags
svelte-next update -pg .
svelte-next update -pst .

# To display version: 
svelte-next --version

# To display help:
svelte-next -h | --help
```

## Option list

```
-h --help: Displays help message.
-L --latest: Update all packages to their latest versions (ignores semver ranges).
-s: Skip running updating svelte.
-p: Skip running package updates.
-t: Skip running integration/e2e tests.
-g: Skip running git commands.
-d: Run in debug mode.
-f: Use -f for starting index of subdirectory
-v --version: version
```

## Package Manager Commands Used

The script translates commands appropriately for each package manager:

| Action        | pnpm         | npm                                     | yarn                  | bun                 |
| ------------- | ------------ | --------------------------------------- | --------------------- | ------------------- |
| Install       | pnpm install | npm install                             | yarn add              | bun add             |
| Update        | pnpm update  | npm update                              | yarn upgrade          | bun update          |
| Update latest | pnpm up -L   | npx npm-check-updates -u && npm install | yarn upgrade --latest | bun update --latest |
| Run           | pnpm         | npm                                     | yarn                  | bun                 |

> **Note for npm users:** The `--latest` flag uses [npm-check-updates](https://github.com/raineorshine/npm-check-updates) (`ncu`) since npm has no native equivalent. This will rewrite your `package.json` before reinstalling.

## Note:

- The script automatically detects and uses the appropriate package manager based on lock files
- The script assumes the target directory structure contains project subdirectories where Svelte is installed
- Ensure you have proper permissions to modify files and run git commands in the target directories
- If no lock file is found, the script defaults to using pnpm
