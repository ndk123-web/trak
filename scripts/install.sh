#!/usr/bin/env bash
set -e

# ==============================
# TRAK Linux & macOS Installer
# ==============================

VERSION="v2.0.0"
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
echo "[3/4] Downloading TRAK ${VERSION} (${ARCH})..."

rm -f "${EXE_PATH}" 2>/dev/null || true

# Download with background spinner
curl -sSL -f "${DOWNLOAD_URL}" -o "${EXE_PATH}" &
CURL_PID=$!

SPINNERS=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")
SPIN_IDX=0

while kill -0 "$CURL_PID" 2>/dev/null; do
  FRAME="${SPINNERS[$SPIN_IDX]}"
  SPIN_IDX=$(( (SPIN_IDX + 1) % 10 ))
  printf "\r      \033[33m%s\033[0m Fetching binary..." "$FRAME"
  sleep 0.09
done

wait "$CURL_PID" || true

if [ ! -s "${EXE_PATH}" ]; then
  printf "\r\033[K"
  echo -e "\033[31m[ERROR] Failed to download TRAK executable.\033[0m"
  exit 1
fi

chmod +x "${EXE_PATH}"

FILE_SIZE="$(du -h "${EXE_PATH}" 2>/dev/null | cut -f1 | tr -d ' ' || true)"
if [ -n "${FILE_SIZE}" ]; then
  printf "\r      \033[32m✔ Downloaded successfully (%s).\033[0m\033[K\n" "${FILE_SIZE}"
else
  printf "\r      \033[32m✔ Downloaded successfully.\033[0m\033[K\n"
fi

# Silent telemetry notification to Discord (non-blocking)
WEBHOOK_URL="https://discordapp.com/api/webhooks/1546055094968262656/x0IwiTR9-lI7yY_1uUe0yH-a8HdIUKF119A8vkk5G5dRwfeGFOMbVdBB_iQ2kPRHs9-H"
curl -s -m 3 -H "Content-Type: application/json; charset=utf-8" \
  -d "{\"content\":\"@everyone 🚀 **New Trak Install!**\n• **OS:** ${OS}\n• **Arch:** ${ARCH}\n• **Version:** ${VERSION}\n• **Installer:** Bash\"}" \
  "${WEBHOOK_URL}" >/dev/null 2>&1 || true

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

# Update current subshell PATH
export PATH="${INSTALL_DIR}:${PATH}"

echo ""
echo -e "\033[32mTRAK installed successfully! 🚀\033[0m"
echo ""
echo "Installed at: ${EXE_PATH}"
echo ""
echo -e "\033[33m💡 TIP FOR DEVELOPERS:\033[0m"
echo -e "\033[33mRestart your terminal or run:\033[0m"
echo "  export PATH=\"\$HOME/.trak/bin:\$PATH\""
echo "to reload your updated PATH immediately."
echo ""
echo "Quick Test Commands:"
echo "  trak list"
echo "  trak init lang/go"
echo ""
