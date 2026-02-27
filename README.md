# whatsapp-bot

A WhatsApp bot built in Go using [whatsmeow](https://github.com/tulir/whatsmeow). Connects via phone-number pairing (no QR scan), persists sessions in SQLite or PostgreSQL, and supports an extensible command plugin system.

## Installation

Pick the script for your platform and run it with elevated privileges. The script handles everything — Go, Git, cloning, building, and PATH setup.

**Windows** (PowerShell as Administrator)
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force
irm https://raw.githubusercontent.com/ortsax/whatsapp-bot/main/scripts/install.ps1 | iex
```

**Linux**
```bash
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/whatsapp-bot/main/scripts/install-linux.sh)
```

**macOS**
```bash
sudo bash <(curl -fsSL https://raw.githubusercontent.com/ortsax/whatsapp-bot/main/scripts/install-mac.sh)
```

Once complete you will see:
```
Orstax is now installed

  Run with    orstax --phone-number <international-number>
  Update with orstax -update
```

> Open a new terminal after install so the updated PATH takes effect.

## First run — pairing

```bash
orstax --phone-number <international-format-number>
```

A pairing code will be printed. On your phone go to **WhatsApp → Linked Devices → Link a Device → Link with phone number instead** and enter the code.

## Subsequent runs

```bash
orstax
```

Press `Ctrl+C` to disconnect.

## Database

By default Orstax uses a local SQLite file. To use PostgreSQL, create a `.env` file next to the binary:

```env
# SQLite (default)
DATABASE_URL=database.db

# PostgreSQL
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
```

## Session management

```bash
orstax -list-sessions                  # list all paired sessions
orstax -delete-session <phone>         # permanently remove a session
orstax -reset-session  <phone>         # remove a session so it can be re-paired
```

## Updating

```bash
orstax -update
```

Pulls the latest source and rebuilds the binary in-place. Stop the bot first on Windows before updating.

## Commands

### Utility

| Command | Alias | Description |
|---|---|---|
| `menu` | `help` | Shows all commands grouped by category |
| `ping` | — | Replies with round-trip latency |

> Typing a category name as a command (e.g. `.settings`) shows only that category's commands.

### AI

| Command | Alias | Description |
|---|---|---|
| `meta <query>` | `ai` | Sends a query to Meta AI and streams the response back |

### Settings *(sudo only)*

| Command | Description |
|---|---|
| `setprefix <p1> <p2> …` | Sets command prefix(es). Use `empty` for no-prefix |
| `setsudo add\|remove <phone>` | Grants or revokes sudo access |
| `setmode public\|private` | Public: anyone can use commands. Private: sudo users only |
| `lang [code]` | Shows current language / available languages, or switches language |

## Languages

The bot ships with 10 languages. Switch with `.lang <code>` (sudo only).

| Code | Language |
|---|---|
| `en` | English (default) |
| `es` | Español |
| `pt` | Português |
| `ar` | العربية |
| `hi` | हिन्दी |
| `fr` | Français |
| `de` | Deutsch |
| `ru` | Русский |
| `tr` | Türkçe |
| `sw` | Kiswahili |

## Project structure

```
main.go               – entry point, device pairing, session management, updater
scripts/
  install.ps1         – Windows installer
  install-linux.sh    – Linux installer
  install-mac.sh      – macOS installer
plugins/
  command.go          – command registry (O(1) lookup), dispatch, Context
  handler.go          – whatsmeow event handler (fully non-blocking)
  menu.go             – menu & category-menu commands
  ping.go             – ping command
  meta.go             – Meta AI integration
  settings.go         – in-memory Settings, SQLite/Postgres persistence
  settings_cmds.go    – setprefix / setsudo / setmode commands
  lang_cmd.go         – lang command
  i18n.go             – translations for all 10 languages
  users.go            – LID ↔ phone mapping helpers
store/                – whatsmeow store package
```

## License

See [LICENSE](LICENSE).
