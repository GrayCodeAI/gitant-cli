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

case "$OS" in
  darwin|linux) ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
esac

if [ "$GITANT_VERSION" = "latest" ]; then
  GITANT_VERSION="$(curl -fsSL "https://api.github.com/repos/${GITANT_REPO}/releases/latest" | grep '"tag_name"' | head -1 | cut -d'"' -f4)"
fi

# Goreleaser archives: gitant-cli_0.1.0_darwin_arm64.tar.gz (no "v", lowercase OS)
VERSION="${GITANT_VERSION#v}"
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="zip"
fi

ARCHIVE="gitant-cli_${VERSION}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/${GITANT_REPO}/releases/download/${GITANT_VERSION}/${ARCHIVE}"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Downloading ${URL}..."
curl -fsSL "$URL" -o "${TMP}/${ARCHIVE}"

if [ "$EXT" = "zip" ]; then
  unzip -q "${TMP}/${ARCHIVE}" -d "$TMP"
else
  tar xzf "${TMP}/${ARCHIVE}" -C "$TMP"
fi

mkdir -p "${INSTALL_DIR}"
install -m 755 "${TMP}/gitant" "${INSTALL_DIR}/gitant"
install -m 755 "${TMP}/git-remote-gitant" "${INSTALL_DIR}/git-remote-gitant"

echo "Installed gitant and git-remote-gitant to ${INSTALL_DIR}"
"${INSTALL_DIR}/gitant" version
