---
layout: page
title: Configuration
nav_order: 3
---

# Configuration
{: .no_toc }

<details open markdown="block">
  <summary>Table of contents</summary>
  {: .text-delta }
- TOC
{:toc}
</details>

---

## Environment variables

Alphonse loads `.env` from the current directory first, then from `~/Documents/Alphonse Files/.env`, then falls back to `.env.example`. You can also export variables directly in your shell or service file.

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `~/Documents/Alphonse Files/database.db` | SQLite file path **or** `postgres://…` connection string |

### SQLite (default)

Any value that does not start with `postgres://` or `postgresql://` is treated as a SQLite file path. If `DATABASE_URL` is not set, the database is stored in `~/Documents/Alphonse Files/database.db` (created automatically).

```env
# Override the default location
DATABASE_URL=/var/data/alphonse.db
```

Alphonse applies WAL mode, busy-timeout, and other recommended pragmas automatically.

### PostgreSQL

```env
DATABASE_URL=postgres://user:password@localhost:5432/alphonse?sslmode=disable
```

### Docker

When running in Docker, mount a volume at `/data` and set `DATABASE_URL` to keep data outside the container:

```yaml
# docker-compose.yml
services:
  alphonse:
    volumes:
      - ./data:/data
```

```env
# data/.env
DATABASE_URL=/data/database.db
```

---

## Runtime settings

Bot settings are stored per-owner in the database and can be changed from chat at any time without restarting.

### Prefix

The character(s) that must precede a command name.

```
.setprefix .
.setprefix . ! empty    # multiple prefixes; "empty" allows no prefix
```

Default: `.`

### Mode

| Value | Who can use commands |
|---|---|
| `public` | Everyone (default) |
| `private` | Sudo users only |

```
.setmode public
.setmode private
```

### Sudo users

Sudo users bypass all permission checks. The owner is added automatically on first start.

```
.setsudo <phone>    # grant sudo
.delsudo <phone>    # revoke sudo
.getsudo            # list all sudo users
```

### Language

```
.lang en   .lang es   .lang pt   .lang ar   .lang hi
.lang fr   .lang de   .lang ru   .lang tr   .lang sw
```

### Other toggles

```
.shh / .shh off          # disable / re-enable group responses
.antidelete on/off        # forward deleted messages to owner
.ban <phone>              # silently ignore a user
.unban <phone>
.disable <command>        # disable a specific command by name
.enable  <command>
```

---

## Build-time flags

When building from source, these ldflags are injected automatically by `make build` and `make release`:

| Flag | Description |
|---|---|
| `-X main.Version=x.y.z` | Version string shown by `--version` |
| `-X main.Commit=<sha>` | Short git commit hash |
| `-X main.BuildDate=<rfc3339>` | Build timestamp |
| `-X main.sourceDir=<path>` | Source directory used by `--update` |

Manual example:

```bash
go build \
  -ldflags="-s -w -X main.Version=1.0.0 -X main.sourceDir=$(pwd)" \
  -trimpath -o alphonse .
```
