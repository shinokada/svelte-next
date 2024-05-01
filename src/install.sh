fn_install() {

  # Get the desired Svelte version from the script argument
  bannerColor 'Welcome to updateSvelte.' "blue" "*"

  svelte_version="$1"

  # Check if an argument is provided
  if [ -z "$svelte_version" ]; then
    bannerColor '"Error: Please provide the Svelte version (e.g., updateSvelte i 120)' "red" "*"
    exit 1
  fi

  # Check for a minimum of two arguments (version and target directory)
  if [ $# -lt 2 ]; then
    bannerColor '"Error: Please provide the Svelte version and target directory (e.g., updateSvelte i 120 ./my-project).' "red" "*"
    exit 1
  fi

  # Get the desired Svelte version (shift argument list to remove the first argument)
  svelte_version="$1"
  shift

  # Get the target directory (remaining arguments after shift)
  target_dir="$*"

  # Check if the target directory exists and is a directory
  if ! [ -d "$target_dir" ]; then
    bannerColor "Error: Target directory $target_dir does not exist or is not a directory." "red" "*"
    exit 1
  fi

  bannerColor "This script update svelte version in $target_dir to $svelte_version." "blue" "*"

  # Loop through each subdirectory within the target directory
  for directory in "$target_dir"/* ; do
    # Check if it's a directory (avoid hidden directories)
    bannerColor "Checking $directory" "blue" "*"

    if [ -d "$directory" -a -d "$directory/.git" ]; then
      cd "$directory"
      bannerColor "Running pnpm update in $directory ..." "magenta" "*"
      # Run pnpm update (assuming you're using pnpm)
      pnpm update

      # Get current Svelte version
      current_version=$(pnpm list svelte --depth=0 | awk '{print $2}')
      # Check for update condition
      update_needed=false  # Initialize a flag for update

      if [[ "$current_version" =~ "next" ]]; then
        # If current_version has "next", check version comparison
        if [[ "$(semver compare "$svelte_version" "$current_version")" -lt "1" ]]; then
          update_needed=true  # Set flag if both conditions met
        fi
      fi
      # Check if update is needed (desired version > current version)
      if [[ $update_needed == true ]]; then

        bannerColor "Running pnpm i -D svelte@$svelte_version ..." "magenta" "*"

        # Install the desired Svelte version as a dev dependency
        pnpm i -D svelte@"5.0.0-next.$svelte_version"

        bannerColor "Running pnpm test:integration ..." "magenta" "*"
        # Run integration tests
        pnpm test:integration

        bannerColor "Running git commands ..." "magenta" "*"
        git add -A && git commit --message "Update Svelte to $svelte_version" && git push
      else
        bannerColor  "Svelte version $current_version (next) is likely up-to-date (not updating)." "blue" "*"
      fi
      cd ..
    elif [ -d "$directory" ]; then
      echo "Directory $directory exists, but it lacks a .git directory."
    else
      echo "Directory $directory does not exist."
    fi
  done
}
