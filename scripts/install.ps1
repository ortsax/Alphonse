#Requires -RunAsAdministrator
<#
.SYNOPSIS
    Installs Orstax on Windows.
    Must be run as Administrator.
#>

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── Configuration ────────────────────────────────────────────────────────────
$REPO_URL   = "https://github.com/ortsax/whatsapp-bot.git"
$INSTALL_DIR = "$env:ProgramFiles\orstax"
$SRC_DIR     = "$INSTALL_DIR\src"
$BIN_PATH    = "$INSTALL_DIR\orstax.exe"
$GO_FALLBACK = "1.25.0"
# ─────────────────────────────────────────────────────────────────────────────

function Write-Step($msg) { Write-Host "`n==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    $msg" -ForegroundColor Green }
function Write-Err($msg)  { Write-Host "    ERROR: $msg" -ForegroundColor Red; exit 1 }

# ── Detect latest Go version ─────────────────────────────────────────────────
Write-Step "Detecting latest Go version"
try {
    $releases  = Invoke-RestMethod -Uri "https://go.dev/dl/?mode=json" -UseBasicParsing
    $GO_VERSION = ($releases | Where-Object { $_.stable -eq $true } | Select-Object -First 1).version -replace '^go',''
    Write-Ok "Latest stable Go: $GO_VERSION"
} catch {
    $GO_VERSION = $GO_FALLBACK
    Write-Ok "Could not fetch version, using fallback: $GO_VERSION"
}

# ── Check / install Git ──────────────────────────────────────────────────────
Write-Step "Checking Git"
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Err "Git is not installed. Please install Git from https://git-scm.com and re-run."
}
Write-Ok "Git found: $(git --version)"

# ── Check / install Go ───────────────────────────────────────────────────────
Write-Step "Checking Go"
$goInstalled = $false
if (Get-Command go -ErrorAction SilentlyContinue) {
    $goVer = (go version) -replace '.*go([0-9.]+).*','$1'
    Write-Ok "Go $goVer already installed"
    $goInstalled = $true
}

if (-not $goInstalled) {
    Write-Step "Installing Go $GO_VERSION"
    $goZip  = "$env:TEMP\go$GO_VERSION.windows-amd64.zip"
    $goRoot = "C:\Program Files\Go"
    $dlUrl  = "https://go.dev/dl/go$GO_VERSION.windows-amd64.zip"

    Write-Ok "Downloading $dlUrl"
    Invoke-WebRequest -Uri $dlUrl -OutFile $goZip -UseBasicParsing

    Write-Ok "Extracting to $goRoot"
    if (Test-Path $goRoot) { Remove-Item -Recurse -Force $goRoot }
    Expand-Archive -Path $goZip -DestinationPath "C:\Program Files" -Force
    Remove-Item $goZip

    # Add Go to system PATH for this session and permanently
    $gobin = "C:\Program Files\Go\bin"
    $env:PATH += ";$gobin"
    $syspath = [System.Environment]::GetEnvironmentVariable("PATH","Machine")
    if ($syspath -notlike "*$gobin*") {
        [System.Environment]::SetEnvironmentVariable("PATH", "$syspath;$gobin", "Machine")
    }
    Write-Ok "Go $GO_VERSION installed"
}

# ── Clone or update repo ─────────────────────────────────────────────────────
Write-Step "Setting up source"
if (Test-Path "$SRC_DIR\.git") {
    Write-Ok "Updating existing clone in $SRC_DIR"
    git -C $SRC_DIR pull
} else {
    Write-Ok "Cloning $REPO_URL into $SRC_DIR"
    New-Item -ItemType Directory -Force -Path $SRC_DIR | Out-Null
    # Clone into parent, targeting SRC_DIR name
    $parent = Split-Path $SRC_DIR -Parent
    $name   = Split-Path $SRC_DIR -Leaf
    git clone $REPO_URL "$parent\$name"
}

# ── Build ────────────────────────────────────────────────────────────────────
Write-Step "Building orstax"
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
$env:CGO_ENABLED = "0"
$ldflags = "-s -w -X main.sourceDir=$SRC_DIR"
Push-Location $SRC_DIR
go build -ldflags $ldflags -trimpath -o $BIN_PATH .
Pop-Location
Write-Ok "Binary written to $BIN_PATH"

# ── Add install dir to system PATH ───────────────────────────────────────────
Write-Step "Updating system PATH"
$syspath = [System.Environment]::GetEnvironmentVariable("PATH","Machine")
if ($syspath -notlike "*$INSTALL_DIR*") {
    [System.Environment]::SetEnvironmentVariable("PATH", "$syspath;$INSTALL_DIR", "Machine")
    $env:PATH += ";$INSTALL_DIR"
    Write-Ok "Added $INSTALL_DIR to system PATH"
} else {
    Write-Ok "Already in PATH"
}

# ── Done ─────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "Orstax is now installed" -ForegroundColor Green
Write-Host ""
Write-Host "  Run with    orstax --phone-number <international-number>"
Write-Host "  Update with orstax -update"
Write-Host "  Sessions    orstax -list-sessions"
Write-Host "              orstax -delete-session <phone>"
Write-Host "              orstax -reset-session  <phone>"
Write-Host ""
Write-Host "Note: open a new terminal for PATH changes to take effect."
