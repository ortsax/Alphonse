---
layout: home
title: Home
nav_order: 1
---

# Alphonse

A self-hosted WhatsApp bot written in Go, built on [whatsmeow](https://github.com/tulir/whatsmeow). Connects via **phone-number pairing** — no QR code, no phone scan — and ships as a single statically-linked binary.

## What it can do

| Category | Highlights |
|---|---|
| **Moderation** | Anti-link, anti-spam, anti-delete, anti-call, anti-word, warn system |
| **Group Admin** | Promote/demote, kick, mute, create groups |
| **Media** | Audio extraction (MP3), video trim, black-border removal |
| **Status** | Auto-save and auto-like contact status updates |
| **AI** | Meta AI integration via `.meta` command |
| **Settings** | Prefix, language, mode, sudo users — all changeable from chat |
| **i18n** | 10 languages: EN ES PT AR HI FR DE RU TR SW |
| **Updates** | `alphonse --update` or `.update` in chat |

## Pages

| | |
|---|---|
| [Installation](installation) | Docker, binary download, pairing, session management |
| [Configuration](configuration) | Environment variables, database, runtime settings |
| [Command Reference](commands) | Every command with usage examples |
| [Plugin Development](plugins) | Writing and registering custom commands |

## Requirements

- A WhatsApp account with an active phone number
- Docker **or** a pre-built binary from the [releases page](https://github.com/ortsax/Alphonse/releases)
- Go 1.25+ only needed if building from source

## License

[MIT](https://github.com/ortsax/Alphonse/blob/master/LICENSE)
