parser_definition_svelte() {
    # from https://github.com/ko1nksm/getoptions/blob/master/examples/basic.sh
    setup REST plus:true help:usage abbr:true -- \
        "Usage: ${2##*/} [options...] [arguments...]" ''
    msg -- 'getoptions basic example' ''
    msg -- 'Options:'
    flag FLAG_P -p +p --{no-}flag-p on:0 no:1 init:@no -- "Use -p not to run pnpm update"
    flag FLAG_T -t +t --{no-}flag-t on:0 no:1 init:@no -- "Use -t not to run pnpm test:integration"
    flag FLAG_G -g +g --{no-}flag-t on:0 no:1 init:@no -- "Use -g not to run git add, commit, and push"
    disp :usage -h --help
}
