#!/usr/bin/env bash
set -e

owner="{{ .Owner }}"
repo="{{ .Repo }}"
version="{{ .Version }}"
binary="{{ .Binary }}"
bin_dir="{{ .BinDir }}" # The executable file will finally be placed here

download_folder="${HOME}/Downloads" # the file should be download to here

if [ -z "$bin_dir" ]; then
    bin_dir="/usr/local/bin"
fi

if [ ! -d "$bin_dir" ]; then
    mkdir -p "$bin_dir"
fi

if [ ! -d "$download_folder" ]; then
    mkdir -p "$download_folder"
fi

get_arch() {
    # https://man7.org/linux/man-pages/man1/uname.1.html
    # https://en.wikipedia.org/wiki/Uname
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" | "i686-64" )
            echo "amd64"
        ;;
        "i386" | "i686" | "i486" | "i586" | "i86pc" | "x86pc" )
            echo "386"
        ;;
        "arm32")
            echo "arm32"
        ;;
        "arm64" | "aarch64" )
            echo "arm64"
        ;;
        "armv8")
            echo "armv8"
        ;;
        "armv7" | "armv7l" )
            echo "armv7"
        ;;
        "armv6" | "armv6l" | "arm" )
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

os=$(get_os)
arch=$(get_arch)
dest_file="${download_folder}/${binary}_${os}_${arch}.tar.gz"
asset_uri="https://github.com/${owner}/${repo}/releases/download/${version}/${binary}_${os}_${arch}.tar.gz"

function clean {
    rm  -r ${dest_file}
}

main() {
    mkdir -p ${download_folder}
    
    echo "[1/3] Download ${asset_uri} to ${download_folder}"
    rm -f ${dest_file}
    curl --location --output "${dest_file}" "${asset_uri}"

    trap clean EXIT
    
    echo "[2/3] Install '${binary}' to the ${bin_dir}"
    mkdir -p ${HOME}/bin
    tar -xz -f ${dest_file} -C ${bin_dir}
    exe=${bin_dir}/${binary}
    chmod +x ${exe}
    
    echo "[3/3] Set environment variables"
    echo "${binary} was installed successfully to ${exe}"
    if command -v ${binary} --version >/dev/null; then
        echo "Run '${binary} --help' to get started"
    else
        echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
        echo "  export PATH=${bin_dir}:\$PATH"
        echo "Run '${binary} --help' to get started"
    fi
    
    exit 0
}

main