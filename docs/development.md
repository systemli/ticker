# Development

The project is three repositories:

| Repository | Stack |
| --- | --- |
| [ticker](https://github.com/systemli/ticker) | Go — Gin, GORM, Cobra |
| [ticker-admin](https://github.com/systemli/ticker-admin) | TypeScript — React, Vite, MUI |
| [ticker-frontend](https://github.com/systemli/ticker-frontend) | TypeScript — React, Vite, Tailwind |

Each repository carries its conventions in `AGENTS.md` — project layout, code style, testing
patterns and commit format. `CLAUDE.md` is a symlink to it, so coding agents pick up the same file
regardless of which name they look for.

## The whole stack, quickly

The repository ships a development stack that mirrors the production topology, over plain HTTP on
`*.localhost` hostnames, building the API from your working tree:

```shell
git clone https://github.com/systemli/ticker.git
cd ticker
docker compose -f compose.dev.yaml up -d --build
```

| URL | Service |
| --- | --- |
| <http://ticker.localhost> | public frontend |
| <http://admin.ticker.localhost> | admin interface |
| <http://localhost:8081> | Traefik dashboard |

The API has no address of its own; it answers under `/api` on both hostnames.

Create a user, then set up a ticker:

```shell
docker compose -f compose.dev.yaml run --rm ticker user create \
  --email dev@localhost --password devpassword123 --super-admin
```

Log in at <http://admin.ticker.localhost>, create a ticker, register the origin
`http://ticker.localhost` under its websites, mark it active, and open <http://ticker.localhost>.

!!! note "If port 5432 is already taken"

    The stack publishes PostgreSQL so you can also reach it from the host. If you already run
    PostgreSQL locally:

    ```shell
    TICKER_DEV_DB_PORT=5433 docker compose -f compose.dev.yaml up -d
    ```

Tear it down, including data:

```shell
docker compose -f compose.dev.yaml down -v
```

## Working on the API

For a fast edit-and-run loop, skip the container:

```shell
go run . run
```

This uses the built-in defaults, which means SQLite in `./ticker.db`. That works locally because your
toolchain has cgo enabled — unlike the released binaries, which cannot use SQLite at all.

```shell
# Create a user for your local instance
go run . user create --email dev@localhost --password devpassword123 --super-admin

curl http://localhost:8080/healthz
```

To use the development stack's PostgreSQL instead:

```shell
docker compose -f compose.dev.yaml up -d postgres

TICKER_DATABASE_TYPE=postgres \
TICKER_DATABASE_DSN='host=127.0.0.1 port=5432 user=ticker password=ticker dbname=ticker sslmode=disable TimeZone=UTC' \
  go run . run
```

### Tests and checks

```shell
go test ./...                                              # everything
go test -coverprofile=coverage.txt -covermode=atomic ./...  # with coverage
go test -run TestTickerTestSuite ./internal/api/...         # one suite

golangci-lint run --timeout 10m
gofmt -w . && goimports -w .
mockery                                                     # regenerate mocks
```

### Docker images

`Dockerfile` is the release image. It is built `FROM scratch` and **copies an already-built binary**,
because goreleaser cross-compiles it and injects the version through linker flags. Building it by
hand therefore requires the binary first:

```shell
go build && docker build .
```

`Dockerfile.dev` is the development image and does compile from source. It is what
`compose.dev.yaml` uses.

## Working on the interfaces

Both use **npm** and Node 24 (`.nvmrc` selects it):

```shell
nvm use
npm install
npm run dev
```

| Repository | Dev server |
| --- | --- |
| ticker-admin | <http://localhost:3000> |
| ticker-frontend | <http://localhost:4000> |

Both dev servers proxy `/api` to `http://localhost:8080/v1`, so run the API alongside them with
`go run . run` and nothing needs configuring.

!!! warning "Delete a leftover `.env`"

    `TICKER_API_URL` overrides the proxy with an absolute address. Requests then work but attachment
    images do not, because their URLs are relative and resolve against the dev server instead. The
    file is gitignored, so an old one may still be sitting in your checkout.

!!! warning "Register the dev server's own origin"

    The proxy sends the **dev server's** address as `Origin`. For the public frontend that means the
    ticker needs `http://localhost:4000` registered under its websites, or you will only ever see
    the inactive page. The admin interface is unaffected.

Other commands, in both repositories:

```shell
npm test           # vitest, watch mode
npm run coverage
npm run lint
npm run build
npm run preview
```

`ticker-admin` additionally has `npm run tsc` for a standalone type check.

## Conventions

Commit messages and pull request titles use [Gitmoji](https://gitmoji.dev/):

| | Meaning |
| --- | --- |
| ✨ | new feature |
| 🐛 | bug fix |
| ♻️ | refactor |
| ✅ | tests |
| ⬆️ | dependency upgrade |
| 📝 | documentation |
| 🧹 | chore |

Releases are drafted automatically from merged pull requests, so label them accordingly
(`feature`, `enhancement`, `fix`, `chore`, `dependencies`).

## Documentation

These pages live in `docs/` in the `ticker` repository and are built with MkDocs Material:

```shell
pip install -r requirements.txt
mkdocs serve          # http://localhost:8000
mkdocs build --strict # fails on broken links
```

`requirements.txt` is a lock file with pinned versions and hashes, generated from
`requirements.in`. To change a documentation dependency, edit `requirements.in` and regenerate:

```shell
pip install pip-tools
pip-compile --generate-hashes --allow-unsafe requirements.in
```

`--strict` is worth running before opening a pull request. Pages must be listed in the `nav` section
of `mkdocs.yml`, and code blocks that embed repository files with `--8<--` fail the build if the file
is moved.

Publishing happens automatically on every push to `main`.
