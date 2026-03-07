<div align="center">
  <img src="media/logo.png" alt="Alphonse" width="550"/>
  <h1>Alphonse</h1>
  <p>A self-hosted WhatsApp bot written in Go.</p>
  <a href="https://ortsax.github.io/Alphonse/"><strong>Full Documentation →</strong></a>
  &nbsp;·&nbsp;
  <a href="https://github.com/ortsax/Alphonse/issues">Report a Bug</a>
  &nbsp;·&nbsp;
  <a href="https://github.com/ortsax/Alphonse/releases">Releases</a>
  <br/><br/>
  <img src="https://img.shields.io/github/v/release/ortsax/Alphonse?style=flat&label=version" alt="Latest release"/>
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?style=flat&logo=docker" alt="Docker"/>
  <img src="https://img.shields.io/github/license/ortsax/Alphonse?style=flat" alt="License"/>
  <img src="https://img.shields.io/github/stars/ortsax/Alphonse?style=flat" alt="Stars"/>
</div>

---

Alphonse connects to WhatsApp via **phone-number pairing** (no QR code needed), persists sessions in SQLite or PostgreSQL, and ships a rich plugin system with moderation, group management, media conversion, AI integration, and more — all in a single statically-linked binary.

## Quick Start

### Docker (recommended)

The fastest way to get running. No Go toolchain or build step required.

```bash
# 1. Create a data directory and drop your config in it
mkdir data
cp .env.example data/.env   # then edit data/.env with your settings

# 2. Pull and run
docker compose up -d
```

Or without Compose:

```bash
docker build -t alphonse .
docker run -it -v "$(pwd)/data:/data" alphonse --phone-number <your-number>
```

### Download the Binary

Grab the pre-compiled binary for your platform from the [latest release](https://github.com/ortsax/Alphonse/releases/latest):

| Platform | File |
| -------- | ---- |
| Linux x86-64 | `alphonse_*_linux_amd64.tar.gz` |
| Linux ARM64 | `alphonse_*_linux_arm64.tar.gz` |
| macOS (Apple Silicon) | `alphonse_*_darwin_arm64.tar.gz` |
| macOS (Intel) | `alphonse_*_darwin_amd64.tar.gz` |
| Windows x86-64 | `alphonse_*_windows_amd64.zip` |

Extract the archive, place the binary on your `PATH`, then continue to [First Run](#first-run).

**Linux / macOS one-liner:**

```bash
curl -fsSL https://github.com/ortsax/Alphonse/releases/latest/download/alphonse_linux_amd64.tar.gz \
  | tar -xz && sudo mv alphonse /usr/local/bin/
```

## First Run

Pair your WhatsApp account once after installation:

```
alphonse --phone-number <international-format-number>
```

A pairing code is printed. On your phone open **WhatsApp → Linked Devices → Link a Device → Link with phone number instead** and enter the code. The session is saved — subsequent starts run the bot automatically.

## Usage

```
alphonse [flags]

Flags:
  --phone-number  <number>   Identify or pair a device
  --update                   Pull latest source and rebuild in-place
  --list-sessions            List all paired sessions
  --delete-session <number>  Permanently delete a session
  --reset-session  <number>  Reset a session for re-pairing
  --version                  Print version and exit
  -h, --help                 Show help
```

## Features

| Category        | Highlights                                                            |
| --------------- | --------------------------------------------------------------------- |
| **Moderation**  | Anti-link, anti-spam, anti-delete, anti-call, anti-word, warn system  |
| **Group Admin** | Promote/demote, kick, mute, create group (`newgc`)                    |
| **Media**       | Audio extraction (`mp3`), video trim, black-border removal            |
| **Status**      | Auto-save and auto-like WhatsApp status updates                       |
| **AI**          | Meta AI integration via `meta` command                                |
| **Settings**    | Per-owner config: prefixes, sudo users, public/private mode, language |
| **i18n**        | 10 built-in languages (EN, ES, PT, AR, HI, FR, DE, RU, TR, SW)        |
| **Updates**     | Self-update via `alphonse --update` or `.update` in chat              |

See the [command reference](https://ortsax.github.io/Alphonse/commands) for the full list.

## Documentation

The complete documentation is hosted on GitHub Pages:

**[ortsax.github.io/Alphonse](https://ortsax.github.io/Alphonse/)**

Topics covered:

- [Installation](https://ortsax.github.io/Alphonse/installation)
- [Configuration](https://ortsax.github.io/Alphonse/configuration)
- [Command Reference](https://ortsax.github.io/Alphonse/commands)
- [Plugin Development](https://ortsax.github.io/Alphonse/plugins)

<details>
<summary><strong>Building from Source</strong></summary>

Requires Go 1.25+ and Git.

**Linux / macOS**

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

**Manual**

```bash
git clone https://github.com/ortsax/Alphonse.git && cd Alphonse
make build            # produces alphonse / alphonse.exe
make release VERSION=0.0.1   # cross-platform archives → dist/
```

</details>

## Contributing

Contributions are **by invitation only**. If you would like to contribute a feature, fix a bug, or improve the documentation, please reach out here.

**[Contact here](mailto:danielpeter0081@gmail.com)**

## License

[MIT](LICENSE)

