# shellcheck disable=SC1083
parser_definition() {
    setup REST help:usage abbr:true -- \
        "Usage: ${2##*/} [command] [options...] [arguments...]"
    msg -- '' 'Options:'
    disp :usage -h --help
    disp VERSION --version

    msg -- '' 'Commands: '
    msg -- 'Use command -h for a command help.'
    cmd svelte -- "Install next version in the subdirectories of the current directory"
    # cmd cmd2 -- "CMD2 description."
    # cmd text_example -- "Print different type of texts."
    # cmd create -- "Create this and that."

    msg -- '' "Examples:
    
    Install svelte@5.0.0-next.120 in the subdirectories of the current directory:
    $SCRIPT_NAME update 120 .
    version: 
    $SCRIPT_NAME --version
    Display help:
    $SCRIPT_NAME -h | --help
"
}
