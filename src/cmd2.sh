fn_cmd2() {
  echo "hello"
  bannerColor "Whew! Finally done. I'm outta here." "blue" "*" 
  QUOTE=$(curl -s https://api.quotable.io/quotes/random | jq -r '.[0].content + " - " + .[0].author')

  if [[ -n "$QUOTE" ]]; then
    bannerColor "$QUOTE" "blue" "*" 
  fi
}