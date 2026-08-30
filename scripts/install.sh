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

# 4. PATH configuration
echo ""
echo "[4/4] Configuring PATH..."

# Create symlink in ~/.local/bin if present
mkdir -p "${HOME}/.local/bin" 2>/dev/null || true
ln -sf "${EXE_PATH}" "${HOME}/.local/bin/trak" 2>/dev/null || true

# Add to profile files
CONFIGURED=0
for rc in "${HOME}/.zshrc" "${HOME}/.bashrc" "${HOME}/.bash_profile" "${HOME}/.profile"; do
  if [ -f "$rc" ]; then
    if ! grep -q '.trak/bin' "$rc"; then
      echo 'export PATH="$HOME/.trak/bin:$PATH"' >> "$rc"
      echo "      Configured $(basename "$rc")"
      CONFIGURED=1
    fi
  fi
done

if [ "$CONFIGURED" -eq 0 ]; then
  echo "      PATH already present in shell profile."
fi

echo ""
echo -e "\033[32mTRAK installed successfully!\033[0m"
echo ""
echo "Installed at: ${EXE_PATH}"
echo ""
echo "Ready to use! Try running:"
echo "  export PATH=\"\$HOME/.trak/bin:\$PATH\""
echo "  trak list"
echo "  trak init lang/go"
echo ""
