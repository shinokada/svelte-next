fn_update() {

  if [[ $# -gt 0  ]];then
    target_dir=$(realpath "$1")
  else
    bannerColor 'Error: Please provide the target directory.' "red" "*"
    exit 1
  fi

  svelte_version="latest"

  if [[ "$SVELTE_NEXT" ]]; then
    svelte_version="$SVELTE_NEXT"
  fi

  # Debug output
  if [[ $DEBUG == 1 ]]; then
    echo ""
    echo "Debug: FROM=$FROM"
    echo "Debug: target_dir=$target_dir"
    echo ""
  fi

  # Check if the target directory exists and is a directory
  if ! [ -d "$target_dir" ]; then
    bannerColor "Error: Target directory $target_dir does not exist or is not a directory." "red" "*"
    exit 1
  fi

  # Count the number of directories
  dir_count=$(find "$target_dir" -maxdepth 1 -type d | wc -l)
  dir_count=$((dir_count - 1))  # Subtract 1 to exclude the target directory itself

  newBannerColor "Total directories found: $dir_count" "blue" "*"

  # Check if FROM is set and valid
  if [[ -n $FROM ]] && (( FROM >= dir_count )); then
    newBannerColor "Error: FROM value ($FROM) is greater than or equal to the number of directories ($dir_count)." "red" "*"
    newBannerColor "Please choose a FROM value less than $dir_count." "yellow" "*"
    exit 1
  fi

  messages=()
  messages+=("🚀 Welcome to svelte-next update. Use -h or --help for help.")
  messages+=("This script will run the following tasks:")
  messages+=("")

  if [[ $FLAG_P == 1 ]]; then
      messages+=("⚡ pnpm update")
  fi

  if [[ $FLAG_S == 1 ]]; then
      messages+=("⚡ pnpm i -D svelte@$svelte_version")
  fi

  if [[ $FLAG_T == 1 ]]; then
      messages+=("⚡ pnpm test:integration")
  fi

  if [[ $FLAG_G == 1 ]]; then
      messages+=("⚡ git add, commit, and push")
  fi

  if [[ $FROM ]]; then
    FROM=$((FROM))
    messages+=("⚡ Starting from index $FROM")
  fi

  # Join all messages with newlines
  formatted_message=$(printf "%s\n" "${messages[@]}")

  # Output all messages at once using bannerColor
  newBannerColor "$formatted_message" "blue" "*"

  # count=0
  cd "$target_dir" || exit
  directories=($(ls -d */))

  if [[ $DEBUG == 1 ]];then
    echo ""
    for ((i=0; i<${#directories[@]}; i++)); do
      echo "Index $i: ${directories[i]}"
    done
    echo "Debug: Loop starting from index $FROM"
    echo ""
  fi

  for ((i=FROM; i<${#directories[@]}; i++)); do
    cd "$target_dir/${directories[$i]}" || exit
    current_dir_name=$(basename "$(pwd)")
    newBannerColor "Started processing $current_dir_name" "green" "*"
    
    if [[ $DEBUG == 1 ]]; then
      echo ""
      echo "Debug: Checking $target_dir/$current_dir_name"
      if [[ -f "$target_dir/$current_dir_name/package.json" ]] then
        echo "Debug: $target_dir/$current_dir_name/package.json exists."
      else
        echo "Debug: $target_dir/$current_dir_name/package.json does not exist."
      fi
      if grep -q '"svelte":' "$target_dir/$current_dir_name/package.json"; then
        echo "Debug: $target_dir/$current_dir_name/package.json contains 'svelte'."
      else
        echo "Debug: $target_dir/$current_dir_name/package.json does not contain 'svelte'."
      fi
      echo ""
    fi

    if [[ -f "$target_dir/$current_dir_name/package.json" ]] && grep -q '"svelte":' "$target_dir/$current_dir_name/package.json"; then

      newBannerColor "🚀 Checking $current_dir_name" "blue" "*"
      # Get current Svelte version
      current_version=$(pnpm list svelte --depth=0 | tail -n 1)
      version_number=$(echo "$current_version" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+(-next\.[0-9]+)?')

      newBannerColor "Your current Svelte version is: $current_version" "green" "*"

      if [[ $DEBUG == 1 ]]; then
        echo ""
        echo "Debug: Full version string: '$current_version'"
        echo "Debug: Extracted version number: '$version_number'"

      if [[ "$version_number" =~ ^5\.0\.0(-next\.[0-9]+)?$ ]] || [[ "$version_number" =~ ^5\.[0-9]+\.[0-9]+$ ]]; then
          echo "Debug: $current_version is a valid Svelte version"
        else
          echo "Debug: $current_version is not a valid Svelte version"
        fi
        echo ""
      fi

      if [[ "$version_number" =~ ^5\.0\.0(-next\.[0-9]+)?$ ]] || [[ "$version_number" =~ ^5\.[0-9]+\.[0-9]+$ ]]; then

        if [[ $DEBUG == 1 ]]; then
          echo ""
          echo "Debug: Working on $current_dir_name"
          echo ""
        fi

        if [[ $FLAG_P == 1 ]];then
          newBannerColor "🔄 Running pnpm update in $current_dir_name ..." "magenta" "*" 
          pnpm update
          newBannerColor "👍 pnpm update completed" "green" "*" 
        else
          newBannerColor "⏭️  Skipping pnpm update." "yellow" "*"
        fi
        
        if [[ $FLAG_S == 1 ]];then
          if [[ "$svelte_version" == "latest" ]];then
            newBannerColor "🏃 Running pnpm i -D svelte@$svelte_version ..." "magenta" "*"
            pnpm i -D svelte@latest
            newBannerColor "🚀 pnpm i -D svelte@latest completed" "green" "*"
          else
            newBannerColor "🏃 Running pnpm i -D svelte@$svelte_version ..." "magenta" "*"
            pnpm i -D svelte@"$svelte_version"
            newBannerColor "🚀 pnpm i -D svelte@$svelte_version completed" "green" "*"
          fi
        else
          newBannerColor "⏭️  Skipping updating svelte." "yellow" "*"
        fi

        if [[ $FLAG_T == 1 ]];then
          newBannerColor "🏃 Running pnpm test:integration ..." "magenta" "*"
          pnpm test:integration
          newBannerColor "🚀 pnpm test:integration completed" "green" "*"
        else
          newBannerColor "⏭️  Skipping pnpm test:integration." "yellow" "*"
        fi
  
        if [[ -d "./.git" ]] && [[ $FLAG_G == 1 ]]; then
          # get the current version installed
          new_version=$(pnpm list svelte --depth=0 | tail -n 1)
          newBannerColor "🏃 Running git commands ..." "magenta" "*"
          git add -A && git commit --message "Update Svelte to $new_version" && git push origin $(git branch --show-current)
          newBannerColor "🚀 Git commands completed" "green" "*"
        else
          newBannerColor "⏭️  Skipping git commands" "yellow" "*"
        fi

      else
        newBannerColor  "Skipping $current_dir_name: No package.json or no Svelte dependency" "yellow" "*"
      fi
      cd "$target_dir" || { echo "Failed to return to $target_dir"; exit 1; }
    else
      newBannerColor "😥 Skipping $current_dir_name: No package.json or no Svelte dependency" "red" "*" 50
    fi
    newBannerColor "Finished processing $current_dir_name. Moving to next." "green" "*"
  done

  newBannerColor "👍 Whew! Finally done. I'm outta here." "blue" "*" 


  # https://api.quotable.io/quotes/random is down right now
  # QUOTE=$(curl -s https://api.quotable.io/quotes/random | jq -r '.[0].content + " - " + .[0].author')
  # At https://quoteslate.vercel.app/api/quotes/randomhe, the JSON output is a single object, not an array
  QUOTE=$(curl -s https://quoteslate.vercel.app/api/quotes/random | jq -r '.quote + " - " + .author')

  if [[ -n "$QUOTE" ]]; then
    echo -e "Here's a quote for you $random_emoji"
    newBannerColor "$QUOTE" "green" "*"
  fi

  # joke_json=$(curl -H "Accept: application/json" https://icanhazdadjoke.com/)
  # joke_text=$(echo "$joke_json" | jq -r '.joke')

  # # Check if joke extraction was successful (non-empty string)
  # if [[ -n "$joke_text" ]]; then
  #   # Echo the joke
  #   echo "Here's a programming joke:"
  #   echo "$joke_text"
  # fi
}
