# Svelte-Next

## Automate Svelte Version Updates

This script automates updating Svelte versions in project directories. If you have mutiple Svelte projects in subdirectories, this script will update them all.

## Features:

- Updates Svelte to the specified version (defaults to "latest").
- Option to run `pnpm update`, `pnpm test:integration`, and git commands (`git add`, `git commit`, `git push`).
- Option to run in a debug mode.
- Option to start a certain index of project subdirectories.
- Displays colored messages for informative progress.

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
svelte-next update -v 5.x.x .

# Use -p flag to NOT to run pnpm update:
svelte-next update -p .

# Use -s flag to NOT to run updating svelte:
svelte-next update -s .

# Use -g flag to NOT to run git add, commit, and push:
svelte-next update -g .

# Use -t flag to NOT to run pnpm test:integration:
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
-s: Skip running updating svelte.
-p: Skip running pnpm update.
-t: Skip running pnpm test:integration.
-g: Skip running git commands.
-d: Run in debug mode.
-f: Use -f for starting index of subdirectory
-v --version: version
```

## Note:

- The script assumes the target directory structure contains project subdirectories where Svelte is installed.
- Ensure you have proper permissions to modify files and run git commands in the target directories.