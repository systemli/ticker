# Configuration

The API reads its configuration from three places, later ones winning:

1. built-in defaults,
2. a YAML file, passed with `--config /path/to/config.yml`,
3. environment variables.

For container deployments, environment variables alone are usually enough.

!!! note "Integrations are not configured here"

    Telegram, Mastodon, Bluesky and Signal are configured at runtime through the admin interface
    and stored in the database — not in this file and not through environment variables. See
    [Integrations](integrations.md).

## Settings

| YAML key | Environment variable | Default | Description |
| --- | --- | --- | --- |
| `listen` | `TICKER_LISTEN` | `:8080` | Address and port for the API. |
| `log_level` | `TICKER_LOG_LEVEL` | `debug` | `debug`, `info`, `warn` or `error`. |
| `log_format` | `TICKER_LOG_FORMAT` | `json` | `json` or `text`. |
| `secret` | `TICKER_SECRET` | *randomly generated* | Signing secret for JSON Web Tokens. **Always set this.** |
| `database.type` | `TICKER_DATABASE_TYPE` | `sqlite` | `postgres`, `mysql` or `sqlite`. |
| `database.dsn` | `TICKER_DATABASE_DSN` | `ticker.db` | Connection string, see below. |
| `metrics_listen` | `TICKER_METRICS_LISTEN` | `:8181` | Address for the Prometheus exporter, on a separate listener. |
| `upload.path` | `TICKER_UPLOAD_PATH` | `uploads` | Directory for uploaded files. |
| `upload.url` | `TICKER_UPLOAD_URL` | `http://localhost:8080` | Public base URL used to build attachment links. |

That is the complete list. There is no environment variable for any setting not named above.

### Example file

```yaml title="config.yml"
--8<-- "config.yml.dist"
```

Run it with:

```shell
ticker run --config config.yml
```

!!! warning "`listen` in containers"

    If you mount a config file into a container, `listen` must bind all interfaces (`:8080`).
    A value like `localhost:8080` makes the API unreachable from outside the container.

## The signing secret

```shell
openssl rand -hex 32
```

When `secret` is unset the API generates a random one **on every start**. Tokens issued before a
restart become invalid, so everyone is logged out whenever the process restarts. Set it once and
treat it like a password: changing it later invalidates all existing sessions.

There is no `TICKER_SECRET_FILE` mechanism, so Docker or Swarm secrets cannot be injected as files.
Either pass the value as an environment variable or mount a full `config.yml`.

## Database

### PostgreSQL (recommended)

```
host=postgres port=5432 user=ticker password=SECRET dbname=ticker sslmode=disable TimeZone=UTC
```

The URL form works as well:

```
postgres://ticker:SECRET@postgres:5432/ticker?sslmode=disable
```

`sslmode=disable` is fine when the database is only reachable over a private container network.

!!! warning "Use `TimeZone=UTC`, not a named zone"

    The published image contains no timezone database, so a value like `TimeZone=Etc/UTC` fails at
    startup with `unknown time zone Etc/UTC`. `UTC` is built in and always works. All timestamps
    are UTC regardless.

### MySQL / MariaDB

```
ticker:SECRET@tcp(mysql:3306)/ticker?charset=utf8mb4&parseTime=True&loc=Local
```

`parseTime=True` is required.

### SQLite

!!! danger "SQLite does not work in the official image or in released binaries"

    The SQLite driver needs cgo, and the released builds are compiled without it. Starting the
    container with `TICKER_DATABASE_TYPE=sqlite` fails immediately:

    ```
    could not connect to database
    error="Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub"
    ```

    Use PostgreSQL or MySQL for any Docker deployment. SQLite is only usable when you build the
    binary yourself with cgo enabled, which is the normal case for local
    [development](development.md).

The schema is migrated automatically at startup; there is no separate migrate command.

## Uploads

Two settings work together:

- `TICKER_UPLOAD_PATH` is where files are written. It must be a **persistent, writable** directory,
  otherwise attachments are lost when the container is replaced while the database still references
  them.
- `TICKER_UPLOAD_URL` is the **public base URL of the API**, and is used to build absolute
  attachment links of the form `<TICKER_UPLOAD_URL>/media/<file>`.

```shell
TICKER_UPLOAD_PATH=/data/uploads
TICKER_UPLOAD_URL=https://api.ticker.example.org
```

`TICKER_UPLOAD_URL` must be the API's own public hostname, with **no `/v1`** and **no trailing
slash**. Media is served by the API at its root, so:

- `https://api.ticker.example.org/v1` produces `…/v1/media/x`, which is a 404.
- Pointing it at the frontend produces links the frontend answers with its own HTML page rather
  than an image.

These URLs are generated per response rather than stored, so correcting the value fixes existing
messages too, once the response cache expires.

Uploads accept `image/jpeg`, `image/gif` and `image/png` only, and the API rejects request bodies
over 10 MB. If a reverse proxy in front of it imposes a smaller limit, uploads fail there first —
nginx defaults to 1 MB, for instance.

## Metrics

Prometheus metrics are served on a **separate** listener, `metrics_listen` (`:8181` by default), at
`/metrics`. Do not expose it publicly; the compose stack keeps it on the internal network only.

## Caching

Some public responses are cached in memory: `/v1/init` and `/v1/feed` for 5 minutes, `/v1/timeline`
for 10 seconds. Configuration changes can therefore take up to five minutes to become visible on
the public page.
