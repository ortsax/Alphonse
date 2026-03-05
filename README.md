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
  <img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go" alt="Go version"/>
  <img src="https://img.shields.io/github/license/ortsax/Alphonse?style=flat" alt="License"/>
  <img src="https://img.shields.io/github/stars/ortsax/Alphonse?style=flat" alt="Stars"/>
</div>

---

Alphonse connects to WhatsApp via **phone-number pairing** (no QR code needed), persists sessions in SQLite or PostgreSQL, and ships a rich plugin system with moderation, group management, media conversion, AI integration, and more — all in a single statically-linked binary.

> **New here?** Read the [full documentation](https://ortsax.github.io/Alphonse/) for installation guides, configuration reference, and the complete command list.

## Quick Install

Pick the one-liner for your platform. Each script installs Go and Git if missing, clones the repo, builds the binary, and adds it to your PATH.

**Windows** — PowerShell as Administrator

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force
irm https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install.ps1 | iex
```

**Linux**

```bash
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install-linux.sh)
```

**macOS**

```bash
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/Alphonse/master/scripts/install-mac.sh)
```

Once the script finishes, pair your WhatsApp account:

```
alphonse --phone-number <international-format-number>
```

A pairing code is printed. On your phone open **WhatsApp → Linked Devices → Link a Device → Link with phone number instead** and enter the code.

## Usage

```
alphonse [flags]

Flags:
  --phone-number  <number>   Identify or pair a device
  --update                   Pull latest source and rebuild in-place
  --list-sessions            List all paired sessions
  --delete-session <number>  Permanently delete a session
  --reset-session  <number>  Reset a session for re-pairing
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

## Contributing

Contributions are **by invitation only**. If you would like to contribute a feature, fix a bug, or improve the documentation, please reach out here.

**[Contact here](mailto:danielpeter0081@gmail.com)**

## License

[MIT](LICENSE)
