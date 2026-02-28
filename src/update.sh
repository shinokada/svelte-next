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

  # FLAG_L=1 means -L was passed (run latest update)
  # FLAG_P/S/T/G use inverted logic: 1=run (flag not passed), 0=skip (flag passed)
  if [[ $FLAG_L == 1 ]]; then
      messages+=("⚡ Update all packages to latest (ignoring semver ranges)")
  fi

  if [[ $FLAG_L == 1 ]]; then
      messages+=("⏭️  Skipping package manager update (-L takes precedence)")
  elif [[ $FLAG_P != 0 ]]; then
      messages+=("⚡ Package manager update")
  else
      messages+=("⏭️  Skipping package manager update (-p)")
  fi

  if [[ $FLAG_S == 1 ]]; then
      messages+=("⚡ Install svelte@\"^$svelte_version\"")
  else
      messages+=("⏭️  Skipping svelte install (-s)")
  fi

  if [[ $FLAG_T == 1 ]]; then
      messages+=("⚡ Run integration/e2e tests")
  else
      messages+=("⏭️  Skipping tests (-t)")
  fi

  if [[ $FLAG_G == 1 ]]; then
      messages+=("⚡ git add, commit, and push")
  else
      messages+=("⏭️  Skipping git commands (-g)")
  fi

  if [[ $FROM -gt 0 ]]; then
    messages+=("⚡ Starting from index $FROM")
  fi

  # Join all messages with newlines
  formatted_message=$(printf "%s\n" "${messages[@]}")

  # Output all messages at once using bannerColor
  newBannerColor "$formatted_message" "blue" "*"

  cd "$target_dir" || exit
  directories=()
  for d in */; do
    [[ -d "$d" ]] && directories+=("$d")
  done

  # Function to detect package manager
  detect_package_manager() {
    local dir="$1"
    if [[ -f "$dir/bun.lockb" ]]; then
      echo "bun"
    elif [[ -f "$dir/pnpm-lock.yaml" ]]; then
      echo "pnpm"
    elif [[ -f "$dir/yarn.lock" ]]; then
      echo "yarn"
    elif [[ -f "$dir/package-lock.json" ]]; then
      echo "npm"
    else
      echo "pnpm"  # Default to pnpm if no lock file is found
    fi
  }

  # Check that the detected package manager is actually installed
  check_package_manager() {
    local pkg_manager="$1"
    if ! command -v "$pkg_manager" &>/dev/null; then
      newBannerColor "Error: $pkg_manager is required for this project but is not installed." "red" "*"
      exit 1
    fi
    if [[ "$pkg_manager" == "npm" && $FLAG_L == 1 ]] && ! command -v npx &>/dev/null; then
      newBannerColor "Error: npx is required for npm latest updates but is not installed." "red" "*"
      exit 1
    fi
  }

  # Read the installed version directly from package.json using jq.
  # This avoids shelling out to the package manager (which may hit the
  # network or hang) just to read a version number we already have on disk.
  get_package_version() {
    local package_name="$1"
    local package_json="$2"
    local version
    version=$(jq -r --arg pkg "$package_name" '(.dependencies[$pkg] // .devDependencies[$pkg] // .peerDependencies[$pkg] // .optionalDependencies[$pkg]) // empty | ltrimstr("^") | ltrimstr("~")' "$package_json")
    echo "${version:-}"
  }

  # Function to run package manager commands
  run_pkg_cmd() {
    local cmd="$1"
    local pkg_manager="$2"
    local args="${3:-}"  # Make args optional with empty default
    
    case "$pkg_manager" in
      "bun")
        case "$cmd" in
          "install") bun add ${args:-} ;;  # Use ${args:-} to handle empty args
          "update") bun update $args ;;
          "update-latest") bun update --latest ;;
          "run") bun $args ;;
        esac
        ;;
      "pnpm")
        case "$cmd" in
          "install") pnpm install ${args:-} ;;
          "update") pnpm update $args ;;
          "update-latest") pnpm up -L ;;
          "run") pnpm $args ;;
        esac
        ;;
      "yarn")
        case "$cmd" in
          "install") yarn add ${args:-} ;;
          "update") yarn upgrade $args ;;
          "update-latest")
            # Detect Yarn major version: Berry (v2+) uses 'yarn up *';
            # Classic (v1) uses 'yarn upgrade --latest'.
            local yarn_major
            yarn_major=$(yarn --version 2>/dev/null | cut -d'.' -f1)
            if [[ "$yarn_major" =~ ^[0-9]+$ ]] && (( yarn_major >= 2 )); then
              yarn up '*'
            else
              yarn upgrade --latest
            fi
            ;;
          "run") yarn $args ;;
        esac
        ;;
      "npm")
        case "$cmd" in
          "install") npm install ${args:-} ;;
          "update") npm update $args ;;
          "update-latest") npx --yes npm-check-updates -u && npm install ;;
          "run") npm $args ;;
        esac
        ;;
    esac
  }


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
      if [[ -f "$target_dir/$current_dir_name/package.json" ]]; then
        echo "Debug: $target_dir/$current_dir_name/package.json exists."
      else
        echo "Debug: $target_dir/$current_dir_name/package.json does not exist."
      fi
      if jq -e '.dependencies.svelte // .devDependencies.svelte // .peerDependencies.svelte // .optionalDependencies.svelte' "$target_dir/$current_dir_name/package.json" &>/dev/null; then
        echo "Debug: $target_dir/$current_dir_name/package.json contains 'svelte'."
      else
        echo "Debug: $target_dir/$current_dir_name/package.json does not contain 'svelte'."
      fi
      echo ""
    fi

    if [[ -f "$target_dir/$current_dir_name/package.json" ]] && jq -e '.dependencies.svelte // .devDependencies.svelte // .peerDependencies.svelte // .optionalDependencies.svelte' "$target_dir/$current_dir_name/package.json" &>/dev/null; then
      # Detect package manager for current directory and verify it is installed
      pkg_manager=$(detect_package_manager "$target_dir/$current_dir_name")
      check_package_manager "$pkg_manager"

      # Calculate current position (i + 1 since array is 0-based)
      current_pos=$((i + 1))
      # Get total number of directories
      total_dirs=${#directories[@]}
      newBannerColor "🚀 Checking $current_dir_name ($current_pos/$total_dirs) using $pkg_manager" "blue" "*"
      
      # Get current Svelte version directly from package.json (no network call)
      local pkg_json="$target_dir/$current_dir_name/package.json"
      current_version=$(get_package_version "svelte" "$pkg_json")
      version_number=$(echo "$current_version" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+(-next\.[0-9]+)?')

      # Extract major version from the validated semver string (works for Svelte 5, 6, 7 ...)
      svelte_major=$(echo "$version_number" | cut -d'.' -f1)

      newBannerColor "Your current Svelte version is: $current_version" "green" "*"

      if [[ $DEBUG == 1 ]]; then
        echo ""
        echo "Debug: Full version string: '$current_version'"
        echo "Debug: Extracted version number: '$version_number'"
        echo "Debug: Svelte major version: '$svelte_major'"
        echo "Debug: Using package manager: '$pkg_manager'"

        if [[ "$svelte_major" =~ ^[0-9]+$ ]] && (( svelte_major >= 5 )); then
          echo "Debug: $current_version is a supported Svelte version (major >= 5)"
        else
          echo "Debug: $current_version is not a supported Svelte version (major < 5)"
        fi
        echo ""
      fi

      if [[ "$svelte_major" =~ ^[0-9]+$ ]] && (( svelte_major >= 5 )); then

        if [[ $DEBUG == 1 ]]; then
          echo ""
          echo "Debug: Working on $current_dir_name"
          echo ""
        fi

        if [[ $FLAG_L == 1 ]]; then
          newBannerColor "🔄 Running update --latest in $current_dir_name using $pkg_manager ..." "magenta" "*"
          run_pkg_cmd "update-latest" "$pkg_manager"
          newBannerColor "👍 update --latest completed" "green" "*"
        elif [[ $FLAG_P == 1 ]]; then
          newBannerColor "🔄 Running $pkg_manager update in $current_dir_name ..." "magenta" "*" 
          run_pkg_cmd "update" "$pkg_manager"
          newBannerColor "👍 $pkg_manager update completed" "green" "*" 
        else
          newBannerColor "⏭️  Skipping package manager update." "yellow" "*"
        fi
        
        if [[ $FLAG_S == 1 ]];then
          if [[ "$svelte_version" == "latest" ]];then
            newBannerColor "🏃 Installing svelte@$svelte_version using $pkg_manager ..." "magenta" "*"
            run_pkg_cmd "install" "$pkg_manager" "-D svelte@latest"
          else
            newBannerColor "🏃 Installing svelte@$svelte_version using $pkg_manager ..." "magenta" "*"
            run_pkg_cmd "install" "$pkg_manager" "-D svelte@$svelte_version"
          fi
          newBannerColor "🚀 svelte installation completed" "green" "*"
        else
          newBannerColor "⏭️  Skipping updating svelte." "yellow" "*"
        fi

        if [[ $FLAG_T == 1 ]]; then
          # Check if package.json has "test:integration" or "test:e2e" scripts
          if grep -q '"test:integration": "playwright test"' "$target_dir/$current_dir_name/package.json"; then
            newBannerColor "🏃 Running test:integration ..." "magenta" "*"
            run_pkg_cmd "run" "$pkg_manager" "test:integration"
            newBannerColor "🚀 test:integration completed" "green" "*"
          elif grep -q '"test:e2e": "playwright test"' "$target_dir/$current_dir_name/package.json"; then
            newBannerColor "🏃 Running test:e2e ..." "magenta" "*"
            run_pkg_cmd "run" "$pkg_manager" "test:e2e"
            newBannerColor "🚀 test:e2e completed" "green" "*"
          else
            newBannerColor "⏭️  No compatible test script found in package.json." "yellow" "*"
          fi
        else
          newBannerColor "⏭️  Skipping test." "yellow" "*"
        fi
  
        if [[ -d "./.git" ]] && [[ $FLAG_G == 1 ]]; then
          # get the post-update version from package.json
          new_version=$(get_package_version "svelte" "$pkg_json")
          newBannerColor "🏃 Running git commands ..." "magenta" "*"
          git add -A
          if git diff --cached --quiet; then
            newBannerColor "ℹ️  No changes to commit" "yellow" "*"
          else
            git commit --message "Update Svelte to $new_version" && git push origin "$(git branch --show-current)"
          fi
          newBannerColor "🚀 Git commands completed" "green" "*"
        else
          newBannerColor "⏭️  Skipping git commands" "yellow" "*"
        fi

      else
        newBannerColor "Skipping $current_dir_name: Svelte major version ($svelte_major) is less than 5" "yellow" "*"
      fi
      cd "$target_dir" || { echo "Failed to return to $target_dir"; exit 1; }
    else
      newBannerColor "😥 Skipping $current_dir_name: No package.json or no Svelte dependency" "red" "*" 50
    fi
    newBannerColor "Finished processing $current_dir_name. Moving to next." "green" "*"
  done

  # Alternative API endpoints (try these if rate-limited)
  APIS=(
    "https://zenquotes.io/api/random"
    "https://quoteslate.vercel.app/api/quotes/random"
  )

  fetch_quote() {
    local api_url=$1
    
    response=$(curl -s -w "\n%{http_code}" \
      -H "Accept: application/json" \
      -H "User-Agent: Mozilla/5.0 (Bash Quote Fetcher)" \
      --max-time 5 \
      "$api_url")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [[ $DEBUG == 1 ]]; then
      echo "Trying: $api_url"
      echo "HTTP Code: $http_code"
    fi
    
    # Check for success
    if [[ "$http_code" == "200" ]] && echo "$body" | jq -e . > /dev/null 2>&1; then
      # Handle different API response formats
      case "$api_url" in
        *quotable.io*)
          QUOTE=$(echo "$body" | jq -r '.content + " - " + .author')
          ;;
        *zenquotes.io*)
          QUOTE=$(echo "$body" | jq -r '.[0].q + " - " + .[0].a')
          ;;
        *quoteslate*)
          QUOTE=$(echo "$body" | jq -r '.quote + " - " + .author')
          ;;
      esac
      
      if [[ -n "$QUOTE" && "$QUOTE" != " - " ]]; then
        return 0
      fi
    fi
    
    return 1
  }

  # Try each API until one works
  for api in "${APIS[@]}"; do
    if fetch_quote "$api"; then
      newBannerColor "$QUOTE" "green" "*"
      exit 0
    fi
  done

  # All APIs failed
  newBannerColor "⚠️ All quote APIs are unavailable (rate limited or down)" "yellow" "*"
  exit 0
}
