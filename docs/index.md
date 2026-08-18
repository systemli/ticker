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
| [ticker](https://github.com/systemli/ticker) | `systemli/ticker` | The API. Stores everything, serves the public endpoints and media, dispatches to integrations. |
| [ticker-admin](https://github.com/systemli/ticker-admin) | `systemli/ticker-admin` | Admin interface. Editors log in here to manage tickers, messages and users. |
| [ticker-frontend](https://github.com/systemli/ticker-frontend) | `systemli/ticker-frontend` | The public page your readers visit. |

Both interfaces are static single-page applications. They talk to the API over HTTP and hold no
data of their own.

## How a request flows

```
                        ┌─────────────────────────────┐
   readers ────────────▶│  ticker.example.org         │
                        │  ticker-frontend (SPA)      │
                        │  /api/**  ──────────────────┼──┐
                        └─────────────────────────────┘  │
                                                         │
                        ┌─────────────────────────────┐  │   ┌──────────────┐
   editors ────────────▶│  admin.ticker.example.org   │  ├──▶│  ticker      │
                        │  ticker-admin (SPA)         │  │   │  (API)       │
                        │  /api/**  ──────────────────┼──┤   │              │
                        └─────────────────────────────┘  │   └──────┬───────┘
                                                         │          │
                        ┌─────────────────────────────┐  │   ┌──────▼───────┐
   feeds, media ───────▶│  api.ticker.example.org     │──┘   │  PostgreSQL  │
                        │  /v1/**, /media/**          │      └──────────────┘
                        └─────────────────────────────┘
```

Two details of this shape matter, and explain most of the configuration:

- **The API needs its own public hostname.** Attachment URLs are absolute and served by the API
  at `/media/...`, outside the `/v1` prefix. RSS readers also fetch feeds directly.
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
