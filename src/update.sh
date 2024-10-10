fn_update() {

  if [[ $# -gt 0  ]];then
    target_dir=$(realpath "$1")
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

  messages=()
  messages+=("Welcome to svelte-next update. Use -h or --help for help.")
  messages+=("This script will run the following tasks:")
  messages+=("")

  if [[ $FLAG_P == 1 ]]; then
      messages+=("- pnpm update")
  fi

  if [[ $FLAG_S == 1 ]]; then
      messages+=("- pnpm i -D svelte@$svelte_version")
  fi

  if [[ $FLAG_T == 1 ]]; then
      messages+=("- pnpm test:integration")
  fi

  if [[ $FLAG_G == 1 ]]; then
      messages+=("- git add, commit, and push")
  fi

  if [[ $FROM ]]; then
      messages+=("- Starting from index $FROM")
  fi

  # Join all messages with newlines
  formatted_message=$(printf "%s\n" "${messages[@]}")

  # Output all messages at once using bannerColor
  newBannerColor "$formatted_message" "blue" "*" 50

  count=0
  for directory in "$target_dir"/* ; do
    if [[ $FROM -ge 1 ]] && (( count < FROM )); then
      ((count++))  # Increment count for skipped directories
      continue
    fi
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
          # get the new next version installed
          new_version=$(pnpm list svelte --depth=0 | tail -n 1)
          bannerColor "Running git commands ..." "magenta" "*"
          git add -A && git commit --message "Update Svelte to $new_version" && git push origin $(git branch --show-current)
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
  # https://api.quotable.io/quotes/random is down right now
  QUOTE=$(curl -s https://quoteslate.vercel.app/api/quotes/random | jq -r '.[0].content + " - " + .[0].author')

  if [[ -n "$QUOTE" ]]; then
    echo "Here's a quote:"
    newBannerColor "$QUOTE" "blue" "*" 30
  fi


  # joke_json=$(curl -s https://v2.jokeapi.dev/joke/Programming,Misc?blacklistFlags=nsfw,sexist&type=single)
  # joke_text=$(echo "$joke_json" | jq -r '.joke')

  # # Check if joke extraction was successful (non-empty string)
  # if [[ -n "$joke_text" ]]; then
  #   # Echo the joke
  #   echo "Here's a programming joke:"
  #   echo "$joke_text"
  # fi
}
