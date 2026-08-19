# Systemli Ticker

**A service to distribute short messages in support of events, demonstrations, or other
time-sensitive events.**

A ticker is a single public page that updates live. Editors post short messages from an admin
interface; readers see them appear without reloading. Messages can additionally be pushed to
Telegram, Mastodon, Bluesky and Signal groups.

See the [Systemli Ticker Project page](https://www.systemli.org/en/service/ticker.html) for the
hosted service.

## Components

A complete installation is three services, published as three Docker images:

| Component | Image | Role |
| --- | --- | --- |
| [ticker](https://github.com/systemli/ticker) | `systemli/ticker` | The API. Stores everything, serves the public endpoints and media, dispatches to integrations. It needs no public hostname of its own. |
| [ticker-admin](https://github.com/systemli/ticker-admin) | `systemli/ticker-admin` | Admin interface. Editors log in here to manage tickers, messages and users. |
| [ticker-frontend](https://github.com/systemli/ticker-frontend) | `systemli/ticker-frontend` | The public page your readers visit. |

Both interfaces are static single-page applications. They talk to the API over HTTP and hold no
data of their own.

## How a request flows

```
   readers         ┌─────────────────────────────┐
   feeds, ────────▶│  ticker.example.org         │
   media           │  ticker-frontend (SPA)      │
                   │  /api/**  ──────────────────┼──┐
                   └─────────────────────────────┘  │      ┌──────────────┐
                                                    ├─────▶│  ticker      │
                   ┌─────────────────────────────┐  │      │  (API)       │
   editors ───────▶│  admin.ticker.example.org   │  │      │              │
                   │  ticker-admin (SPA)         │  │      └──────┬───────┘
                   │  /api/**  ──────────────────┼──┘             │
                   └─────────────────────────────┘         ┌──────▼───────┐
                                                           │  PostgreSQL  │
                                                           └──────────────┘
```

Two details of this shape matter, and explain most of the configuration:

- **The API has no public hostname of its own.** Everything it serves — the public endpoints,
  attachments, RSS feeds — lives below `/v1`, and both interfaces expose that as `/api` on their own
  address. Attachment URLs in API responses are relative for the same reason.
- **The API works out which ticker a request is for from the browser's `Origin` header.** That is
  why the public frontend's address must be registered on the ticker itself, and why the reverse
  proxy has to pass a correct `Origin` along. See [Installation](installation.md).

## Where to go next

- **[Installation](installation.md)** — run the whole stack with Docker Compose and Traefik.
- **[Configuration](configuration.md)** — every setting the API understands.
- **[Integrations](integrations.md)** — Telegram, Mastodon, Bluesky, Signal.
- **[Operations](operations.md)** — upgrades, backups, users, monitoring.
- **[Troubleshooting](troubleshooting.md)** — symptoms and their causes.
- **[Development](development.md)** — work on the code.

## Licence

GPL-3.0.
