# ticker

[![Integration](https://github.com/systemli/ticker/actions/workflows/integration.yaml/badge.svg)](https://github.com/systemli/ticker/actions/workflows/integration.yaml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=systemli_ticker&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=systemli_ticker)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=systemli_ticker&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=systemli_ticker)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=systemli_ticker&metric=coverage)](https://sonarcloud.io/summary/new_code?id=systemli_ticker)
[![Docker Image](https://img.shields.io/docker/v/systemli/ticker?label=docker)](https://hub.docker.com/r/systemli/ticker)

The API for the [Systemli Ticker Project](https://www.systemli.org/en/service/ticker.html) — a
service to distribute short messages in support of events, demonstrations, or other time-sensitive
events.

This repository contains the backend. A complete installation also needs
[ticker-admin](https://github.com/systemli/ticker-admin) (admin interface) and
[ticker-frontend](https://github.com/systemli/ticker-frontend) (public page).

## Documentation

**<https://systemli.github.io/ticker/>**

The documentation covers the whole stack, for all three repositories:

- [Installation](https://systemli.github.io/ticker/installation/) — Docker Compose and Traefik
- [Configuration](https://systemli.github.io/ticker/configuration/)
- [Integrations](https://systemli.github.io/ticker/integrations/) — Telegram, Mastodon, Bluesky, Signal
- [Operations](https://systemli.github.io/ticker/operations/) — upgrades, backups, users
- [Troubleshooting](https://systemli.github.io/ticker/troubleshooting/)
- [Development](https://systemli.github.io/ticker/development/)

The sources live in [docs/](docs/).

## Quick start

```shell
git clone https://github.com/systemli/ticker.git
cd ticker
cp .env.example .env      # set the hostnames and TICKER_SECRET
docker compose up -d
docker compose run --rm ticker user create --email admin@example.org --super-admin
```

See the [installation guide](https://systemli.github.io/ticker/installation/) for DNS requirements
and the remaining setup — in particular registering your frontend's address on the ticker, which is
what makes the public page work.

## Contributing

Contributions are welcome. Fork the repository, commit your changes, and open a pull request.

[AGENTS.md](AGENTS.md) documents the project layout, code style, testing conventions and commit
format. For a local environment, see the
[development guide](https://systemli.github.io/ticker/development/):

```shell
docker compose -f compose.dev.yaml up -d --build   # the whole stack
go test ./...
```

## Licence

GPL-3.0. See [LICENSE](LICENSE).
