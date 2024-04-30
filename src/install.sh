fn_install() {
    # echo "FLAG_A: $FLAG_A"
    # echo "FLAG_B: $FLAG_B"
    # echo "FLAG_F: $FLAG_F"
    # echo "FLAG_W: $FLAG_W"
    # echo "VERBOSE: $VERBOSE"
    # echo "PARAM: $PARAM"
    # echo "NUMBER: $NUMBER"
    # echo "OPTION: $OPTION"
    # echo "VERSION: $VERSION"

    # echo "My variable VAR1: $VAR1."

    # i=0
    # while [ $# -gt 0 ] && i=$((i + 1)); do
    #     echo "$i: $1"
    #     shift
    # done

    # Get the desired Svelte version from the script argument
 bannerColor 'Welcome to updateSvelte.' "blue" "*"

svelte_version="$1"

# Check if an argument is provided
if [ -z "$svelte_version" ]; then
  echo "Error: Please provide the Svelte version (e.g., updateSvelte 120)"
  exit 1
fi

# Check for a minimum of two arguments (version and target directory)
if [ $# -lt 2 ]; then
  echo "Error: Please provide the Svelte version and target directory (e.g., my-script 120 ./my-project)."
  exit 1
fi

# Get the desired Svelte version (shift argument list to remove the first argument)
svelte_version="$1"
shift

# Get the target directory (remaining arguments after shift)
target_dir="$*"

# Check if the target directory exists and is a directory
if ! [ -d "$target_dir" ]; then
  echo "Error: Target directory '$target_dir' does not exist or is not a directory."
  exit 1
fi

 bannerColor "This script update svelte version in $target_dir to $svelte_version." "blue" "*"

# Loop through each subdirectory within the target directory
for directory in "$target_dir"/* ; do
  # Check if it's a directory (avoid hidden directories)
  if [ -d "$directory" ] && [[ ! $directory =~ ^\./\.* ]]; then
    # Change directory to the current subdirectory
    cd "$directory"

    echo "Running pnpm update ..."
    # Run pnpm update (assuming you're using pnpm)
    pnpm update

    # Get current Svelte version
    current_version=$(pnpm list svelte --depth=0 | awk '{print $2}')
    # Check if update is needed (desired version > current version)
    if [[ "$(semver compare "$svelte_version" "$current_version")" == "1" ]]; then

      echo "Running pnpm i -D svelte@$svelte_version ..."

      # Install the desired Svelte version as a dev dependency
      pnpm i -D svelte@"5.0.0-next.$svelte_version"
      
    else
      echo "Svelte version $current_version in $directory is already up-to-date."
    fi
    
    echo "Running pnpm test:integration ..."
    # Run integration tests
    pnpm test:integration

    if [ -d "$directory/.git" ]; then
      echo "Running git add ..."
      git add -A

      echo "Running git commit ..."
      git commit --message "Update Svelte to $svelte_version"

      echo "Running git push ..."
      git push
    fi
    cd ..
  fi
done

echo "Svelte updated to 5.0.0-next.$svelte_version and integration tests run in subdirectories."
}
