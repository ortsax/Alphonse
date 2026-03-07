<div align="center">
  <img src="media/logo.png" alt="Alphonse" width="550"/>
  <h1>Alphonse</h1>
  <p>A self-hosted WhatsApp bot written in Go.</p>
  <a href="https://ortsax.github.io/Alphonse/"><strong>Documentation →</strong></a>
  &nbsp;·&nbsp;
  <a href="https://github.com/ortsax/Alphonse/releases/latest">Latest Release</a>
  &nbsp;·&nbsp;
  <a href="https://github.com/ortsax/Alphonse/releases/tag/nightly">Nightly Build</a>
  &nbsp;·&nbsp;
  <a href="https://github.com/ortsax/Alphonse/issues">Report a Bug</a>
  <br/><br/>
  <img src="https://img.shields.io/github/v/release/ortsax/Alphonse?style=flat&label=version" alt="Latest release"/>
  <img src="https://img.shields.io/badge/Docker-ready-2496ED?style=flat&logo=docker" alt="Docker"/>
  <img src="https://img.shields.io/github/license/ortsax/Alphonse?style=flat" alt="License"/>
  <img src="https://img.shields.io/github/stars/ortsax/Alphonse?style=flat" alt="Stars"/>
</div>

---

Alphonse connects to WhatsApp via **phone-number pairing**, runs as a single binary, and features wa moderation, media tools, AI integration, and a plugin system.

## Quick Start

**1 — Run with Docker**

```bash
mkdir data && cp .env.example data/.env   # edit data/.env first
docker compose up -d
```

**2 — Or download the binary** from the [latest release](https://github.com/ortsax/Alphonse/releases/latest) for your platform, then run it directly. Data is stored in `~/Documents/Alphonse Files` automatically.

**3 — Pair your phone** (first run only)

```
alphonse --phone-number <international-number>
```

WhatsApp → Linked Devices → Link with phone number → enter the printed code.

**4 — Done.** Subsequent starts reconnect automatically.

> See the [full installation guide](https://ortsax.github.io/Alphonse/installation) for Docker details, session management, and updating.

## Usage

```
alphonse [flags]

  --phone-number  <number>   Pair or identify a device
  --update                   Pull latest source and rebuild in-place
  --list-sessions            List all paired sessions
  --delete-session <number>  Permanently delete a session
  --reset-session  <number>  Reset a session for re-pairing
  --version                  Print version and exit
  -h, --help                 Show help
```

## Features

| Category        | Highlights                                                           |
| --------------- | -------------------------------------------------------------------- |
| **Moderation**  | Anti-link, anti-spam, anti-delete, anti-call, anti-word, warn system |
| **Group Admin** | Promote/demote, kick, mute, create group                             |
| **Media**       | Audio extraction (`mp3`), video trim, black-border removal           |
| **Status**      | Auto-save and auto-like WhatsApp status updates                      |
| **AI**          | Meta AI via `meta` command                                           |
| **Settings**    | Prefix, language, mode, sudo users — all changeable from chat        |
| **i18n**        | 10 languages: EN ES PT AR HI FR DE RU TR SW                         |
| **Updates**     | `alphonse --update` or `.update` in chat                             |

Full command reference at **[ortsax.github.io/Alphonse/commands](https://ortsax.github.io/Alphonse/commands)**.

## Documentation

**[ortsax.github.io/Alphonse](https://ortsax.github.io/Alphonse/)** — installation, configuration, commands, plugin development.

<details>
<summary><strong>Building from Source</strong></summary>

Requires Go 1.25+ and Git. The `patched/` directory (whatsmeow fork) is included in the repo.

```bash
git clone https://github.com/ortsax/Alphonse.git && cd Alphonse
make build              # alphonse.exe / alphonse
make release VERSION=x.y.z   # cross-platform archives → dist/
```

Install scripts (also build from source):

```bash
# Linux
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install-linux.sh)
# macOS
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install-mac.sh)
# Windows (PowerShell as Administrator)
irm https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install.ps1 | iex
```

</details>

## Contributing

Contributions are **by invitation only** — [get in touch](mailto:danielpeter0081@gmail.com) if you'd like to help.

## License

[MIT](LICENSE)


