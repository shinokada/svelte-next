fn_update() {

  # Get the desired Svelte version from the script argument
  bannerColor 'Welcome to svelte update.' "blue" "*"

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

  bannerColor "This script will update svelte version, run pnpm update, git add, commit, and push, pnpm test:integration in $target_dir to $svelte_version. Use -h or --help for help." "blue" "*"


  for directory in "$target_dir"/* ; do
    bannerColor "Checking $directory" "blue" "*"

    if [ -d "$directory" -a -d "$directory/.git" ]; then
      cd "$directory"

      # Get current Svelte version
      current_version=$(pnpm list svelte --depth=0 | awk '{print $2}')
      update_needed=false

      if [[ "$current_version" =~ "next" ]]; then
        # If current_version has "next", check version comparison
        if [[ "$(semver compare "$svelte_version" "$current_version")" -lt "1" ]]; then
          update_needed=true 
        fi
      fi

      if [[ $update_needed == true ]]; then

        if [[ $FLAG_P == 1 ]];then
          bannerColor "Running pnpm update in $directory ..." "magenta" "*"
          pnpm update
          bannerColor "Update completed" "green" "*"
        fi
        
        
        bannerColor "Running pnpm i -D svelte@$svelte_version ..." "magenta" "*"
        pnpm i -D svelte@"5.0.0-next.$svelte_version"
        bannerColor "Update completed" "green" "*"
        
        if [[ $FLAG_T == 1 ]];then
          bannerColor "Running pnpm test:integration ..." "magenta" "*"
          pnpm test:integration
          bannerColor "Test completed" "green" "*"
        fi
  
        if [[ $FLAG_G == 1 ]];then
          bannerColor "Running git commands ..." "magenta" "*"
          git add -A && git commit --message "Update Svelte to $svelte_version" && git push
          bannerColor "Git commands completed" "green" "*"
        fi
        
      else
        bannerColor  "Svelte version $current_version (next) is likely up-to-date (not updating)." "blue" "*"
      fi
      cd ..
    elif [ -d "$directory" ]; then
      bannerColor "Directory $directory exists, but it lacks a .git directory." "red" "*"
    else
      bannerColor "Directory $directory does not exist." "red" "*"
    fi
  done
}
