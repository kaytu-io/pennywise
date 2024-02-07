#!/usr/bin/env bash
set -e

DEFAULT_BIN_DIR="/usr/local/bin"
BIN_DIR=${1:-"${DEFAULT_BIN_DIR}"}
GITHUB_REPO="kaytu-io/pennywise"
TMP_DIR=$(mktemp -d -t pennywise-install.XXXXXXXXXX)

# Helper functions for logs
info() {
    echo '[INFO] ' "$@"
}

warn() {
    echo '[WARN] ' "$@" >&2
}

fatal() {
    echo '[ERROR] ' "$@" >&2
    exit 1
}

# Set os, fatal if operating system not supported
setup_verify_os() {
    OS=$(uname -s)
    case ${OS} in
        Darwin)
            OS=macos
            ;;
        Linux)
            OS=linux
            ;;
        *)
            fatal "Unsupported operating system ${OS}"
    esac
}

# Set arch, fatal if architecture not supported
setup_verify_arch() {
    ARCH=$(uname -m)

    case ${ARCH} in
        arm64|aarch64|armv8l)
            ARCH=arm64
            ;;
        amd64)
            ARCH=amd64
            ;;
        x86_64)
            ARCH=amd64
            ;;
        *)
            fatal "Unsupported architecture ${ARCH}"
    esac
}

# Verify existence of downloader executable
verify_downloader() {
    # Return failure if it doesn't exist or is no executable
    [ -x "$(which "$1")" ] || return 1

    # Set verified executable as our downloader program and return success
    DOWNLOADER=$1
    return 0
}


# Find version from Github metadata
get_release_version() {
    if [ -n "${FLUX_VERSION}" ]; then
      SUFFIX_URL="tags/v${FLUX_VERSION}"
    else
      SUFFIX_URL="latest"
    fi

    METADATA_URL="https://api.github.com/repos/${GITHUB_REPO}/releases/${SUFFIX_URL}"
    TMP_METADATA="${TMP_DIR}/pennywise.json"
    info "Downloading metadata ${METADATA_URL}"
    download "${TMP_METADATA}" "${METADATA_URL}"
    VERSION=$(grep '"tag_name":' "${TMP_METADATA}" | sed -E 's/.*"([^"]+)".*/\1/' | cut -c 2-)
    if [ -n "${VERSION}" ]; then
        info "Using ${VERSION} as release"
    else
        fatal "Unable to determine release version"
    fi
}

# Download from file from URL
download() {
    [ $# -eq 2 ] || fatal 'download needs exactly 2 arguments'

    case $DOWNLOADER in
        curl)
            curl -o "$1" -sfL "$2"
            ;;
        wget)
            wget --auth-no-challenge -qO "$1" "$2"
            ;;
        *)
            fatal "Incorrect executable '${DOWNLOADER}'"
            ;;
    esac

    # Abort if download command failed
    [ $? -eq 0 ] || fatal 'Download failed'
}


# Download binary from Github URL
download_binary() {
    BIN_URL="https://kaytu.s3.amazonaws.com/pennywise/releases/tag/v1.6.16/pennywise-${OS}-${ARCH}"
    info "Downloading binary ${BIN_URL}"
    TMP_BIN="${TMP_DIR}/pennywise"

    download "${TMP_BIN}" "${BIN_URL}"
}


# Setup permissions and move binary
setup_binary() {
    chmod +x "${TMP_BIN}"
    info "Installing pennywise to ${BIN_DIR}/pennywise"

    local CMD_MOVE="mv -f \"${TMP_DIR}/pennywise\" \"${BIN_DIR}\""
    if [ -w "${BIN_DIR}" ]; then
        eval "${CMD_MOVE}"
    else
        eval "sudo ${CMD_MOVE}"
    fi
}

# Run the install process
{
    setup_verify_os
    setup_verify_arch
    verify_downloader curl || verify_downloader wget || fatal 'Can not find curl or wget for downloading files'
#    get_release_version
    download_binary
    setup_binary
}