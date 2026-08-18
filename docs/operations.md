# Operations

## Managing users

User management is a CLI concern; there is no self-registration.

```shell
# Create a user. Without --password one is generated and printed.
docker compose run --rm ticker user create --email editor@example.org

# Create an administrator
docker compose run --rm ticker user create --email admin@example.org --super-admin

# Change a password
docker compose run --rm ticker user password --email editor@example.org --password NEW

# Delete a user
docker compose run --rm ticker user delete --email editor@example.org
```

Super admins may manage tickers, users and integration settings. Regular users only see the tickers
they were assigned to.

!!! note

    These commands need the database, so they only work while it is running. The same applies to
    `ticker version`.

## Upgrading

Pin image tags in `.env` rather than tracking `latest`, so upgrades are deliberate:

```shell title=".env"
TICKER_TAG=3.2.1
ADMIN_TAG=3.15.0
FRONTEND_TAG=3.4.0
```

Then:

```shell
docker compose pull
docker compose up -d
```

Database migrations run automatically at startup. Take a backup first, since there is no automatic
downgrade path.

Releases are listed on the [releases page](https://github.com/systemli/ticker/releases). The API and
the two interfaces are versioned independently.

## Backups

Two things must be saved **together**, because attachment records in the database point at files on
disk:

```shell
# Database
docker compose exec -T postgres pg_dump -U ticker ticker | gzip > ticker-db.sql.gz

# Uploaded files
docker run --rm \
  -v ticker_ticker-data:/data:ro \
  -v "$PWD":/backup \
  alpine:3.24 tar czf /backup/ticker-uploads.tar.gz -C /data .
```

Keep `.env` too — losing `TICKER_SECRET` logs everyone out, and losing the database password makes
the dump useless.

Restore the database into an empty instance with:

```shell
gunzip -c ticker-db.sql.gz | docker compose exec -T postgres psql -U ticker ticker
```

!!! warning "PostgreSQL major upgrades"

    Moving between PostgreSQL major versions is not a matter of changing the image tag; the data
    directory format differs. Dump with the old version, then restore into the new one.

## Health and monitoring

```shell
curl https://ticker.example.org/healthz
# OK
```

Prometheus metrics are exposed on a **separate** listener, port `8181` by default, at `/metrics`.
The compose stack keeps this on the internal network — it is not routed publicly. To scrape it, put
your collector on the same network:

```shell
docker compose exec ticker-init sh -c 'wget -qO- http://ticker:8181/metrics' | head
```

The published API image is built `FROM scratch` and contains no shell, so it cannot carry a Docker
`HEALTHCHECK`. The stack instead lets Traefik poll `/healthz`, which needs nothing inside the
container.

## Logs

```shell
docker compose logs -f ticker
```

Control verbosity and shape with `TICKER_LOG_LEVEL` (`debug`, `info`, `warn`, `error`) and
`TICKER_LOG_FORMAT` (`json` for log collectors, `text` when reading by eye). Every line carries a
`package` field, and integration failures carry `bridge_name`.

## Restarts and shutdown

The API handles `SIGTERM`, so `docker compose stop`, `restart` and `up -d` shut it down gracefully:
open WebSocket connections are closed and in-flight requests are given up to five seconds to finish.

Provided `TICKER_SECRET` is set, restarts do not log anyone out.

## Resetting a ticker

The admin interface can reset a ticker, which deletes its messages and disconnects its integrations
while keeping the ticker and its configuration. This cannot be undone.
