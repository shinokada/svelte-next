parser_definition_svelte() {
    # from https://github.com/ko1nksm/getoptions/blob/master/examples/basic.sh
    setup REST plus:true help:usage abbr:true -- \
        "Usage: ${2##*/} [options...] [arguments...]" ''
    msg -- 'getoptions basic example' ''
    msg -- 'Options:'
    flag FLAG_P -p +p --{no-}flag-p on:0 no:1 init:@no -- "Use -p not to run pnpm update"
    flag FLAG_S -s +s --{no-}flag-s on:0 no:1 init:@no -- "Use -s not to run pnpm i -D svelte@next"
    flag FLAG_T -t +t --{no-}flag-t on:0 no:1 init:@no -- "Use -t not to run pnpm test:integration"
    flag FLAG_G -g +g --{no-}flag-t on:0 no:1 init:@no -- "Use -g not to run git add, commit, and push"
    param FROM -f --from validate:number -- "Use -f for starting index of subdirectory"
    param SVELTE_NEXT -n --next validate:number -- "Svelte 5 version number"
    disp :usage -h --help
}

error() {
	case $2 in
		unknown) echo "$1" ;;
		number:*) echo "Not a number: $3" ;;
		*) return 0 ;; # Display default error
	esac
	return 1
}

number() {
	case $OPTARG in (*[!0-9]*) return 1; esac
}
