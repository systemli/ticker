# Troubleshooting

Organised by what you actually see.

## The public page says "The ticker is currently inactive"

By far the most common problem. The API could not match the request to a ticker, and falls back to
the inactive page rather than reporting an error.

Check, in order:

1. **Is the ticker marked active?** An inactive ticker shows exactly this page.
2. **Is the frontend's address registered on the ticker, character for character?** The comparison
   is literal:

    - :white_check_mark: `https://ticker.example.org`
    - :x: `https://ticker.example.org/` — a trailing slash in the **stored** value never matches
    - :x: `https://ticker.example.org/live` — no paths
    - :x: `http://…` when the site is served over HTTPS

3. **Is the proxy passing an `Origin` header?** Browsers omit it on same-origin `GET` requests, so
   the proxy has to add it. Compare the two responses:

    ```shell
    # With an origin: "settings" comes back populated
    curl -s -H 'Origin: https://ticker.example.org' \
      https://api.ticker.example.org/v1/init

    # Without one: "settings" is empty — this is what a missing Origin looks like
    curl -s https://api.ticker.example.org/v1/init
    ```

    If the request through your frontend (`https://ticker.example.org/api/init`) looks like the
    second, the proxy is not setting `Origin`.

!!! note "Changes take up to five minutes"

    `/v1/init` is cached for 5 minutes, so a correct fix can look like it did nothing. Verify with
    the `curl` commands above, which bypass the browser cache, or wait it out.

## Everything in the interface is empty, with no error

The `/api/` path is not routed to the API. The web server answers with the single-page app's own
HTML instead, the frontend fails to parse it, and silently renders nothing.

```shell
# Should return JSON, not HTML
curl -s https://ticker.example.org/api/init
```

Check that your proxy maps `/api/**` to `/v1/**` on the frontend and admin hostnames.

## Every API call returns 404

The API base URL is missing its `/v1` prefix. The applications request paths like
`<base>/admin/users`, so the base must end in `/v1` — for example `http://ticker:8080/v1`.

With the supplied Traefik stack this is handled by the path rewrite and needs no configuration.

## Admins are logged out after every restart

`TICKER_SECRET` is not set, so the API generates a new random signing secret each time it starts,
invalidating all existing tokens.

```shell
openssl rand -hex 32
```

Set it in `.env` and recreate the container. See [Configuration](configuration.md#the-signing-secret).

## The container will not start: "could not connect to database"

Read the accompanying `error` field.

**`go-sqlite3 requires cgo to work. This is a stub`**

: SQLite cannot be used with the official image or the released binaries. Switch to PostgreSQL or
  MySQL — see [Configuration](configuration.md#sqlite).

**`unknown time zone Etc/UTC`**

: The image has no timezone database. Use `TimeZone=UTC` in the connection string, not a named zone.

**`connection refused`**

: The database is not up yet or the host is wrong. The supplied stack waits for a health check;
  confirm with `docker compose ps`.

## New messages only appear after reloading the page

The realtime WebSocket connection is not getting through. It is opened at `/api/ws` on the
frontend's own hostname, so the proxy must forward the upgrade.

```shell
curl -i -N \
  -H 'Connection: Upgrade' -H 'Upgrade: websocket' \
  -H 'Sec-WebSocket-Version: 13' -H 'Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==' \
  https://ticker.example.org/api/ws
```

A working setup answers `HTTP/1.1 101 Switching Protocols`. If it returns 200 or 404 instead, the
proxy is not upgrading the connection. With nginx that means `proxy_http_version 1.1` plus the
`Upgrade` and `Connection` headers; Traefik handles it natively.

Also make sure the proxy does not time out idle connections aggressively — the server pings every
54 seconds, so a 60 second read timeout leaves almost no margin.

## Images are broken

**Everywhere, including in Telegram or Mastodon posts** — `TICKER_UPLOAD_URL` is wrong. It must be
the API's public base URL, with no `/v1` and no trailing slash:

```shell
TICKER_UPLOAD_URL=https://api.ticker.example.org
```

Verify a link resolves:

```shell
curl -I https://api.ticker.example.org/media/<uuid>.png
```

Note `/media` is served at the API's root, so it is **not** reachable through the `/api` rewrite —
`https://ticker.example.org/api/media/...` is expected to 404.

**Only for older messages** — the uploads directory was not persistent and the files are gone, while
the database still references them. Confirm `TICKER_UPLOAD_PATH` points into a named volume, and
restore the files from a backup.

## Uploads fail

**With a generic form error** — the uploads directory is not writable. The API runs as UID `10001`,
and a freshly created volume is owned by `root`. The supplied stack fixes ownership with its
`ticker-init` service; check that it completed:

```shell
docker compose logs ticker-init
docker compose logs ticker | grep -i mkdir
```

**Only for larger files** — something in front of the API imposes a smaller body limit than its own
10 MB. nginx defaults to 1 MB; raise it with `client_max_body_size 10m`.

**For a valid image that is rejected** — only JPEG, GIF and PNG are accepted.

## Certificates are not issued

```shell
docker compose logs traefik | grep -i acme
```

Usual causes: port 80 not reachable from the internet (required for the HTTP challenge, even though
traffic is redirected to HTTPS), DNS not yet pointing at the host, or the Let's Encrypt rate limit
after repeated failures. Confirm each hostname resolves to this machine before retrying.

## Traefik returns 404 for everything

Every hostname answers 404 while the Traefik dashboard itself works. Traefik discovered no
containers, so it has no routes — the dashboard is served internally and does not depend on
discovery.

```shell
docker compose logs traefik | grep -i "permission denied"
```

If that reports `permission denied while trying to connect to the docker API at
unix:///var/run/docker.sock`, Traefik cannot read the socket.

**On SELinux hosts — Podman, Fedora, RHEL — this happens even though the file permissions are
correct.** SELinux denies the container access to the mounted socket and the denial surfaces as a
permission error. Disabling the label for the Traefik container is enough, and does not relabel
anything on the host:

```yaml
  traefik:
    security_opt:
      - label=disable
```

`compose.dev.yaml` already sets this. Add it to your production stack only if you hit the error;
it is a no-op on hosts without SELinux.

To confirm the socket is reachable at all:

```shell
docker run --rm --security-opt label=disable \
  -v /var/run/docker.sock:/var/run/docker.sock alpine:3.24 \
  test -w /var/run/docker.sock && echo reachable
```

Otherwise check that the socket is mounted, and that its path is right for your runtime — rootless
Podman also exposes a per-user socket under `$XDG_RUNTIME_DIR/podman/podman.sock`.

Once discovery works, the routers appear here:

```shell
curl -s http://localhost:8081/api/http/routers | jq '.[].name'   # dev stack
```

You should see `ticker-api`, `frontend`, `frontend-api`, `admin` and `admin-api`, each suffixed
`@docker`. A suffix of `@internal` only means Traefik's own dashboard routes.

## Getting more detail

```shell
docker compose logs ticker | tail -50
```

Set `TICKER_LOG_LEVEL=debug` and `TICKER_LOG_FORMAT=text` while investigating. If you believe you
have found a bug, please [open an issue](https://github.com/systemli/ticker/issues) with the
relevant log lines and your configuration, credentials removed.
