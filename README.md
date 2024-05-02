# svelte-update

`svelte-next` updates svelte next version, run `pnmp update`, git add, commit, and push, `pnpm test:integration` in the subdirectories.

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

svelte-next .

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

