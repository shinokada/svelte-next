check_cmd() {
    if [ ! "$(command -v "$1")" ]; then
        app=$1
        redprint "It seems like you don't have ${app}."
        redprint "Please install ${app}."
        exit 1
    fi
}

# bash version check
check_bash() {
    # If you are using Bash, check Bash version
    if ((BASH_VERSINFO[0] < $1)); then
        printf '%s\n' "Error: This requires Bash v${1} or higher. You have version $BASH_VERSION." 1>&2
        exit 2
    fi
}

### Colors ##
ESC=$(printf '\033')
RESET="${ESC}[0m"
BLACK="${ESC}[30m"
RED="${ESC}[31m"
GREEN="${ESC}[32m"
YELLOW="${ESC}[33m"
BLUE="${ESC}[34m"
MAGENTA="${ESC}[35m"
CYAN="${ESC}[36m"
WHITE="${ESC}[37m"
DEFAULT="${ESC}[39m"

### Color Functions ##

blackprint() {
    printf "${BLACK}%s${RESET}\n" "$1"
}

blueprint() {
    printf "${BLUE}%s${RESET}\n" "$1"
}

cyanprint() {
    printf "${CYAN}%s${RESET}\n" "$1"
}

defaultprint() {
    printf "${DEFAULT}%s${RESET}\n" "$1"
}

greenprint() {
    printf "${GREEN}%s${RESET}\n" "$1"
}

magentaprint() {
    printf "${MAGENTA}%s${RESET}\n" "$1"
}

redprint() {
    printf "${RED}%s${RESET}\n" "$1"
}

whiteprint() {
    printf "${WHITE}%s${RESET}\n" "$1"
}

yellowprint() {
    printf "${YELLOW}%s${RESET}\n" "$1"
}

# Text attribute
BOLD="${ESC}[1m"
FAINT="${ESC}[2m"
ITALIC="${ESC}[3m"
UNDERLINE="${ESC}[4m"
SLOWBLINK="${ESC}[5m"
SWAP="${ESC}[7m"
STRIKE="${ESC}[9m"

boldprint() {
    printf "${BOLD}%s${RESET}\n" "$1"
}

faintprint() {
    printf "${FAINT}%s${RESET}\n" "$1"
}

italicprint() {
    printf "${ITALIC}%s${RESET}\n" "$1"
}

underlineprint() {
    printf "${UNDERLINE}%s${RESET}\n" "$1"
}

slowblinkprint() {
    printf "${SLOWBLINK}%s${RESET}\n" "$1"
}

swapprint() {
    printf "${SWAP}%s${RESET}\n" "$1"
}

strikeprint() {
    printf "${STRIKE}%s${RESET}\n" "$1"
}


# lib/banners
# Usage: bannerSimple "my title" "*"
function bannerSimple() {
    msg="${2} ${1} ${2}"
    edge=$(echo "${msg}" | sed "s/./"${2}"/g")
    echo "${edge}"
    echo "$(tput bold)${msg}$(tput sgr0)"
    echo "${edge}"
    echo
}

# Usage: bannerColor "my title" "red" "*"
function bannerColor() {
    case ${2} in
    black)
        color=0
        ;;
    red)
        color=1
        ;;
    green)
        color=2
        ;;
    yellow)
        color=3
        ;;
    blue)
        color=4
        ;;
    magenta)
        color=5
        ;;
    cyan)
        color=6
        ;;
    white)
        color=7
        ;;
    *)
        echo "color is not set"
        exit 1
        ;;
    esac

    msg="${3} ${1} ${3}"
    edge=$(echo "${msg}" | sed "s/./${3}/g")
    tput setaf ${color}
    tput bold
    echo "${edge}"
    echo "${msg}"
    echo "${edge}"
    tput sgr 0
    echo
}


# Usage: bannerColor "my title" "blue" "*" [border_width]
function newBannerColor() {
    case ${2} in
    black)
        color=0
        ;;
    red)
        color=1
        ;;
    green)
        color=2
        ;;
    yellow)
        color=3
        ;;
    blue)
        color=4
        ;;
    magenta)
        color=5
        ;;
    cyan)
        color=6
        ;;
    white)
        color=7
        ;;
    *)
        echo "color is not set"
        exit 1
        ;;
    esac

    # Set border width to 4th argument if provided, otherwise default to 10
    border_width=${4:-10}
    
    # Create border string with specified width
    border=$(printf "%0.s${3}" $(seq 1 $border_width))
    
    tput setaf ${color}
    tput bold
    echo "${border}"
    echo "${1}"  # Print the message as-is, without adding border characters
    echo "${border}"
    tput sgr 0
    echo
}

# Array of emoji Unicode codes
emojis=(
    "\U1F921"  # 🤡 clown face
    "\U1F479"  # 👹 ogre
    "\U1F47A"  # 👺 goblin
    "\U1F47B"  # 👻 ghost
    "\U1F47D"  # 👽 alien
    "\U1F47E"  # 👾 alien monster
    "\U1F916"  # 🤖 robot
    "\U1F348"  # 🍈 melon
    "\U1F349"  # 🍉 watermelon
    "\U1F34A"  # 🍊 tangerine
    "\U1F34B"  # 🍋 lemon
    "\U1F34C"  # 🍌 banana
    "\U1F34D"  # 🍍 pineapple
    "\U1F96D"  # 🥭 mango
    "\U1F34E"  # 🍎 red apple
    "\U1F34F"  # 🍏 green apple
    "\U1F350"  # 🍐 pear
    "\U1F351"  # 🍑 peach
    "\U1F352"  # 🍒 cherries
    "\U1F353"  # 🍓 strawberry
    "\U1F433"  # 🐳 whale
    "\U1F419"  # 🐙 octopus
    "\U1F98B"  # 🦋 butterfly
    "\U1F439"  # 🐹 hamster
    "\U1F431"  # 🐱 cat
    "\U1F436"  # 🐶 dog
    "\U1F680"  # 🚀 rocket
    "\U1F525"  # 🔥 fire
    "\U1F355"  # 🍕 pizza
    "\U1F354"  # 🍔 hamburger
    "\U1F338"  # 🌸 cherry
    "\U1F339"  # 🌹 rose
    "\U1F33A"  # 🌺 hibiscus
    "\U1F33B"  # 🌻 sunflower
    "\U1F33C"  # 🌼 blossom
    "\U1F337"  # 🌷 tulip
    "\U1F331"  # 🌱 seedling
    "\U1F332"  # 🌲 evergreen tree
    "\U1F333"  # 🌳 deciduous tree
)

# Get a random emoji from the array
random_emoji=${emojis[$RANDOM % ${#emojis[@]}]}