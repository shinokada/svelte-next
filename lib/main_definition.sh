# shellcheck disable=SC1083
parser_definition() {
    setup REST help:usage abbr:true -- \
        "Usage: ${2##*/} [command] [options...] [arguments...]"
    msg -- '' 'Options:'
    disp :usage -h --help
    disp VERSION --version

    msg -- '' 'Commands: '
    msg -- 'Use command -h for a command help.'
    cmd svelte -- "Install svelte version"
    # cmd cmd2 -- "CMD2 description."
    # cmd text_example -- "Print different type of texts."
    # cmd create -- "Create this and that."

    msg -- '' "Examples:
    
    i 
    $SCRIPT_NAME i 118 .
    version 
    $SCRIPT_NAME --version
    Display help:
    $SCRIPT_NAME -h | --help
"
}
