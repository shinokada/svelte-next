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

  # Check if the target directory exists and is a directory
  if ! [ -d "$target_dir" ]; then
    bannerColor "Error: Target directory $target_dir does not exist or is not a directory." "red" "*"
    exit 1
  fi

  bannerColor "This script will run the following tasks:" "blue" "*"

  if [[ $FLAG_P == 1 ]];then
    bannerColor "pnpm update" "blue" "*"
  fi
  if [[ $FLAG_S == 1 ]];then
    bannerColor "pnpm i -D svelte@$svelte_version" "blue" "*"
  fi
  if [[ $FLAG_T == 1 ]];then
    bannerColor "pnpm test:integration" "blue" "*"
  fi
  if [[ $FLAG_G == 1 ]];then
    bannerColor "git add, commit, and push" "blue" "*"
  fi
  
  bannerColor "Use -h or --help for help. " "blue" "*"

  for directory in "$target_dir"/* ; do
    if [[ -d "$directory" && -f "$directory/package.json" && $(grep -q '"svelte":' "$directory/package.json" && echo $? ) ]]; then
      cd "$directory"
      bannerColor "Checking $directory" "blue" "*"
      # Get current Svelte version
      current_version=$(pnpm list svelte --depth=0 | tail -n 1)
      bannerColor "Your current Svelte version is: $current_version" "green" "*"

      if [[ "$current_version" =~ "next" ]]; then

        if [[ $FLAG_P == 1 ]];then
          bannerColor "Running pnpm update in $directory ..." "magenta" "*"
          pnpm update
          bannerColor "pnpm update completed" "green" "*"
        else
          bannerColor "Skipping git pnpm update." "yellow" "*"
        fi
        
        if [[ $FLAG_S == 1 ]];then
          if [[ "$svelte_version" == "next" ]];then
            bannerColor "Running pnpm i -D svelte@$svelte_version ..." "magenta" "*"
            pnpm i -D svelte@next
            bannerColor "pnpm i -D svelte@next completed" "green" "*"
          else
            bannerColor "Running pnpm i -D svelte@5.0.0-next.$svelte_version ..." "magenta" "*"
            pnpm i -D svelte@"5.0.0-next.$svelte_version"
            bannerColor "pnpm i -D svelte@$svelte_version completed" "green" "*"
          fi
        else
          bannerColor "Skipping updating svelte." "yellow" "*"
        fi

        if [[ $FLAG_T == 1 ]];then
          bannerColor "Running pnpm test:integration ..." "magenta" "*"
          pnpm test:integration
          bannerColor "pnpm test:integration completed" "green" "*"
        else
          bannerColor "Skipping pnpm test:integration." "yellow" "*"
        fi
  
        if [[ -d "./.git" ]] && [[ $FLAG_G == 1 ]]; then
          bannerColor "Running git commands ..." "magenta" "*"
          git add -A && git commit --message "Update Svelte to $svelte_version" && git push
          bannerColor "Git commands completed" "green" "*"
        else
          bannerColor "Skipping git commands" "yellow" "*"
        fi

      else
        bannerColor  "Your subdirectory $directory does not have Svelte version. (Not updating)." "yellow" "*"
      fi
      cd ..
    else
      bannerColor "The directory $directory either doesn't exist, doesn't have a package.json, or Svelte isn't mentioned." "red" "*"
    fi
  done

  bannerColor "Whew! Finally done. I'm outta here." "blue" "*" 
  QUOTE=$(curl -s https://api.quotable.io/quotes/random | jq -r '.[0].content + " - " + .[0].author')

  if [[ -n "$QUOTE" ]]; then
    bannerColor "$QUOTE" "blue" "*" 
  fi


  # joke_json=$(curl -s https://v2.jokeapi.dev/joke/Programming,Misc?blacklistFlags=nsfw,sexist&type=single)
  # joke_text=$(echo "$joke_json" | jq -r '.joke')

  # # Check if joke extraction was successful (non-empty string)
  # if [[ -n "$joke_text" ]]; then
  #   # Echo the joke
  #   echo "Here's your joke:"
  #   echo "$joke_text"
  # fi
}
