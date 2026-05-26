#!/usr/bin/env bash
set -euo pipefail

GITANT_REPO="${GITANT_REPO:-GrayCodeAI/gitant-cli}"
GITANT_VERSION="${GITANT_VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

if [ "$OS" = "darwin" ]; then
  OS="Darwin"
elif [ "$OS" = "linux" ]; then
  OS="Linux"
else
  echo "Unsupported OS: $OS" >&2
  exit 1
fi

if [ "$GITANT_VERSION" = "latest" ]; then
  GITANT_VERSION="$(curl -fsSL "https://api.github.com/repos/${GITANT_REPO}/releases/latest" | grep '"tag_name"' | head -1 | cut -d'"' -f4)"
fi

ARCHIVE="gitant-cli_${GITANT_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${GITANT_REPO}/releases/download/${GITANT_VERSION}/${ARCHIVE}"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Downloading ${URL}..."
curl -fsSL "$URL" -o "${TMP}/${ARCHIVE}"
tar xzf "${TMP}/${ARCHIVE}" -C "$TMP"

install -m 755 "${TMP}/gitant" "${INSTALL_DIR}/gitant"
install -m 755 "${TMP}/git-remote-gitant" "${INSTALL_DIR}/git-remote-gitant"

echo "Installed gitant and git-remote-gitant to ${INSTALL_DIR}"
gitant version
