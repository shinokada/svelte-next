# svelte-update: Automate Svelte Version Updates

This script automates updating Svelte versions (specifically targeting next versions) in project directories.

## Features:

- Updates Svelte to the specified version (defaults to "next").
- Runs pnpm update, pnpm test:integration, and optionally git add, commit, and push.
- Displays colored messages for informative progress.
- Supports providing the target directory as a command-line argument.

## Requirements:

- pnpm package manager
- git version control system (optional - for git commands)

## Installation

Install [awesome package manager](https://github.com/shinokada/awesome)

```sh
curl -s https://raw.githubusercontent.com/shinokada/awesome/main/install | bash -s install
```

Add the following to your terminal config file, such as .zshrc or .bashrc.

```sh
export PATH=$HOME/.local/share/bin:$PATH
```

Then source the config file or open a new terminal tab.

```sh
# for example
. ~/.zshrc
```

Install `svelte-next`:

```sh
awesome install shinokada/svelte-next
```

## Usage

```sh
# Install the latest and run pnpm update, pnpm test:integration and git add, commit, and push if it is a git repo in subdirectories of the current directory

svelte-next update .

# Use -v param to install a certain Svelte next version.

svelte-next update -v 120 .

# Use -p flag to NOT to run pnpm update:

svelte-next update -p .

# Use -g flag to NOT to run git add, commit, and push:
svelte-next update -g .

# Use -t flag to NOT to run pnpm test:integration:
svelte-next update -t .

# Combine the flags
svelte-next update -pg .

# To display version: 
svelte-next --version

# To display help:
svelte-next -h | --help
```

## Optional Flags:

```
-h or --help: Displays help message.
-p: Skip running pnpm update (default: runs).
-t: Skip running pnpm test:integration (default: runs).
-g: Skip running git commands (default: runs if git repo present).
```

## Note:

- The script assumes the target directory structure contains project subdirectories where Svelte is installed.
- Ensure you have proper permissions to modify files and run git commands in the target directories.