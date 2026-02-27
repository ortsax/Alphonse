#!/usr/bin/env bash
# Orstax installer for Linux.
# Must be run as root (or with sudo).
set -euo pipefail

# ── Configuration ────────────────────────────────────────────────────────────
REPO_URL="https://github.com/ortsax/whatsapp-bot.git"
SRC_DIR="/opt/orstax/src"
BIN_PATH="/usr/local/bin/orstax"
GO_FALLBACK="1.25.0"
GOROOT="/usr/local/go"
# ─────────────────────────────────────────────────────────────────────────────

step()  { echo; echo "==> $*"; }
ok()    { echo "    $*"; }
err()   { echo "    ERROR: $*" >&2; exit 1; }

# ── Require root ─────────────────────────────────────────────────────────────
if [ "$(id -u)" -ne 0 ]; then
    err "This script must be run as root. Try: sudo bash $0"
fi

# ── Detect architecture ───────────────────────────────────────────────────────
step "Detecting system architecture"
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    armv6l)  GOARCH="armv6l" ;;
    i686)    GOARCH="386" ;;
    *) err "Unsupported architecture: $ARCH" ;;
esac
ok "Architecture: $GOARCH"

# ── Detect latest Go version ─────────────────────────────────────────────────
step "Detecting latest Go version"
if command -v curl &>/dev/null; then
    GO_VERSION=$(curl -fsSL "https://go.dev/dl/?mode=json" 2>/dev/null \
        | grep -o '"version":"go[^"]*"' | head -1 \
        | grep -o '[0-9][0-9.]*' | head -1) || true
elif command -v wget &>/dev/null; then
    GO_VERSION=$(wget -qO- "https://go.dev/dl/?mode=json" 2>/dev/null \
        | grep -o '"version":"go[^"]*"' | head -1 \
        | grep -o '[0-9][0-9.]*' | head -1) || true
fi
GO_VERSION="${GO_VERSION:-$GO_FALLBACK}"
ok "Using Go $GO_VERSION"

# ── Check / install dependencies ─────────────────────────────────────────────
step "Checking dependencies"

# git
if ! command -v git &>/dev/null; then
    ok "Installing git"
    if command -v apt-get &>/dev/null; then
        apt-get update -qq && apt-get install -y -qq git
    elif command -v yum &>/dev/null; then
        yum install -y -q git
    elif command -v dnf &>/dev/null; then
        dnf install -y -q git
    elif command -v pacman &>/dev/null; then
        pacman -Sy --noconfirm git
    else
        err "Cannot install git automatically. Please install it manually."
    fi
fi
ok "Git: $(git --version)"

# curl/wget (for Go download)
if ! command -v curl &>/dev/null && ! command -v wget &>/dev/null; then
    if command -v apt-get &>/dev/null; then
        apt-get install -y -qq curl
    elif command -v yum &>/dev/null; then
        yum install -y -q curl
    fi
fi

# ── Check / install Go ───────────────────────────────────────────────────────
step "Checking Go"
NEED_GO=true
if command -v go &>/dev/null; then
    INSTALLED=$(go version | grep -o '[0-9][0-9.]*' | head -1)
    ok "Go $INSTALLED already installed"
    NEED_GO=false
fi

if [ "$NEED_GO" = true ]; then
    step "Installing Go $GO_VERSION"
    TARBALL="/tmp/go${GO_VERSION}.linux-${GOARCH}.tar.gz"
    DL_URL="https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz"
    ok "Downloading $DL_URL"
    if command -v curl &>/dev/null; then
        curl -fsSL "$DL_URL" -o "$TARBALL"
    else
        wget -q "$DL_URL" -O "$TARBALL"
    fi
    rm -rf "$GOROOT"
    tar -C /usr/local -xzf "$TARBALL"
    rm "$TARBALL"

    # Add to PATH system-wide
    echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
    chmod +x /etc/profile.d/go.sh
    export PATH="$PATH:/usr/local/go/bin"
    ok "Go $GO_VERSION installed"
fi

export PATH="$PATH:$GOROOT/bin"

# ── Clone or update repo ─────────────────────────────────────────────────────
step "Setting up source"
mkdir -p "$(dirname "$SRC_DIR")"
if [ -d "$SRC_DIR/.git" ]; then
    ok "Updating existing clone"
    git -C "$SRC_DIR" pull
else
    ok "Cloning $REPO_URL"
    git clone "$REPO_URL" "$SRC_DIR"
fi

# ── Build ────────────────────────────────────────────────────────────────────
step "Building orstax"
cd "$SRC_DIR"
CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.sourceDir=${SRC_DIR}" \
    -trimpath \
    -o "$BIN_PATH" \
    .
chmod +x "$BIN_PATH"
ok "Binary written to $BIN_PATH"

# ── Done ─────────────────────────────────────────────────────────────────────
echo
echo "Orstax is now installed"
echo
echo "  Run with    orstax --phone-number <international-number>"
echo "  Update with orstax -update"
echo "  Sessions    orstax -list-sessions"
echo "              orstax -delete-session <phone>"
echo "              orstax -reset-session  <phone>"
echo
echo "Note: open a new shell or run 'source /etc/profile.d/go.sh' if go was just installed."
