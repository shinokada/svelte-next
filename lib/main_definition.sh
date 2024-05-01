# shellcheck disable=SC1083
parser_definition() {
    setup REST help:usage abbr:true -- \
        "Usage: ${2##*/} [command] [options...]"
    msg -- '' 'Options:'
    disp :usage -h --help
    disp VERSION --version

    msg -- '' 'Commands: '
    msg -- 'Use command -h for a command help.'
    cmd update -- "Install next version in the subdirectories of the current directory"

    msg -- '' "Examples:
    
    - Use $SCRIPT_NAME update target-directory. For example, installing latest svelte@5.0.0-next, running pnmp update, git add, commit, and push, pnpm test:integration in the subdirectories of the current directory:

    $SCRIPT_NAME .

    - Use -v param to install a certain Svelte next version.

    $SCRIPT_NAME update -v 120 .

    - Use -p flag to NOT to run pnpm update:

    $SCRIPT_NAME update -p .

    - Use -g flag to NOT to run git add, commit, and push:
    $SCRIPT_NAME update -g .

    - Use -t flag to NOT to run pnpm test:integration:
    $SCRIPT_NAME update -t .

    - To display version: 
    $SCRIPT_NAME --version
    
    - To display help:
    $SCRIPT_NAME -h | --help
"
}
