#!/usr/bin/env bash
set -e

owner="{{ .Owner }}"
repo="{{ .Repo }}"
version="{{ .Version }}"
binary="{{ .Binary }}"

# the file should be download to here
downloadFolder="${HOME}/Downloads"
# The executable file will finally be placed here
binDir=/usr/local/bin

mkdir -p ${downloadFolder}

get_arch() {
    # https://man7.org/linux/man-pages/man1/uname.1.html
    # https://en.wikipedia.org/wiki/Uname
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" | "i686-64" )
            echo "amd64"
        ;;
        "i386" | "i686" | "i486" | "i586" | "i86pc" ｜ "x86pc")
            echo "386"
        ;;
        "arm32")
            echo "arm32"
        ;;
        "arm64" | "aarch64")
            echo "arm64"
        ;;
        "armv8")
            echo "armv8"
        ;;
        "armv7" | "armv7l")
            echo "armv7"
        ;;
        "armv6" | "armv6l" | "arm")
            echo "armv6"
        ;;
        "armv5")
            echo "armv5"
        ;;
        *)
            echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

main() {
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${downloadFolder}/${binary}_${os}_${arch}.tar.gz"

    asset_uri="https://github.com/${owner}/${repo}/releases/download/${version}/${binary}_${os}_${arch}.tar.gz"

    mkdir -p ${downloadFolder}

    echo "[1/3] Download ${asset_uri} to ${downloadFolder}"
    rm -f ${dest_file}
    curl --location --output "${dest_file}" "${asset_uri}"

    echo "[2/3] Install '${binary}' to the ${binDir}"
    mkdir -p ${HOME}/bin
    tar -xz -f ${dest_file} -C ${binDir}
    exe=${binDir}/${binary}
    chmod +x ${exe}

    echo "[3/3] Set environment variables"
    echo "${binary} was installed successfully to ${exe}"
    if command -v ${binary} --version >/dev/null; then
        echo "Run '${binary} --help' to get started"
    else
        echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
        echo "  export PATH=${HOME}/bin:\$PATH"
        echo "Run '${binary} --help' to get started"
    fi

    exit 0
}

main