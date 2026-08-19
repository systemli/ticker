# Integrations

Besides its own public page, a ticker can push every message to Telegram, Mastodon, Bluesky and
Signal groups.

!!! important "Integrations are configured at runtime, not in a config file"

    All credentials live in the database and are managed through the admin interface. There are no
    environment variables and no `config.yml` keys for them. Older documentation described
    `telegram:` and `signal_group:` config blocks and variables such as `TICKER_TELEGRAM_TOKEN` —
    these no longer exist and are silently ignored if present.

There are two levels:

| Integration | Instance-wide setup | Per-ticker setup |
| --- | --- | --- |
| Telegram | required (bot token) | channel |
| Signal | required (signal-cli endpoint) | group |
| Mastodon | none | server and credentials |
| Bluesky | none | handle and app password |

Telegram and Signal need an instance-wide step by a super admin before editors can use them. Until
that is done, the admin interface hides them.

## Telegram

**Instance-wide.** Create a bot with [@BotFather](https://t.me/BotFather) and note the token. As a
super admin, enter it in the admin interface under the Telegram settings.

Stored as the `telegram_settings` record, with fields `token` and `botUsername`.

**Per ticker.** Create a Telegram channel, add the bot as an administrator with permission to post,
then set the channel name on the ticker.

## Signal

**Instance-wide.** Signal support talks to [signal-cli](https://github.com/AsamK/signal-cli) running
in daemon mode, which is **not part of the Ticker stack** — you run it yourself.

!!! important "signal-cli's own JSON-RPC interface is required"

    Ticker calls signal-cli's native JSON-RPC methods (`send`, `updateGroup`, `listGroups`,
    `quitGroup`, `remoteDelete`) directly. A REST wrapper around signal-cli exposes different
    endpoints and will **not** work — point Ticker at the official image's `--http` endpoint.

Add it to the stack on the internal network:

```yaml
  signal-cli:
    image: ghcr.io/asamk/signal-cli:0.13.2
    # The published image is amd64 only.
    platform: linux/amd64
    restart: unless-stopped
    # The image's entrypoint is already
    # "signal-cli --config=/var/lib/signal-cli", so only the subcommand is given.
    # --http must bind 0.0.0.0: the default is localhost, which is unreachable
    # from other containers.
    command: ["daemon", "--http=0.0.0.0:8080"]
    volumes:
      - signal-cli-data:/var/lib/signal-cli
    networks:
      - internal
```

…and add `signal-cli-data:` to the `volumes:` block at the bottom of the file.

The container runs as UID `999` and `/var/lib/signal-cli` already exists in the image, so a fresh
named volume inherits the right ownership — unlike the API's uploads volume, this one needs no
preparation.

### Register the account first

signal-cli needs a registered or linked Signal account before the daemon is useful. Do this once,
against the same volume, following
[signal-cli's own documentation](https://github.com/AsamK/signal-cli/wiki):

```shell
# Link this instance as a secondary device to an existing Signal account
docker compose run --rm signal-cli link -n "Ticker"
```

This prints a `sgnl://linkdevice?uuid=…` URI. Turn it into a QR code (for example with `qrencode`)
and scan it from Signal on your phone under *Linked devices*.

Registering a dedicated number instead of linking is also possible; see signal-cli's documentation,
as it involves a captcha and an SMS verification step.

Afterwards start the daemon and confirm it answers:

```shell
docker compose up -d signal-cli
docker compose exec ticker-init sh -c \
  'wget -qO- --post-data="{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"listGroups\",\"params\":{\"account\":\"+491234567890\"}}" \
   --header="Content-Type: application/json" http://signal-cli:8080/api/v1/rpc'
```

In daemon mode signal-cli continuously receives messages, which the Signal protocol requires for
group functionality to work reliably.

### Configure Ticker

As a super admin, fill in the Signal settings in the admin interface:

| Field | Meaning |
| --- | --- |
| `apiUrl` | Full URL of signal-cli's JSON-RPC endpoint, including the path: `http://signal-cli:8080/api/v1/rpc` |
| `account` | The registered phone number in international format, e.g. `+491234567890` |
| `avatar` | Optional path to an avatar image for created groups |

The endpoint path `/api/v1/rpc` is not optional — it is where `--http` serves JSON-RPC, and Ticker
POSTs to exactly the URL you configure.

`avatar` is passed through to signal-cli, which reads it from **its own** filesystem. The path must
therefore exist inside the signal-cli container, not in the Ticker container. Mount the file in if
you want one.

Stored as the `signal_group_settings` record. Signal is considered enabled once **both** `apiUrl`
and `account` are set; `avatar` is optional.

**Per ticker.** Enabling Signal on a ticker creates a Signal group and returns an invite link that
you can publish. Ticker admins can be promoted to group admins. Messages posted to the ticker are
sent to the group, and deleting a message removes it there too.

## Mastodon

Configured entirely per ticker: the server URL and the application credentials. Register an
application on your Mastodon instance and enter the resulting values on the ticker. It is
considered connected once the token, secret and access token are all present.

## Bluesky

Configured entirely per ticker: the handle and an **app password** — generate one in the Bluesky
settings rather than using your account password.

Optionally restrict who may reply to the posts. Values are `followers`, `following`, `mentioned`
and `nobody`; leave empty to allow anyone. Several can be combined with commas, for example
`followers,mentioned`.

## Behaviour

Dispatch happens when a message is created, and integrations are also told about deletions and
ticker updates. A failing integration is logged and does **not** block the message from appearing
on the ticker itself, so a broken Telegram token will not take your public page down. Check the API
logs if a message did not arrive somewhere:

```shell
docker compose logs ticker | grep bridge_name
```

Attachments are sent along as files, read straight from `TICKER_UPLOAD_PATH` — no public URL is
involved, so an integration keeps working even if the interfaces are unreachable.
