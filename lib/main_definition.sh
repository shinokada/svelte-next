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
    
    - Use $SCRIPT_NAME update target-directory. For example, installing  svelte@5.x.x, running pnmp update, git add, commit, and push, pnpm test:integration in the subdirectories of the current directory:

    $SCRIPT_NAME update .

    - Use -n or --next number to install a certain Svelte next version.
    $SCRIPT_NAME update -n 5.0.0 .

    - Use -L or --latest to update all packages to their latest versions (ignores semver ranges):
    $SCRIPT_NAME update -L .

    - Use -p flag to NOT to run pnpm update:
    $SCRIPT_NAME update -p .

    - Use -g flag to NOT to run git add, commit, and push:
    $SCRIPT_NAME update -g .

    - Use -t flag to NOT to run pnpm test:integration:
    $SCRIPT_NAME update -t .

    - Use -f or --from number to start from a certain index:
    $SCRIPT_NAME update -f 2 .

    - Use -d flag to run in debug mode:
    $SCRIPT_NAME update -d .

    - Combine the flags:
    $SCRIPT_NAME update -dtgp -n 5.0.0 .

    - To display version: 
    $SCRIPT_NAME --version
    
    - To display help:
    $SCRIPT_NAME -h | --help
"
}
