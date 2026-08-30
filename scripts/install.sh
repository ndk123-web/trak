#!/usr/bin/env bash
set -e

# ==============================
# TRAK Linux & macOS Installer
# ==============================

VERSION="v1.0.0"
INSTALL_DIR="${HOME}/.trak/bin"
EXE_PATH="${INSTALL_DIR}/trak"

echo ""
echo -e "\033[36mTRAK Installer\033[0m"
echo "=============="
echo ""

# 1. Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux)
    TARGET="trak-linux-${ARCH}"
    ;;
  darwin)
    TARGET="trak-darwin-${ARCH}"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

DOWNLOAD_URL="https://github.com/ndk123-web/trak/releases/download/${VERSION}/${TARGET}"

echo "[1/4] Detected ${OS} (${ARCH})"

# 2. Create install directory
mkdir -p "${INSTALL_DIR}"
echo "[2/4] Install directory: ${INSTALL_DIR}"

# 3. Download binary
echo ""
echo "[3/4] Downloading TRAK ${VERSION}..."
curl -fsSL "${DOWNLOAD_URL}" -o "${EXE_PATH}"
chmod +x "${EXE_PATH}"
echo "      Downloaded successfully."

# 4. PATH instructions
echo ""
echo "[4/4] Configuring PATH..."

SHELL_RC=""
if [ -n "$ZSH_VERSION" ] || [ "$SHELL" = "/bin/zsh" ]; then
  SHELL_RC="${HOME}/.zshrc"
elif [ -n "$BASH_VERSION" ] || [ "$SHELL" = "/bin/bash" ]; then
  SHELL_RC="${HOME}/.bashrc"
fi

if [ -n "$SHELL_RC" ] && [ -f "$SHELL_RC" ]; then
  if ! grep -q 'export PATH="$HOME/.trak/bin:$PATH"' "$SHELL_RC"; then
    echo 'export PATH="$HOME/.trak/bin:$PATH"' >> "$SHELL_RC"
    echo "      Added to ${SHELL_RC}"
  else
    echo "      Already present in ${SHELL_RC}"
  fi
fi

echo ""
echo -e "\033[32mTRAK installed successfully!\033[0m"
echo ""
echo "Installed at: ${EXE_PATH}"
echo ""
echo "Run:"
echo "  export PATH=\"\$HOME/.trak/bin:\$PATH\""
echo "  trak list"
echo ""
