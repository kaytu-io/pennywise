#!/bin/bash

sudo -v

# Detect architecture
arch=$(uname -m | tr '[:upper:]' '[:lower:]' | sed -e s/x86_64/amd64/)
if [ "$arch" = "aarch64" ]; then
  arch="arm64"
fi

# Detect operating system
OS=$(uname -s)

# Determine the release based on OS and architecture
case "$OS" in
  Linux)
    DOWNLOAD_URL=$(curl -s -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kaytu-io/pennywise/releases/latest | jq -r --arg ARCH "$arch" '.assets[] | select(.browser_download_url | contains("linux") and contains($ARCH)) | .browser_download_url')
    ;;
  Darwin)
    DOWNLOAD_URL=$(curl -s -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kaytu-io/pennywise/releases/latest | jq -r --arg ARCH "$arch" '.assets[] | select(.browser_download_url | contains("mac") and contains($ARCH)) | .browser_download_url')
    ;;
  *)
    echo "Unsupported operating system: $OS"
    exit 1
    ;;
esac

# Download and install
curl -L "$DOWNLOAD_URL" -o ./pennywise
chmod +x ./pennywise
sudo mv ./pennywise /usr/local/bin/

echo "Pennywise installed successfully version $(pennywise version)"