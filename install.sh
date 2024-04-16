#!/usr/bin/env bash
set -euo pipefail


# Reset
Color_Off=''

# Regular Colors
Red=''
Green=''
Dim='' # White

# Bold
Bold_White=''
Bold_Green=''

if [[ -t 1 ]]; then
    # Reset
    Color_Off='\033[0m' # Text Reset

    # Regular Colors
    Red='\033[0;31m'   # Red
    Green='\033[0;32m' # Green
    Dim='\033[0;2m'    # White

    # Bold
    Bold_Green='\033[1;32m' # Bold Green
    Bold_White='\033[1m'    # Bold White
fi

error() {
    echo -e "${Red}error${Color_Off}:" "$@" >&2
    exit 1
}

info() {
    echo -e "${Dim}$@ ${Color_Off}"
}

info_bold() {
    echo -e "${Bold_White}$@ ${Color_Off}"
}

success() {
    echo -e "${Green}$@ ${Color_Off}"
}

command -v unzip >/dev/null ||
    error 'unzip is required to install edge'


case $(uname -ms) in
'Linux aarch64' | 'Linux arm64')
    target=linux-amd64
    ;;
'Linux x86_64' | *)
    target=linux-386
    ;;
esac

GITHUB=${GITHUB-"https://github.com"}

github_repo="$GITHUB/liteseed/edge"

exe_name=edge

edge_uri=$github_repo/releases/latest/download/edge-$target.zip

install_env=EDGE_INSTALL
bin_env=\$$install_env/edge

install_dir=${!install_env:-$HOME/.edge}
exe=$install_dir/edge

if [[ ! -d $install_dir ]]; then
    mkdir -p "$install_dir" ||
        error "Failed to create install directory \"$install_dir\""
fi

curl --fail --location --progress-bar --output "$exe.zip" "$edge_uri" ||
    error "Failed to download edge from \"$edge_uri\""

unzip -oqd "$install_dir" "$exe.zip" ||
    error 'Failed to extract edge'

chmod +x "$exe" ||
    error 'Failed to set permissions on edge executable'

rm -r "$install_dir/edge-$target" "$exe.zip"