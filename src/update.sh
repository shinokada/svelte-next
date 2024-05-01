fn_update() {

  bannerColor 'Welcome to svelte-next update.' "blue" "*"

  if [[ $# -gt 0  ]];then
    target_dir="$1"
  else
    bannerColor 'Error: Please provide the target directory.' "red" "*"
    exit 1
  fi

  svelte_version="next"

  if [[ "$SVELTE_NEXT" ]];then
    svelte_version="$SVELTE_NEXT"
  fi

  # Check for a minimum of two arguments (version and target directory)
  # if [ $# -lt 2 ]; then
  #   bannerColor '"Error: Please provide the Svelte version and target directory (e.g., updateSvelte i 120 ./my-project).' "red" "*"
  #   exit 1
  # fi

  # Check if the target directory exists and is a directory
  if ! [ -d "$target_dir" ]; then
    bannerColor "Error: Target directory $target_dir does not exist or is not a directory." "red" "*"
    exit 1
  fi

  bannerColor "This script will update svelte version 5 $svelte_version in the $target_dir directories." "blue" "*"

  bannerColor "And run pnpm update, git add, commit, push, pnpm test:integration." "blue" "*"
  
  bannerColor "Use -h or --help for help. " "blue" "*"

  for directory in "$target_dir"/* ; do
    if [ -d "$directory" ]; then 
      cd "$directory"
      bannerColor "Checking $directory" "blue" "*"
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
          bannerColor "pnpm update completed" "green" "*"
        fi
        
        if [[ "$svelte_version" == "next" ]];then
          bannerColor "Running pnpm i -D svelte@$svelte_version ..." "magenta" "*"
          pnpm i -D svelte@next
          bannerColor "pnpm i -D svelte@next completed" "green" "*"
        else
          bannerColor "Running pnpm i -D svelte@5.0.0-next.$svelte_version ..." "magenta" "*"
          pnpm i -D svelte@"5.0.0-next.$svelte_version"
          bannerColor "pnpm i -D svelte@$svelte_version completed" "green" "*"
        fi

        if [[ $FLAG_T == 1 ]];then
          bannerColor "Running pnpm test:integration ..." "magenta" "*"
          pnpm test:integration
          bannerColor "pnpm test:integration completed" "green" "*"
        fi
  
        if [[ -d "$directory/.git" ]] && [[ $FLAG_G == 1 ]]; then
          bannerColor "Running git commands ..." "magenta" "*"
          git add -A && git commit --message "Update Svelte to $svelte_version" && git push
          bannerColor "Git commands completed" "green" "*"
        else
          bannerColor "Skipping git commands" "yellow" "*"
        fi
      else
        bannerColor  "(Not updating) Not able find Svelte version or is likely up-to-date ." "yellow" "*"
      fi
      cd ..
    else
      bannerColor "$directory is not a directory." "red" "*"
    fi
  done

  bannerColor "Whew! Finally done. I'm outta here." "blue" "*" 
  bannerColor "I may be over, but the bugs are eternal. - Some Programmer." "blue" "*"

  # joke_json=$(curl -s https://v2.jokeapi.dev/joke/Programming,Misc?blacklistFlags=nsfw,sexist&type=single)
  # joke_text=$(echo "$joke_json" | jq -r '.joke')

  # # Check if joke extraction was successful (non-empty string)
  # if [[ -n "$joke_text" ]]; then
  #   # Echo the joke
  #   echo "Here's your joke:"
  #   echo "$joke_text"
  # fi
}
