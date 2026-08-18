# Installation

This guide installs the complete Ticker stack with Docker Compose, using
[Traefik](https://traefik.io/) as reverse proxy so certificates are obtained and renewed
automatically.

## Requirements

- A host with Docker and the Compose plugin.
- Ports **80** and **443** reachable from the internet. Port 80 is required for the Let's Encrypt
  HTTP challenge, even though all traffic is redirected to HTTPS.
- **Three hostnames** pointing at the host, each with an `A` record and — please — an `AAAA`
  record:

    | Example | Purpose |
    | --- | --- |
    | `ticker.example.org` | the public page readers visit |
    | `admin.ticker.example.org` | the admin interface |
    | `api.ticker.example.org` | the API, media files and feeds |

!!! note "Why three hostnames?"

    The admin and frontend are separate applications, and the API needs its own name because it
    serves attachments at `/media/...` and RSS feeds directly to readers. You can use any names
    you like; only their DNS records and your `.env` need to agree.

## 1. Get the files

```shell
git clone https://github.com/systemli/ticker.git
cd ticker
cp .env.example .env
```

## 2. Configure

Edit `.env`. Every value below must be set — `docker compose` refuses to start otherwise and tells
you which one is missing.

```shell title=".env"
--8<-- ".env.example"
```

Generate the JWT secret once and keep it:

```shell
openssl rand -hex 32
```

!!! warning "Always set `TICKER_SECRET`"

    Without it the API generates a **new random secret every time it starts**. Every admin is
    silently logged out on each restart, upgrade, or crash. This is the single most common
    misconfiguration.

## 3. Start

```shell
docker compose up -d
```

Watch Traefik obtain the certificates:

```shell
docker compose logs -f traefik
```

Then check that the API is alive:

```shell
curl https://api.ticker.example.org/healthz
# OK
```

The stack that is started looks like this:

```yaml title="compose.yaml"
--8<-- "compose.yaml"
```

## 4. Create the first user

There is **no default account and no generated password**. Create a super admin explicitly:

```shell
docker compose run --rm ticker user create \
  --email admin@example.org --super-admin
```

Omit `--password` and one is generated and printed:

```
Created user 1
Password: 6mXq...
```

Copy it now — it is not stored anywhere in readable form.

`--super-admin` is required for the first account. Only super admins can create tickers or change
integration settings.

Log in at `https://admin.ticker.example.org` and change the password.

## 5. Create a ticker and register its origin

In the admin interface, create a ticker. Then open its **Websites** configuration and add the
public address of your frontend, exactly:

```
https://ticker.example.org
```

!!! danger "This step is what makes the public page work"

    The API decides which ticker to serve from the browser's `Origin` header, and it compares the
    value **literally** against the origins you registered. Use scheme and host only:

    - :white_check_mark: `https://ticker.example.org`
    - :x: `https://ticker.example.org/` — a trailing slash never matches
    - :x: `https://ticker.example.org/live` — no paths
    - :x: `http://...` when the site is served over HTTPS

    If it does not match, the frontend shows *"The ticker is currently inactive"* rather than an
    error. Add one entry per domain if the same ticker is served on several.

Finally mark the ticker **active**, otherwise the same inactive page is shown.

## 6. Verify end to end

```shell
# Through the frontend, exactly as a browser does it
curl -s https://ticker.example.org/api/init | jq .data.ticker

# Directly against the API, supplying the origin yourself
curl -s -H 'Origin: https://ticker.example.org' \
  https://api.ticker.example.org/v1/init | jq .data.ticker
```

Both must return your ticker rather than `null`.

Now open the frontend, post a message from the admin interface, and confirm it appears **without
reloading the page** — that proves the realtime WebSocket connection works. Upload an image to a
message and confirm it renders, which proves `TICKER_UPLOAD_URL` is correct.

## Next steps

- Connect [Integrations](integrations.md) such as Telegram or Mastodon.
- Read [Operations](operations.md) for upgrades and backups.
- If something is not working, [Troubleshooting](troubleshooting.md) lists the symptoms.

## How the routing works

You do not need this to run the stack, but it helps when adapting it.

The published admin and frontend images contain a **relative** API base URL, `/api`. They are used
unmodified, and Traefik does two things for the `/api/` path on each of their hostnames:

1. **Rewrites the path.** `/api/**` becomes `/v1/**`, via `stripPrefix` followed by `addPrefix`.
2. **Sets the `Origin` header** to that hostname.

Step 2 is not cosmetic. Browsers do **not** send an `Origin` header on same-origin `GET` requests,
so without it the API cannot tell which ticker is being requested and returns the inactive page for
every visitor. It would also collapse its response cache into a single shared entry across all
tickers.

The API's own hostname is routed straight through with no rewriting, because `/media/...`, `/feed`
and `/healthz` are served outside the `/v1` prefix.

## Using a different reverse proxy

Any proxy works, provided it does all of the following for `/api/` on the admin and frontend
hostnames:

- rewrites `/api/**` to `/v1/**`;
- sets `Origin` to the public origin of that hostname;
- forwards WebSocket upgrades (`Connection`, `Upgrade`, HTTP/1.1) for `/api/ws`;
- allows request bodies of at least 10 MB, which is the API's own limit;
- and serves the API on its own hostname for `/media/**` and feeds.

An nginx equivalent of the `/api/` block:

```nginx
map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
}

server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    client_max_body_size 10m;

    location /api/ {
        proxy_pass http://ticker:8080/v1/;
        proxy_set_header Origin $scheme://$http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_read_timeout 3600s;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## Hardening

The stack is deliberately minimal. Two things are worth adding for a public deployment:

**Restrict the admin interface.** It is reachable by anyone who knows the hostname. Traefik can
limit it by IP:

```yaml
labels:
  - traefik.http.middlewares.admin-allow.ipallowlist.sourcerange=203.0.113.0/24
  - traefik.http.routers.admin.middlewares=admin-allow
```

Note this protects the admin *interface* only. The API hostname must stay public for media and
feeds, so apply the same middleware to the `admin-api` router if you want the API path restricted
too.

**Avoid mounting the Docker socket directly.** Traefik reads it to discover containers, and even
mounted read-only it is equivalent to root on the host. In a more sensitive setup, put a
[socket proxy](https://github.com/Tecnativa/docker-socket-proxy) in front of it and expose only
container listings.
