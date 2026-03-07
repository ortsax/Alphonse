---
layout: page
title: Installation
nav_order: 2
---

# Installation
{: .no_toc }

<details open markdown="block">
  <summary>Table of contents</summary>
  {: .text-delta }
- TOC
{:toc}
</details>

---

## Releases

Two release channels are available:

| Channel | When to use | Link |
|---|---|---|
| **Stable** | Production — tested, versioned | [Latest release](https://github.com/ortsax/Alphonse/releases/latest) |
| **Nightly** | Try the newest features (may be unstable) | [Nightly build](https://github.com/ortsax/Alphonse/releases/tag/nightly) |

---

## Docker (recommended)

No Go toolchain required. The image downloads the pre-built binary and bundles `ffmpeg`.

```bash
# 1. Create a data directory and add your config
mkdir data
cp .env.example data/.env    # then edit data/.env

# 2. Start the bot
docker compose up -d
```

To pair a phone number on first run:

```bash
docker compose run --rm alphonse --phone-number <number>
```

To upgrade to a newer release, update the `VERSION` arg in `docker-compose.yml` and rebuild:

```bash
docker compose build --no-cache && docker compose up -d
```

---

## Download the Binary

Download the archive for your platform from the [latest release](https://github.com/ortsax/Alphonse/releases/latest):

| Platform | File |
|---|---|
| Linux x86-64 | `alphonse_*_linux_amd64.tar.gz` |
| Linux ARM64 | `alphonse_*_linux_arm64.tar.gz` |
| macOS (Apple Silicon) | `alphonse_*_darwin_arm64.tar.gz` |
| macOS (Intel) | `alphonse_*_darwin_amd64.tar.gz` |
| Windows x86-64 | `alphonse_*_windows_amd64.zip` |

Extract the archive and place the binary on your `PATH`.

**Linux / macOS one-liner:**

```bash
curl -fsSL https://github.com/ortsax/Alphonse/releases/latest/download/alphonse_linux_amd64.tar.gz \
  | tar -xz && sudo mv alphonse /usr/local/bin/
```

### Data directory

When using SQLite (the default), Alphonse stores `database.db` and reads `.env` from:

```
~/Documents/Alphonse Files/
```

This directory is created automatically on first run. To use a different location, set `DATABASE_URL` — see [Configuration](configuration#environment-variables).

---

## Pairing your phone

Run this once after installation:

```
alphonse --phone-number <international-number>
```

`<number>` must be in international format without the leading `+` — for example `2348012345678`.

Alphonse prints a pairing code:

```
Your pairing code is: ABCD-1234
```

On your phone open **WhatsApp → Settings → Linked Devices → Link a Device → Link with phone number instead** and enter the code. The session is saved automatically — subsequent runs reconnect without a code.

---

## Session management

| Command | Effect |
|---|---|
| `alphonse --list-sessions` | List all paired sessions |
| `alphonse --delete-session <number>` | Permanently remove a session |
| `alphonse --reset-session <number>` | Clear a session so it can be re-paired |

---

## Updating

```
alphonse --update
```

Pulls the latest source from GitHub, rebuilds the binary in-place, and exits. Restart the bot to run the new version. The same operation is available as the `.update` chat command.

> Docker users: update by changing the version in `docker-compose.yml` and rebuilding — `--update` is not applicable inside a container.

---

## Building from source

Requires Go 1.25+ and Git. The `patched/` directory (a local whatsmeow fork) is included in the repo.

```bash
git clone https://github.com/ortsax/Alphonse.git
cd Alphonse
make build          # alphonse.exe (Windows) or alphonse (Linux/macOS)
```

One-liner install scripts are also available for each platform — they clone the repo, build, and add the binary to `PATH`:

**Linux**
```bash
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install-linux.sh)
```

**macOS**
```bash
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install-mac.sh)
```

**Windows** (PowerShell as Administrator)
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force
irm https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install.ps1 | iex
```
