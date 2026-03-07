# patch-whatsmeow.ps1
# Copies the current whatsmeow version from the module cache into ./patched/
# and applies all local patches.

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$root = Split-Path $PSScriptRoot -Parent

# ── 1. Resolve version from go.mod ──────────────────────────────────────────
Push-Location $root
$modLine = go list -m go.mau.fi/whatsmeow 2>&1
Pop-Location

if ($modLine -notmatch "go\.mau\.fi/whatsmeow\s+(\S+)") {
    Write-Error "Could not determine whatsmeow version from go.mod"
    exit 1
}
$version = $Matches[1]
Write-Host "whatsmeow version: $version"

# ── 2. Ensure module is in cache ─────────────────────────────────────────────
Push-Location $root
go mod download go.mau.fi/whatsmeow@$version | Out-Null
Pop-Location

$gopath = (go env GOPATH)
$src = "$gopath\pkg\mod\go.mau.fi\whatsmeow@$version"
if (-not (Test-Path $src)) {
    Write-Error "Module cache path not found: $src"
    exit 1
}

# ── 3. Replace ./patched with fresh copy ─────────────────────────────────────
$dst = Join-Path $root "patched"
if (Test-Path $dst) {
    Write-Host "Removing old patched/ ..."
    Remove-Item $dst -Recurse -Force
}
Write-Host "Copying $src -> patched/ ..."
Copy-Item $src $dst -Recurse -Force
# Make all files writable
Get-ChildItem $dst -Recurse | ForEach-Object {
    $_.Attributes = $_.Attributes -band (-bnot [System.IO.FileAttributes]::ReadOnly)
}

# ── 4. Patch send.go ─────────────────────────────────────────────────────────
$sendFile = Join-Path $dst "send.go"
$sendContent = [System.IO.File]::ReadAllText($sendFile)

# 4a. add PinInChatMessage to getEditAttribute
if ($sendContent -notlike "*EditAttributePinInChat*") {
    # Anchor: closing of the switch + default return — unique in getEditAttribute
    # Handles both LF and CRLF by matching \r?\n
    $patched = $sendContent -replace `
        '(\r?\n\t\}(\r?\n)\treturn types\.EditAttributeEmpty(\r?\n)\})', `
        "`n`tcase msg.PinInChatMessage != nil:`n`t`treturn types.EditAttributePinInChat`$1"
    if ($patched -eq $sendContent) {
        Write-Warning "send.go: PinInChatMessage anchor not found — patch may need updating"
    } else {
        $sendContent = $patched
        Write-Host "Patched: send.go (PinInChatMessage)"
    }
} else {
    Write-Host "send.go (PinInChatMessage) already patched"
}

# 4b. Release messageSendLock before the server-ACK wait so queued sends don't
#     stall behind a 100-500 ms round-trip.  Replace the bare `defer Unlock`
#     with a one-shot helper and call it right after the message hits the wire.
if ($sendContent -notlike "*doUnlock*") {
    # Replace the defer line
    $patched = $sendContent -replace `
        'resp\.DebugTimings\.Queue = time\.Since\(start\)\r?\n\tdefer cli\.messageSendLock\.Unlock\(\)', `
        ("resp.DebugTimings.Queue = time.Since(start)`n" +
         "`t// The lock must be held during encrypt+write so Signal session state`n" +
         "`t// advances correctly. It does NOT need to be held while waiting for the`n" +
         "`t// server ACK, which is pure network I/O. Releasing early lets the next`n" +
         "`t// queued send start its encrypt+write immediately in parallel with our ACK`n" +
         "`t// wait, cutting queue latency from (N x RTT) to (N x encrypt_ms + 1 RTT).`n" +
         "`tunlocked := false`n" +
         "`tdoUnlock := func() {`n" +
         "`t`tif !unlocked {`n" +
         "`t`t`tunlocked = true`n" +
         "`t`t`tcli.messageSendLock.Unlock()`n" +
         "`t`t}`n" +
         "`t}`n" +
         "`tdefer doUnlock()")
    if ($patched -eq $sendContent) {
        Write-Warning "send.go: messageSendLock defer anchor not found — patch may need updating"
    } else {
        $sendContent = $patched
        Write-Host "Patched: send.go (early unlock)"
    }

    # Insert doUnlock() call right before 'var respNode *waBinary.Node' (after send, before ACK wait)
    $patched = $sendContent -replace `
        '(\tcli\.cancelResponse\(req\.ID, respChan\)\r?\n\t\treturn\r?\n\t\}\r?\n)\tvar respNode \*waBinary\.Node', `
        "`$1`t// Message is written to the socket. Unlock now so the next sender can`n`t// begin its encrypt+write while we await the server ACK.`n`tdoUnlock()`n`tvar respNode *waBinary.Node"
    if ($patched -eq $sendContent) {
        Write-Warning "send.go: early-unlock insertion anchor not found — patch may need updating"
    } else {
        $sendContent = $patched
        Write-Host "Patched: send.go (doUnlock call before ACK wait)"
    }

    [System.IO.File]::WriteAllText($sendFile, $sendContent)
} else {
    Write-Host "send.go (early unlock) already patched"
}

# ── 5. Patch appstate/encode.go — add BuildClearChat ─────────────────────────
$encodeFile = Join-Path $dst "appstate\encode.go"
$encodeContent = [System.IO.File]::ReadAllText($encodeFile)

$clearFunc = @'

// BuildClearChat builds an app state patch for clearing a chat's message history.
func BuildClearChat(target types.JID, lastMessageTimestamp time.Time, lastMessageKey *waCommon.MessageKey) PatchInfo {
	return PatchInfo{
		Type: WAPatchRegularHigh,
		Mutations: []MutationInfo{{
			Index:   []string{IndexClearChat, target.String(), "1", "0"},
			Version: 6,
			Value: &waSyncAction.SyncActionValue{
				ClearChatAction: &waSyncAction.ClearChatAction{
					MessageRange: newMessageRange(lastMessageTimestamp, lastMessageKey),
				},
			},
		}},
	}
}
'@

if ($encodeContent -notlike "*BuildClearChat*") {
    # Insert after BuildDeleteChat's closing brace (before newMessageRange)
    $anchor = "func newMessageRange("
    if ($encodeContent.Contains($anchor)) {
        $encodeContent = $encodeContent.Replace($anchor, $clearFunc + "`n" + $anchor)
        [System.IO.File]::WriteAllText($encodeFile, $encodeContent)
        Write-Host "Patched: appstate/encode.go (BuildClearChat)"
    } else {
        Write-Warning "appstate/encode.go: anchor not found — patch may need updating"
    }
} else {
    Write-Host "appstate/encode.go already patched"
}

Write-Host "`nAll patches applied. Run 'go build ./...' to verify."

# ── 6. Patch store/clientpayload.go — add SetAndroidMode ─────────────────────
$cpFile = Join-Path $dst "store\clientpayload.go"
$cpContent = [System.IO.File]::ReadAllText($cpFile)

$androidFunc = @'

// SetAndroidMode configures the client payload to identify as an Android companion,
// matching Baileys' Browsers.android(name) setup. name is shown on the linked devices list.
func SetAndroidMode(name string) {
	DeviceProps.Os = proto.String(name)
	DeviceProps.PlatformType = waCompanionReg.DeviceProps_ANDROID_PHONE.Enum()
	BaseClientPayload.UserAgent.Platform = waWa6.ClientPayload_UserAgent_ANDROID.Enum()
	BaseClientPayload.UserAgent.OsVersion = proto.String("")
	BaseClientPayload.UserAgent.OsBuildNumber = proto.String("")
	BaseClientPayload.WebInfo = nil
}

'@

if ($cpContent -notlike "*SetAndroidMode*") {
    $anchor = "func SetOSInfo("
    if ($cpContent.Contains($anchor)) {
        $cpContent = $cpContent.Replace($anchor, $androidFunc + $anchor)
        [System.IO.File]::WriteAllText($cpFile, $cpContent)
        Write-Host "Patched: store/clientpayload.go (SetAndroidMode)"
    } else {
        Write-Warning "store/clientpayload.go: anchor not found — patch may need updating"
    }
} else {
    Write-Host "store/clientpayload.go already patched"
}

Write-Host "`nAll patches applied. Run 'go build ./...' to verify."
