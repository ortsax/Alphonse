#Requires -RunAsAdministrator

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── Configuration ────────────────────────────────────────────────────────────
$REPO_URL    = "https://github.com/ortsax/whatsapp-bot.git"
$INSTALL_DIR = "$env:ProgramFiles\orstax"
$BIN_PATH    = "$INSTALL_DIR\orstax.exe"
$TEMP_SRC    = "$env:TEMP\orstax_build_$(Get-Random)"
# ─────────────────────────────────────────────────────────────────────────────

function Write-Step($msg) { Write-Host "`n==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    $msg" -ForegroundColor Green }
function Write-Err($msg)  { Write-Host "    ERROR: $msg" -ForegroundColor Red; exit 1 }

# ── Setup Install Dir ────────────────────────────────────────────────────────
Write-Step "Preparing installation directory"
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
}

# ── Clone to Temp ────────────────────────────────────────────────────────────
Write-Step "Cloning repository to temporary folder"
if (Test-Path $TEMP_SRC) { Remove-Item -Recurse -Force $TEMP_SRC }
git clone $REPO_URL $TEMP_SRC

# ── Build ────────────────────────────────────────────────────────────────────
Write-Step "Building orstax"
$env:CGO_ENABLED = "0"

Push-Location $TEMP_SRC
# Simplified build command to avoid ldflags issues unless you explicitly need them
go build -trimpath -o $BIN_PATH .
Pop-Location

# ── Cleanup Source ───────────────────────────────────────────────────────────
Write-Step "Cleaning up temporary files"
Remove-Item -Recurse -Force $TEMP_SRC

if (Test-Path $BIN_PATH) {
    Write-Ok "Binary successfully installed to $BIN_PATH"
} else {
    Write-Err "Build failed: Binary not found."
}

# ── Update PATH ──────────────────────────────────────────────────────────────
Write-Step "Updating system PATH"
$syspath = [System.Environment]::GetEnvironmentVariable("PATH","Machine")
$normDir = $INSTALL_DIR.TrimEnd('\')

if ($syspath -split ';' -notcontains $normDir) {
    [System.Environment]::SetEnvironmentVariable("PATH", "$syspath;$normDir", "Machine")
    $env:PATH += ";$normDir"
    Write-Ok "Added $normDir to system PATH"
} else {
    Write-Ok "Already in PATH"
}

Write-Host "`nDone! Orstax is installed. Restart your terminal to use it." -ForegroundColor Green