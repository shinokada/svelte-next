# Svelte-Next

## Automate Svelte Version Updates

This script automates updating Svelte versions (specifically targeting next versions) in project directories.

## Features:

- Updates Svelte to the specified version (defaults to "next").
- Runs pnpm update, pnpm test:integration, and optionally git add, commit, and push.
- Displays colored messages for informative progress.
- Supports providing the target directory as a command-line argument.

## Requirements:

- pnpm package manager
- git
- jq

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
# Install the latest and run pnpm update, 
# pnpm test:integration and git add, commit, and push 
# if it is a git repo in subdirectories of the CURRENT directory.
svelte-next update .

# run the script in the ./Runes directory
svelte-next update ./Runes

# Use -v param to install a certain Svelte next version.
svelte-next update -v 120 .

# Use -p flag to NOT to run pnpm update:
svelte-next update -p .

# Use -s flag to NOT to run updating svelte:
svelte-next update -s .

# Use -g flag to NOT to run git add, commit, and push:
svelte-next update -g .

# Use -t flag to NOT to run pnpm test:integration:
svelte-next update -t .

# Use -f <number> for starting index of subdirectory:
svelte-next update -f 3

# Combine the flags
svelte-next update -pg .
svelte-next update -pst .

# To display version: 
svelte-next --version

# To display help:
svelte-next -h | --help
```

## Optional Flags:

```
-h or --help: Displays help message.
-s: Skip running updating svelte.
-p: Skip running pnpm update.
-t: Skip running pnpm test:integration.
-g: Skip running git commands.
-f: Use -f for starting index of subdirectory
```

## Note:

- The script assumes the target directory structure contains project subdirectories where Svelte is installed.
- Ensure you have proper permissions to modify files and run git commands in the target directories.